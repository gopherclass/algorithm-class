package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"unsafe"
)

var treeBits uint

var runnerMap = map[string]DrawRunner{
	"red-black":   redblacktree{},
	"avl":         avltree{},
	"binary-tree": binarytree{},
	"patricia":    patriciatree{},
	"digital":     digitaltree{},
	"radix":       radixtree{},
}

func proc() error {
	var size int
	var aux bool
	var iteration uint
	var scale string
	flag.IntVar(&size, "size", 500, "size")
	flag.BoolVar(&aux, "aux", true, "draw auxiliary lines")
	flag.UintVar(&iteration, "iter", 10, "iteration")
	flag.StringVar(&scale, "scale", "", "scale: log, square")
	flag.UintVar(&treeBits, "bits", uint(unsafe.Sizeof(int(0))*8-2),
		"tree bits")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(0)
	}
	draw := newDrawRunner(defaultStyle, size)
	var options drawOptions
	switch scale {
	case "log":
		options.scale = logScale{}
	case "square":
		options.scale = squareScale{}
	}
	for _, arg := range args {
		runner, ok := runnerMap[arg]
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown runner %s\n", arg)
			os.Exit(0)
		}
		err := draw.draw(runner, iteration, &options)
		if err != nil {
			return err
		}
	}
	if aux {
		draw.drawAux()
	}
	return draw.store()
}

func printUsage() {
	names := make([]string, 0, len(runnerMap))
	for name := range runnerMap {
		names = append(names, name)
	}
	sort.Strings(names)
	fmt.Fprintln(os.Stderr, "need runner")
	fmt.Fprintln(os.Stderr, "supported runners =", strings.Join(names, ", "))
}

func main() {
	err := proc()
	if err != nil {
		log.Fatal(err)
	}
}
