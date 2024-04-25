package data

import (
	"context"
	"fmt"
	"sync"
)

type BackgroundDataProcessor struct {
	c  chan DataWithKey
	p  map[string]*KeySpecificDataProcessor
	mu sync.RWMutex
}

func NewDataProcessor() *BackgroundDataProcessor {
	return &BackgroundDataProcessor{
		c: make(chan DataWithKey),
		p: make(map[string]*KeySpecificDataProcessor),
	}
}

// GetOrCreateDataProcessor is a threadsafe method that either returns or create a processor based on the specified key
func (p *BackgroundDataProcessor) GetOrCreateDataProcessor(ctx context.Context, key string) *KeySpecificDataProcessor {
	// should lock this here
	// to guard against unsynchronized concurrent access
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.p[key]; !exists {
		p.p[key] = CreateAndStart(ctx, key)
	}

	return p.p[key]
}

// Push receives data into the backlog, internally it has a buffer to hold data, it should ideally be called in a goroutine to avoid blocking on congestion
func (p *BackgroundDataProcessor) Push(entry DataWithKey) {
	// fmt.Printf("before pushing size: %v and cap %v\n", len(p.c), cap(p.c))
	p.c <- entry
	// fmt.Printf("after pushing size: %v and cap %v\n", len(p.c), cap(p.c))
}

func (p *BackgroundDataProcessor) Execute(ctx context.Context) {
	fmt.Println("executing background consumer")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[background] context timeout exceeded")
			return
		case data := <-p.c:
			fmt.Printf("[background] received data: %v\n", data)
			fmt.Printf("after receiving size: %v and cap %v\n", len(p.c), cap(p.c))
			processor := p.GetOrCreateDataProcessor(ctx, data.Key)
			go processor.Schedule(data)
		}

	}
}
