package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var runnerMap = map[string]DrawRunner{
	"red-black":   redblacktree{},
	"avl":         avltree{},
	"binary-tree": binarytree{},
}

func proc() error {
	var size int
	var aux bool
	flag.IntVar(&size, "size", 1000, "size")
	flag.BoolVar(&aux, "aux", true, "draw auxiliary lines")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "need runner")
		os.Exit(0)
	}
	draw := newDrawRunner(defaultStyle, size)
	for _, arg := range args {
		runner, ok := runnerMap[arg]
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown runner %s\n", arg)
			os.Exit(0)
		}
		err := draw.draw(runner)
		if err != nil {
			return err
		}
	}
	if aux {
		draw.drawAux()
	}
	return draw.store()
}

func main() {
	err := proc()
	if err != nil {
		log.Fatal(err)
	}
}
