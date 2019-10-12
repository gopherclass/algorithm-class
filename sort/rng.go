package main

import (
	"time"

	"golang.org/x/exp/rand"
)

var rngSeed = time.Now().UnixNano()

var rngSource = rand.New(rand.NewSource(uint64(rngSeed)))
