package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
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

func showProgress(gid int, i int, total int) {
	if i%50000 == 0 {
		fmt.Printf("Goroutine %d. Processed %d of %d items...\n", gid, i, total)
	}
}

func benchmark(opts options) error {
	db, err := kv.NewStore(opts.engine, opts.path)
	if err != nil {
		return err
	}

	fmt.Printf("Number of keys: %d\n", opts.numKeys)
	fmt.Printf("Minimum key size: %d, maximum key size: %d\n", opts.minKeySize, opts.maxKeySize)
	fmt.Printf("Minimum value size: %d, maximum value size: %d\n", opts.minValueSize, opts.maxValueSize)
	fmt.Printf("Concurrency: %d\n", opts.concurrency)
	fmt.Printf("Running %s benchmark...\n", opts.engine)

	keys := generateKeys(opts.numKeys, opts.minKeySize, opts.maxKeySize)
	valSrc := make([]byte, opts.maxValueSize)
	if _, err := rand.Read(valSrc); err != nil {
		return err
	}
	forceGC()

	// Put.
	start := time.Now()
	err = concurrentBatch(keys, opts.concurrency, func(gid int, batch [][]byte) error {
		rnd := rand.New(rand.NewSource(int64(rand.Uint64())))
		for i, k := range batch {
			if err := db.Put(k, randValue(rnd, valSrc, opts.minValueSize, opts.maxValueSize)); err != nil {
				return err
			}
			if opts.progress {
				showProgress(gid, i, len(batch))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	endsecs := time.Since(start).Seconds()
	totalalsecs := endsecs
	fmt.Printf("Put: %.3f sec, %d ops/sec\n", endsecs, int(float64(opts.numKeys)/endsecs))

	// Reopen DB.
	if err := db.Close(); err != nil {
		return err
	}
	shuffle(keys)
	db, err = kv.NewStore(opts.engine, opts.path)
	if err != nil {
		return err
	}
	forceGC()

	// Read.
	start = time.Now()
	err = concurrentBatch(keys, opts.concurrency, func(gid int, batch [][]byte) error {
		for i, k := range batch {
			v, err := db.Get(k)
			if err != nil {
				return err
			}
			if v == nil {
				return errors.New("key doesn't exist")
			}
			if opts.progress {
				showProgress(gid, i, len(batch))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	endsecs = time.Since(start).Seconds()
	totalalsecs += endsecs
	fmt.Printf("Get: %.3f sec, %d ops/sec\n", endsecs, int(float64(opts.numKeys)/endsecs))

	// Total stats.
	fmt.Printf("Put + Get time: %.3f sec\n", totalalsecs)
	if err := db.Close(); err != nil {
		return err
	}
	sz, err := dirSize(opts.path)
	if err != nil {
		return err
	}
	fmt.Printf("File size: %s\n", byteSize(sz))
	return nil
}
