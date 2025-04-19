package util

import (
	"log/slog"
	"runtime"
	"sync"
	"time"
)

func LogRuntimeStats(period time.Duration) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	slog.Info(
		"go runtime stats",
		"goroutines", runtime.NumGoroutine(),
		"cpu_cores", runtime.NumCPU(),
		"cgo_calls", runtime.NumCgoCall(),
		"allocated_memory_bytes", memStats.Alloc,
		"total_allocated_memory_bytes", memStats.TotalAlloc,
		"heap_allocated_memory_bytes", memStats.HeapAlloc,
		"heap_system_memory_bytes", memStats.HeapSys,
		"memory_from_system_bytes", memStats.Sys,
		"mallocs", memStats.Mallocs,
		"frees", memStats.Frees,
		"live_objects", memStats.Mallocs-memStats.Frees,
		"heap_objects", memStats.HeapObjects,
		"next_GC_bytes", memStats.NextGC,
		"total_GC_cycles", memStats.NumGC,
		"GC_pause_total_ns", memStats.PauseTotalNs,
		"Last_GC_pause_ns", memStats.PauseNs[(memStats.NumGC+255)%256],
	)
}

func LogRuntimeStatsBasic() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	slog.Info(
		"go runtime short stats",
		"goroutines", runtime.NumGoroutine(),
		"allocated_memory_bytes", memStats.Alloc,
	)
}

func Periodic(period time.Duration, funcs ...func()) {
	if len(funcs) == 0 {
		return
	}

	go func() {
		for {
			wg := sync.WaitGroup{}
			wg.Add(len(funcs))
			for i := range funcs {
				go func(i int) {
					defer wg.Done()
					funcs[i]()
				}(i)
			}
			wg.Wait()
			time.Sleep(period)
		}
	}()
}
