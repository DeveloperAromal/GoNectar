package collector

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Event struct {
	Type string                 `json:"type"`
	Time time.Time              `json:"time"`
	Data map[string]interface{} `json:"data"`
}

type Collector struct {
	mu     sync.Mutex
	file   *os.File
	logger *log.Logger
	closed bool
}

func NewCollector(logger *log.Logger) *Collector {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("failed to get executable path: %v", err)
	}

	rootDir := filepath.Dir(filepath.Dir(exePath))
	eventsDir := filepath.Join(rootDir, "events")

	if err := os.MkdirAll(eventsDir, os.ModePerm); err != nil {
		logger.Fatalf("failed to create events directory: %v", err)
	}

	filePath := filepath.Join(eventsDir, "events.jsonl")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatalf("failed to open events file: %v", err)
	}

	logger.Printf("Collector initialized, writing to: %s\n", filePath)

	return &Collector{
		file:   f,
		logger: logger,
	}
}

func (s *Collector) IngestEvent(e Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	b, _ := json.Marshal(e)
	_, _ = s.file.Write(append(b, '\n'))
	s.logger.Printf("Ingested: %s %v\n", e.Type, e.Data["path"])
}

func (s *Collector) Stop(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	_ = s.file.Close()
	s.closed = true
	s.logger.Println("Collector stopped")
}
