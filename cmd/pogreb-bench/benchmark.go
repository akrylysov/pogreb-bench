package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/akrylysov/pogreb-bench/kv"
)

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func randKey(minL int, maxL int) string {
	n := rand.Intn(maxL-minL+1) + minL
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		buf[i] = byte(rand.Intn(95) + 32)
	}
	return string(buf)
}

func randValue(rnd *rand.Rand, src []byte, minS int, maxS int) []byte {
	n := rnd.Intn(maxS-minS+1) + minS
	return src[:n]
}

func forceGC() {
	runtime.GC()
	time.Sleep(time.Millisecond * 500)
}

func shuffle(a [][]byte) {
	for i := len(a) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func generateKeys(count int, minL int, maxL int) [][]byte {
	keys := make([][]byte, 0, count)
	seen := make(map[string]struct{}, count)
	for len(keys) < count {
		k := randKey(minL, maxL)
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		keys = append(keys, []byte(k))
	}
	return keys
}

func concurrentBatch(keys [][]byte, concurrency int, cb func(gid int, batch [][]byte) error) error {
	eg := &errgroup.Group{}
	batchSize := len(keys) / concurrency
	for i := 0; i < concurrency; i++ {
		batchStart := i * batchSize
		batchEnd := (i + 1) * batchSize
		if batchEnd > len(keys) {
			batchEnd = len(keys)
		}
		gid := i
		batch := keys[batchStart:batchEnd]
		eg.Go(func() error {
			return cb(gid, batch)
		})
	}
	return eg.Wait()
}

func showProgress(cur int, total int) {
	const (
		width float32 = 40
		freq          = 10000
	)
	complete := int(float32(cur) / float32(total) * width)
	if cur%freq != 0 {
		return
	}
	fmt.Printf("\r[%-40s] %d/%d", strings.Repeat("-", complete), cur, total)
}

func clearLine() {
	fmt.Print("\r\033[K")
}

func benchmarkPut(opts options, db kv.Store, keys [][]byte) error {
	valSrc := make([]byte, opts.maxValueSize)
	if _, err := rand.Read(valSrc); err != nil {
		return err
	}

	var keysProcessed int64
	err := concurrentBatch(keys, opts.concurrency, func(gid int, batch [][]byte) error {
		rnd := rand.New(rand.NewSource(int64(rand.Uint64())))
		for _, k := range batch {
			if err := db.Put(k, randValue(rnd, valSrc, opts.minValueSize, opts.maxValueSize)); err != nil {
				return err
			}
			showProgress(int(atomic.AddInt64(&keysProcessed, 1)), opts.numKeys)
		}
		return nil
	})
	if err != nil {
		return err
	}

	showProgress(int(keysProcessed), opts.numKeys)
	clearLine()
	return nil
}

func benchmarkGet(opts options, db kv.Store, keys [][]byte) error {
	var keysProcessed int64
	err := concurrentBatch(keys, opts.concurrency, func(gid int, batch [][]byte) error {
		for _, k := range batch {
			v, err := db.Get(k)
			if err != nil {
				return err
			}
			if v == nil {
				return errors.New("key doesn't exist")
			}
			showProgress(int(atomic.AddInt64(&keysProcessed, 1)), opts.numKeys)
		}
		return nil
	})
	if err != nil {
		return err
	}

	showProgress(int(keysProcessed), opts.numKeys)
	clearLine()
	return nil
}

func benchmark(opts options) error {
	db, err := kv.NewStore(opts.engine, opts.path)
	if err != nil {
		return err
	}

	fmt.Printf("engine: %s\n", opts.engine)
	fmt.Printf("keys: %d\n", opts.numKeys)
	fmt.Printf("key size: %d-%d\n", opts.minKeySize, opts.maxKeySize)
	fmt.Printf("value size %d-%d\n", opts.minValueSize, opts.maxValueSize)
	fmt.Printf("concurrency: %d\n\n", opts.concurrency)

	keys := generateKeys(opts.numKeys, opts.minKeySize, opts.maxKeySize)
	clearLine()

	var totalElapsed float64

	// Put.
	forceGC()
	start := time.Now()
	if err := benchmarkPut(opts, db, keys); err != nil {
		return err
	}
	elapsed := time.Since(start).Seconds()
	totalElapsed += elapsed
	fmt.Printf("put: %.3fs\t%d ops/s\n", elapsed, int(float64(opts.numKeys)/elapsed))

	// Reopen DB.
	if err := db.Close(); err != nil {
		return err
	}
	shuffle(keys)
	db, err = kv.NewStore(opts.engine, opts.path)
	if err != nil {
		return err
	}

	// Get.
	forceGC()
	start = time.Now()
	if err := benchmarkGet(opts, db, keys); err != nil {
		return err
	}
	elapsed = time.Since(start).Seconds()
	totalElapsed += elapsed
	fmt.Printf("get: %.3fs\t%d ops/s\n", elapsed, int(float64(opts.numKeys)/elapsed))

	// Total stats.
	fmt.Printf("\nput + get: %.3fs\n", totalElapsed)
	if err := db.Close(); err != nil {
		return err
	}
	sz, err := dirSize(opts.path)
	if err != nil {
		return err
	}
	fmt.Printf("file size: %s\n", byteSize(sz))
	return nil
}
