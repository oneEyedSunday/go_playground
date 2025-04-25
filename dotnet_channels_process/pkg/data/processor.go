package data

import (
	"context"
	"fmt"
	"time"
)

type KeySpecificDataProcessor struct {
	processorKey string
	// c holds data to be worked on
	c chan DataWithKey
	// abortC is a cancel function that is used internally to signal to the work functions to abort processing
	// abortC blends both the context of CreateAndStart with a cancel function that is run on `StopProcessing`
	abortC                      context.CancelFunc
	processingFinishedTimestamp time.Time
}

func (p *KeySpecificDataProcessor) GetProcessorKey() string {
	return p.processorKey
}

func (p *KeySpecificDataProcessor) LastProcessingTimestamp() time.Time {
	if !p.processingFinishedTimestamp.IsZero() {
		return p.processingFinishedTimestamp
	}
	return time.Now().UTC()
}

func (p *KeySpecificDataProcessor) Schedule(d DataWithKey) {
	if d.Key != p.processorKey {
		panic(fmt.Sprintf("Data with key %s scheduled for KeySpecificDataProcessor with key %s", d.Key, p.processorKey))
	}

	p.c <- d
}

func (p *KeySpecificDataProcessor) StartProcessing(ctx context.Context) {
	// read from queue until end
	// or until context cancelled
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[processor][%s] context timeout exceeded\n", p.processorKey)
			p.StopProcessing()
			return
		case data := <-p.c:
			p.processingFinishedTimestamp = time.Now().UTC()
			p.Work(ctx, data)
			return
		}

	}
}

func (p *KeySpecificDataProcessor) Work(ctx context.Context, data DataWithKey) {
	fmt.Printf("[processor][%s] received data: %v\n", p.processorKey, data)

}

func (p *KeySpecificDataProcessor) StopProcessing() {
	// close(p.c)
	// do not receive any more work
	p.c = nil
	// but do not close since the pump may continue to push work in
	// close(p.c)
	p.abortC()
}

func CreateAndStart(ctx context.Context, processorKey string) *KeySpecificDataProcessor {
	p := KeySpecificDataProcessor{
		processorKey: processorKey,
		c:            make(chan DataWithKey),
	}

	// can now pass in context from above to worker
	// this is safe as StartProcessing is called on init
	abortableCtx, cancel := context.WithCancel(ctx)

	// and call p.abortC to ensure it cleanups
	p.abortC = cancel
	go p.StartProcessing(abortableCtx)

	return &p
}

var _ IDataProcessor = (*KeySpecificDataProcessor)(nil)
