package collector

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"
	"os"
)


type Event struct {
	Type string    				`json: "type"`
	Time time.Time 				`json: "time"`
	Date map[string]interface{}  `json: "date"`
}


type Collector struct {
	mu      sync.Mutex
	file *	os.File
	logger *log.Logger
	closed 	bool
} 

func NewCollector(logger *log.Logger) *Collector {
	f, err := os.OpenFile("events.jsonl", os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0644)
	
	if err != nil {
		logger.Fatal("Open events file:", err)
	}

	return &Collector {
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

	b, _:= json.Marshal(e)
	_, _ = s.file.Write(append(b, '\n'))

	s.logger.Printf("Ingested: %s %s \n", e.Type, e.Date["path"])
}


func (s *Collector) Stop(ctx context.Context){
	s.mu.Lock()

	defer s.mu.Unlock()

	if s.closed {
		return
	}

	_ = s.file.Close()
	s.closed = true
	s.logger.Println("Collector stopped")
}