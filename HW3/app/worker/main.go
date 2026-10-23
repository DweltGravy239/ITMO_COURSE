// worker — фоновый обработчик заказов shop.
//
// Раз в WORKER_INTERVAL забирает необработанные заказы (status='new') из postgres
// и помечает их обработанными. Таблицу orders создаёт api; пока её нет, worker ждёт.
//
// Чтобы к нему применимы те же пробы, что и к api, worker поднимает маленький HTTP-сервер:
//
//	GET /health   — "ok"; при HEALTH_FAIL=true отвечает 500
//	GET /metrics  — метрики Prometheus
//
// Переменные окружения:
//
//	PORT             порт HTTP-сервера (по умолчанию 8080)
//	DATABASE_URL     строка подключения к Postgres
//	WORKER_INTERVAL  пауза между проходами, например 2s (по умолчанию 2s)
//	HEALTH_FAIL      "true" — /health отвечает 500
package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	processed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "shop_worker_orders_processed_total",
		Help: "Количество заказов, помеченных обработанными.",
	})
	errorsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "shop_worker_errors_total",
		Help: "Количество неудачных проходов worker (ошибки базы).",
	})
	lastSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "shop_worker_last_success_timestamp_seconds",
		Help: "Unix-время последнего успешного прохода.",
	})
	pending = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "shop_orders_pending",
		Help: "Сколько заказов ждёт обработки (status='new').",
	})
)

const takeBatch = `
UPDATE orders
SET status = 'processed', processed_at = now()
WHERE id IN (
    SELECT id FROM orders
    WHERE status = 'new'
    ORDER BY id
    LIMIT 10
    FOR UPDATE SKIP LOCKED
)
RETURNING id`

// pass делает один проход: помечает до 10 заказов и обновляет gauge очереди.
// SKIP LOCKED позволяет нескольким репликам worker работать параллельно,
// не мешая друг другу и не обрабатывая один заказ дважды.
func pass(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, takeBatch)
	if err != nil {
		return err
	}
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return err
	}
	_ = rows.Close()
	if len(ids) > 0 {
		processed.Add(float64(len(ids)))
		log.Printf("processed orders: %v", ids)
	}

	var left int
	if err := db.QueryRowContext(ctx, `SELECT count(*) FROM orders WHERE status = 'new'`).Scan(&left); err != nil {
		return err
	}
	pending.Set(float64(left))
	return nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	prometheus.MustRegister(processed, errorsTotal, lastSuccess, pending)

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	interval, err := time.ParseDuration(getenv("WORKER_INTERVAL", "2s"))
	if err != nil || interval <= 0 {
		log.Fatalf("bad WORKER_INTERVAL: %v", err)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	db.SetMaxOpenConns(3)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	healthFail := os.Getenv("HEALTH_FAIL") == "true"
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		if healthFail {
			http.Error(w, "unhealthy", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.Handle("GET /metrics", promhttp.Handler())
	srv := &http.Server{Addr: ":" + getenv("PORT", "8080"), Handler: mux, ReadHeaderTimeout: 5 * time.Second}

	go func() {
		log.Printf("worker listening on %s, interval %s", srv.Addr, interval)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	// Ошибки базы (недоступна, таблицы ещё нет) не роняют процесс: worker логирует и пробует снова.
	// Иначе при каждом сбое postgres под уходил бы в CrashLoopBackOff.
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		if err := pass(pctx, db); err != nil {
			if ctx.Err() == nil {
				errorsTotal.Inc()
				log.Printf("pass failed: %v", err)
			}
		} else {
			lastSuccess.SetToCurrentTime()
		}
		cancel()

		select {
		case <-ctx.Done():
			log.Println("shutting down")
			sctx, scancel := context.WithTimeout(context.Background(), 5*time.Second)
			_ = srv.Shutdown(sctx)
			scancel()
			_ = db.Close()
			return
		case <-ticker.C:
		}
	}
}
