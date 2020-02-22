package main

import (
	"context"
	"runtime"
	"time"
)

func monitorMemory(interval time.Duration) (*uint64, context.CancelFunc) {
	ticker := time.NewTicker(interval)
	ctx, cancel := context.WithCancel(context.Background())
	var peakSys uint64
	go func() {
		for {
			select {
			case <-ticker.C:
				ms := runtime.MemStats{}
				runtime.ReadMemStats(&ms)
				if ms.Sys > peakSys {
					peakSys = ms.Sys
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
	return &peakSys, cancel
}
