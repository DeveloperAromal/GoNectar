package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DeveloperAromal/GoNectar/internal/collector"
	"github.com/DeveloperAromal/GoNectar/internal/config"
	"github.com/DeveloperAromal/GoNectar/utils"
	httptrap "github.com/DeveloperAromal/GoNectar/internal/trap"

)

func main() {
	banner.Banner()
	cfg := config.Config{
		HTTPAddr: ":8080",
	}

	logger := log.Default()

	col := collector.NewCollector(logger)

	trap := httptrap.NewHTTPTrap(&cfg, col, logger)
	go func() {
		if err := trap.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("trap failed:", err)
		}
	}()

	logger.Printf("http trap listening on %s\n", cfg.HTTPAddr)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = trap.Stop(ctx)
	col.Stop(ctx)
	logger.Println("shutting down")
}
