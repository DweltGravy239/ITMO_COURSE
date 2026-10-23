// api — HTTP-сервис магазина shop.
//
// Эндпоинты:
//
//	GET  /health   — "ok"; при HEALTH_FAIL=true отвечает 500 (нужно для проверки rollout/readiness)
//	POST /order    — создать заказ {"item":"coffee","qty":2}
//	GET  /orders   — последние 100 заказов
//	GET  /version  — значение APP_VERSION (удобно видеть при rolling update)
//	GET  /metrics  — метрики Prometheus
//
// Переменные окружения:
//
//	PORT          порт (по умолчанию 8080)
//	DATABASE_URL  строка подключения к Postgres (postgresql://user:pass@host:5432/db)
//	HEALTH_FAIL   "true" — /health отвечает 500
//	APP_VERSION   версия, которую приложение сообщает о себе (по умолчанию "dev")
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "shop_api_http_requests_total",
		Help: "Количество HTTP-запросов к api.",
	}, []string{"method", "path", "code"})

	httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "shop_api_http_request_duration_seconds",
		Help:    "Длительность обработки HTTP-запросов.",
		Buckets: prometheus.DefBuckets,
	}, []string{"path"})

	ordersCreated = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "shop_api_orders_created_total",
		Help: "Количество созданных заказов.",
	})

	dbUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "shop_api_db_up",
		Help: "1, если api подключился к базе и схема готова; иначе 0.",
	})
)

type server struct {
	db         *sql.DB
	dbReady    atomic.Bool
	healthFail bool
	version    string
}

type statusRecorder struct {
	http.ResponseWriter
	code int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.code = code
	r.ResponseWriter.WriteHeader(code)
}

// instrument оборачивает обработчик: считает запросы, меряет время, добавляет заголовок версии.
func (s *server) instrument(path string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, code: http.StatusOK}
		w.Header().Set("X-App-Version", s.version)
		h(rec, r)
		httpRequests.WithLabelValues(r.Method, path, strconv.Itoa(rec.code)).Inc()
		httpDuration.WithLabelValues(path).Observe(time.Since(start).Seconds())
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) {
	if s.healthFail {
		http.Error(w, "unhealthy", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("ok\n"))
}

func (s *server) versionHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(s.version + "\n"))
}

type orderIn struct {
	Item string `json:"item"`
	Qty  int    `json:"qty"`
}

func (s *server) createOrder(w http.ResponseWriter, r *http.Request) {
	if !s.dbReady.Load() {
		http.Error(w, "database is not ready", http.StatusServiceUnavailable)
		return
	}
	var in orderIn
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&in); err != nil {
		http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
		return
	}
	if in.Item == "" {
		http.Error(w, "field item is required", http.StatusBadRequest)
		return
	}
	if in.Qty == 0 {
		in.Qty = 1
	}
	if in.Qty < 0 {
		http.Error(w, "qty must be positive", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	var id int64
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO orders (item, qty) VALUES ($1, $2) RETURNING id`, in.Item, in.Qty).Scan(&id)
	if err != nil {
		log.Printf("insert order: %v", err)
		http.Error(w, "database error", http.StatusServiceUnavailable)
		return
	}
	ordersCreated.Inc()
	writeJSON(w, http.StatusCreated, map[string]any{"id": id, "item": in.Item, "qty": in.Qty, "status": "new"})
}

type orderOut struct {
	ID          int64      `json:"id"`
	Item        string     `json:"item"`
	Qty         int        `json:"qty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
}

func (s *server) listOrders(w http.ResponseWriter, r *http.Request) {
	if !s.dbReady.Load() {
		http.Error(w, "database is not ready", http.StatusServiceUnavailable)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, item, qty, status, created_at, processed_at FROM orders ORDER BY id DESC LIMIT 100`)
	if err != nil {
		log.Printf("select orders: %v", err)
		http.Error(w, "database error", http.StatusServiceUnavailable)
		return
	}
	defer rows.Close()
	out := []orderOut{}
	for rows.Next() {
		var o orderOut
		var processed sql.NullTime
		if err := rows.Scan(&o.ID, &o.Item, &o.Qty, &o.Status, &o.CreatedAt, &processed); err != nil {
			log.Printf("scan order: %v", err)
			http.Error(w, "database error", http.StatusServiceUnavailable)
			return
		}
		if processed.Valid {
			t := processed.Time
			o.ProcessedAt = &t
		}
		out = append(out, o)
	}
	writeJSON(w, http.StatusOK, out)
}

const schema = `
CREATE TABLE IF NOT EXISTS orders (
    id           BIGSERIAL PRIMARY KEY,
    item         TEXT        NOT NULL,
    qty          INT         NOT NULL DEFAULT 1,
    status       TEXT        NOT NULL DEFAULT 'new',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at TIMESTAMPTZ
)`

// ensureSchema создаёт таблицу под advisory-блокировкой: несколько реплик api
// стартуют одновременно и не должны спорить за CREATE TABLE.
func (s *server) ensureSchema(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(727001)`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, schema); err != nil {
		return err
	}
	return tx.Commit()
}

// connectLoop в фоне ждёт базу и готовит схему. Сервис при этом уже отвечает на /health:
// если бы api падал без базы, под бы уходил в CrashLoopBackOff при каждом сбое postgres.
func (s *server) connectLoop(ctx context.Context) {
	for ctx.Err() == nil {
		cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := s.db.PingContext(cctx)
		if err == nil {
			err = s.ensureSchema(cctx)
		}
		cancel()
		if err == nil {
			s.dbReady.Store(true)
			dbUp.Set(1)
			log.Println("database is ready")
			break
		}
		dbUp.Set(0)
		log.Printf("waiting for database: %v", err)
		select {
		case <-ctx.Done():
		case <-time.After(2 * time.Second):
		}
	}
	// После первого успеха следим за связью: если база пропала, фиксируем это в метрике и /order.
	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
		}
		cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		err := s.db.PingContext(cctx)
		cancel()
		if err != nil {
			if s.dbReady.Swap(false) {
				log.Printf("database lost: %v", err)
			}
			dbUp.Set(0)
		} else if !s.dbReady.Swap(true) {
			log.Println("database is back")
			dbUp.Set(1)
		}
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	prometheus.MustRegister(httpRequests, httpDuration, ordersCreated, dbUp)

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	s := &server{
		db:         db,
		healthFail: os.Getenv("HEALTH_FAIL") == "true",
		version:    getenv("APP_VERSION", "dev"),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go s.connectLoop(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.instrument("/health", s.health))
	mux.HandleFunc("GET /version", s.instrument("/version", s.versionHandler))
	mux.HandleFunc("POST /order", s.instrument("/order", s.createOrder))
	mux.HandleFunc("GET /orders", s.instrument("/orders", s.listOrders))
	mux.Handle("GET /metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:              ":" + getenv("PORT", "8080"),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		log.Println("shutting down: finishing in-flight requests")
		sctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(sctx)
	}()

	log.Printf("api %s listening on %s (health_fail=%v)", s.version, srv.Addr, s.healthFail)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
	_ = db.Close()
	log.Println("bye")
}
