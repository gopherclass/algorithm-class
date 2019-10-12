package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"golang.org/x/xerrors"
)

type targetFunc func(context.Context) error

var targets = make(map[string]targetFunc)

func registerSorter(sorter sorter) {
	forTest := func(ctx context.Context) error {
		runTest(sorter)
		return nil
	}
	forDraw := func(ctx context.Context) error {
		const maxsize = 500
		const iteration = 10
		res := benchmark(sorter, maxsize, iteration)
		for _, record := range res.records {
			err := documentDrawing().storeRecord(record)
			if err != nil {
				return err
			}
		}
		return nil
	}
	register("test-"+sorter.epithet(), forTest)
	register("draw-"+sorter.epithet(), forDraw)
}

func register(name string, target targetFunc) {
	targets[name] = target
}

func separateArgs(args []string) ([]string, []string) {
	for i, arg := range args {
		if arg == "--" {
			return args[:i], args[i+1:]
		}
	}
	return args, nil
}

func invalidTargetNames(targetNames []string) (invalids []string) {
	for _, targetName := range targetNames {
		if targets[targetName] == nil {
			invalids = append(invalids, targetName)
		}
	}
	return invalids
}

func proc(ctx context.Context) error {
	targetArgs, args := separateArgs(os.Args[1:])
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	err := fs.Parse(targetArgs)
	if err != nil {
		return err
	}
	targetNames := fs.Args()
	invalids := invalidTargetNames(targetNames)
	if len(invalids) > 0 {
		errorInvalidTargets(invalids)
	}
	if len(targetNames) == 0 {
		helpTargets()
		abort()
	}
	flag.CommandLine.Parse(args)
	return executeTargets(ctx, targetNames)
}

func errorInvalidTargets(invalids []string) {
	fmt.Fprintf(os.Stderr, "invalid targets: %s\n", invalids)
	helpTargets()
	abort()
}

func helpTargets() {
	validTargetNames := make([]string, 0, len(targets))
	for targetName := range targets {
		validTargetNames = append(validTargetNames, targetName)
	}
	sort.Strings(validTargetNames)
	fmt.Fprintf(os.Stderr, "available targets:\n")
	for _, targetName := range validTargetNames {
		fmt.Fprintf(os.Stderr, "  - %s\n", targetName)
	}
}

func executeTargets(ctx context.Context, targetNames []string) error {
	for _, targetName := range targetNames {
		err := targets[targetName](ctx)
		if err != nil {
			return xerrors.Errorf("when executing target %s: %w",
				targetName, err)
		}
	}
	return nil
}

func main() {
	ctx := context.Background()
	err := proc(ctx)
	if err != nil {
		log.Fatalf("%+v", err)
		abort()
	}
}

func abort() {
	os.Exit(1)
}
