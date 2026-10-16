package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

var (
	mu      sync.Mutex
	held    [][]byte
	totalMB int
)

func main() {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("shutting down")
		os.Exit(0)
	}()
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	// Выделяет N МБ и держит их. Страницы трогаем явно, иначе ядро
	// не выделит физическую память (ленивое выделение) и RSS не вырастет.
	http.HandleFunc("/eat", func(w http.ResponseWriter, r *http.Request) {
		mb, err := strconv.Atoi(r.URL.Query().Get("mb"))
		if err != nil || mb <= 0 {
			http.Error(w, "usage: /eat?mb=N", http.StatusBadRequest)
			return
		}
		buf := make([]byte, mb<<20)
		for i := 0; i < len(buf); i += 4096 {
			buf[i] = 1
		}
		mu.Lock()
		held = append(held, buf)
		totalMB += mb
		total := totalMB
		mu.Unlock()
		fmt.Fprintf(w, "holding %d MB\n", total)
	})

	// Запускает бесконечный цикл на одной горутине, то есть грузит одно ядро.
	http.HandleFunc("/burn", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			for {
			}
		}()
		fmt.Fprintln(w, "burning one core")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("listening on :" + port)
	http.ListenAndServe(":"+port, nil)

}
