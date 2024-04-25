package data

import (
	"context"
	"fmt"
)

func (p *KeySpecificDataProcessor) Schedule(d DataWithKey) {
	if d.Key != p.ProcessorKey {
		panic(fmt.Sprintf("Data with key %s scheduled for KeySpecificDataProcessor with key %s", d.Key, p.ProcessorKey))
	}

	p.c <- d
}

func (p *KeySpecificDataProcessor) StartProcessing(ctx context.Context) {
	// read from queue until end
	// or until context cancelled
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context timeout exceeded")
			return
		case data := <-p.c:
			fmt.Printf("received data: %v\n", data)
			return
		}

	}
}

func (p *KeySpecificDataProcessor) StopProcessing() {
	// TODO send signal to current task to stop via ctx
	close(p.c)
}

func CreateAndStart(ctx context.Context, processorKey string) *KeySpecificDataProcessor {
	p := KeySpecificDataProcessor{
		ProcessorKey: processorKey,
		c:            make(chan DataWithKey),
	}

	p.StartProcessing(ctx)

	return &p
}
