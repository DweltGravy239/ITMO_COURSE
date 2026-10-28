package main

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

func envInt(name string, def int) int {
	if v, err := strconv.Atoi(os.Getenv(name)); err == nil {
		return v
	}
	return def
}

func main() {
	cpu := envInt("CPU_WORKERS", 1)
	memMB := envInt("MEM_MB", 64)

	// память: выделяем и записываем по байту в каждую страницу (4 КБ),
	// иначе ядро не выдаст её физически
	buf := make([]byte, memMB*1024*1024)
	for i := 0; i < len(buf); i += 4096 {
		buf[i] = 1
	}

	// CPU: по одной горутине на каждый «поток» нагрузки, бесконечный цикл
	for i := 0; i < cpu; i++ {
		go func() {
			for {
			}
		}()
	}

	log.Printf("batch started: cpu=%d mem=%dMB", cpu, memMB)
	for {
		time.Sleep(10 * time.Second)
		runtime.KeepAlive(buf)
	}
}
