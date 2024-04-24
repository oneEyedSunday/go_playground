package data

import (
	"context"
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

func (d *BackgroundDataProcessor) write(data DataWithKey) {
	d.c <- data
}

func (p *BackgroundDataProcessor) GetOrCreateDataProcessor(key string) *KeySpecificDataProcessor {
	// should lock this here
	// to guard against unsynchronized concurrent access
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.p[key]; !exists {
		p.p[key] = CreateAndStart(context.TODO(), key)
	}

	return p.p[key]
}

// func (d *DataProcessor) Schedule(f data.DataWithKey) chan DataWithKey {
// 	d.write(f)

// 	return d.c
// }

// func (d *DataProcessor) ReadAndFn(fn func(item data.DataWithKey) error) {
// 	for {
// 		entry := <-d.c
// 		processor := d.GetOrCreateDataProcessor(entry.Key)
// 		fmt.Printf("Processor found: %v \n", processor)
// 		// err := processor(data)
// 		// if err != nil {
// 		// 	fmt.Printf("error handling data %s\n", err.Error())
// 		// }
// 	}
// }
