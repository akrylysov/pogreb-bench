package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/pkg/profile"
)

type options struct {
	engine       string
	numKeys      int
	minKeySize   int
	maxKeySize   int
	minValueSize int
	maxValueSize int
	concurrency  int
	path         string
	progress     bool
	compact      bool
	profileMode  string
}

func main() {
	var opts options
	flag.StringVar(&opts.engine, "e", "pogreb", "database engine name. pogreb, goleveldb, bbolt or badger")
	flag.IntVar(&opts.numKeys, "n", 100000, "number of keys")
	flag.IntVar(&opts.minKeySize, "mink", 16, "minimum key size")
	flag.IntVar(&opts.maxKeySize, "maxk", 64, "maximum key size")
	flag.IntVar(&opts.minValueSize, "minv", 128, "minimum value size")
	flag.IntVar(&opts.maxValueSize, "maxv", 512, "maximum value size")
	flag.IntVar(&opts.concurrency, "c", 1, "number of concurrent goroutines")
	flag.StringVar(&opts.path, "p", "", "database path")
	flag.BoolVar(&opts.progress, "progress", false, "show progress")
	flag.BoolVar(&opts.compact, "compact", false, "write keys twice and run compaction after")
	flag.StringVar(&opts.profileMode, "profile", "", "enable profile. cpu, mem, block or mutex")

	flag.Parse()

	if opts.maxKeySize < opts.minKeySize {
		opts.maxKeySize = opts.minKeySize
	}

	if opts.maxValueSize < opts.minValueSize {
		opts.maxValueSize = opts.minValueSize
	}

	if opts.path == "" {
		flag.Usage()
		return
	}

	switch opts.profileMode {
	case "cpu":
		defer profile.Start(profile.CPUProfile).Stop()
	case "mem":
		defer profile.Start(profile.MemProfile).Stop()
	case "block":
		defer profile.Start(profile.BlockProfile).Stop()
	case "mutex":
		defer profile.Start(profile.MutexProfile).Stop()
	}

	peakSysMem, cancelMon := monitorMemory(time.Millisecond * 100)
	if err := benchmark(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error running benchmark: %v\n", err)
	}
	cancelMon()
	fmt.Printf("Peak Sys Mem: %s\n", byteSize(*peakSysMem))
}
