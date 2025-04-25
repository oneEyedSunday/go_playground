package data

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type BackgroundDataProcessorMonitor struct {
	mu                            sync.Mutex
	log                           *log.Logger
	processors                    TProcess
	processorExpiryScanningPeriod time.Duration
	processorExpiryThreshold      time.Duration
	abortTaskC                    context.CancelFunc
}

type TProcess = map[string]*KeySpecificDataProcessor

func newDataProcessorMonitor(dataProcessors TProcess, l *log.Logger) *BackgroundDataProcessorMonitor {
	return &BackgroundDataProcessorMonitor{
		log:                           l,
		processors:                    dataProcessors,
		processorExpiryScanningPeriod: time.Duration(500 * time.Millisecond),
		processorExpiryThreshold:      time.Duration(1 * time.Second),
	}
}

func (m *BackgroundDataProcessorMonitor) StartMonitoring(ctx context.Context) {
	newCtx, newCtxCancelFn := context.WithCancel(ctx)
	m.abortTaskC = newCtxCancelFn
	ticker := time.NewTicker(m.processorExpiryScanningPeriod)
	go func() {
		for {
			select {
			case <-newCtx.Done():
				// cancellation from context
				fmt.Println("monitoring stopped")
				ticker.Stop()
				return
			case <-ticker.C:
				// our processor expiry scanner has ticked
				func() {
					m.log.Println("expiry period reached")
					m.mu.Lock()
					defer m.mu.Unlock()

					m.log.Printf("checking %d processors for work expiry\n", len(m.processors))

					for index, p := range m.processors {
						if m.isExpired(p) {
							// block till we stop processing
							m.log.Printf("processor (%s) is expired, stopping...\n", p.GetProcessorKey())
							p.StopProcessing()

							delete(m.processors, index)

							m.log.Printf("Removed data processor for data key %s\n", p.GetProcessorKey())
						}
					}

				}()
			}

		}
	}()

}

func (m *BackgroundDataProcessorMonitor) StopMonitoring() {
	// trigger the cancelFn of
	if m.abortTaskC != nil {
		m.abortTaskC()
	}

	m.abortTaskC = nil
}

func (m *BackgroundDataProcessorMonitor) isExpired(p *KeySpecificDataProcessor) bool {
	return time.Now().UTC().Sub(p.LastProcessingTimestamp()) > m.processorExpiryThreshold
}

func CreateAndStartMonitor(ctx context.Context, dataProcessors TProcess, l *log.Logger) *BackgroundDataProcessorMonitor {
	m := newDataProcessorMonitor(dataProcessors, l)

	m.StartMonitoring(ctx)
	return m
}
