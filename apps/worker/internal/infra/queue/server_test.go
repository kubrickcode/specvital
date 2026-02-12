package queue

import (
	"testing"

	"github.com/riverqueue/river"
)

func TestBuildQueueConfig_MultiQueue(t *testing.T) {
	cfg := ServerConfig{
		Queues: []QueueAllocation{
			{Name: "analysis:priority", MaxWorkers: 5},
			{Name: "analysis:default", MaxWorkers: 3},
			{Name: "analysis:scheduled", MaxWorkers: 2},
		},
	}

	result := buildQueueConfig(cfg)

	if len(result) != 3 {
		t.Errorf("expected 3 queues, got %d", len(result))
	}

	tests := []struct {
		name       string
		maxWorkers int
	}{
		{"analysis:priority", 5},
		{"analysis:default", 3},
		{"analysis:scheduled", 2},
	}

	for _, tt := range tests {
		if q, ok := result[tt.name]; !ok {
			t.Errorf("queue %q not found", tt.name)
		} else if q.MaxWorkers != tt.maxWorkers {
			t.Errorf("queue %q: expected MaxWorkers %d, got %d", tt.name, tt.maxWorkers, q.MaxWorkers)
		}
	}
}

func TestBuildQueueConfig_MultiQueueWithZeroWorkers(t *testing.T) {
	cfg := ServerConfig{
		Queues: []QueueAllocation{
			{Name: "test:queue", MaxWorkers: 0},
		},
	}

	result := buildQueueConfig(cfg)

	if q, ok := result["test:queue"]; !ok {
		t.Error("queue not found")
	} else if q.MaxWorkers != DefaultConcurrency {
		t.Errorf("expected default concurrency %d, got %d", DefaultConcurrency, q.MaxWorkers)
	}
}

func TestBuildQueueConfig_MultiQueueWithEmptyName(t *testing.T) {
	cfg := ServerConfig{
		Queues: []QueueAllocation{
			{Name: "", MaxWorkers: 5},
		},
	}

	result := buildQueueConfig(cfg)

	if q, ok := result[river.QueueDefault]; !ok {
		t.Errorf("expected empty name to fall back to default queue %q", river.QueueDefault)
	} else if q.MaxWorkers != 5 {
		t.Errorf("expected MaxWorkers 5, got %d", q.MaxWorkers)
	}
}

func TestBuildQueueConfig_EmptyQueues(t *testing.T) {
	cfg := ServerConfig{
		Queues: []QueueAllocation{},
	}

	result := buildQueueConfig(cfg)

	if len(result) != 0 {
		t.Errorf("expected empty queue map, got %d queues", len(result))
	}
}
