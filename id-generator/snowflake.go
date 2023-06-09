package main

import (
	"errors"
	"math"
	"sync"
	"time"
)

const (
	BitLenTime          = 39
	BitLenMachineID     = 10
	BitLenSequence      = 63 - BitLenTime - BitLenMachineID
	TimeWindowInMiliSec = 10
	BitShiftMachineID   = BitLenSequence
	BitShiftTime        = BitShiftMachineID + BitLenMachineID
)

var (
	MaxTime     int64 = int64(math.Pow(float64(2), float64(BitLenTime))) - 1
	MaxSequence int64 = int64(math.Pow(float64(2), float64(BitLenSequence))) - 1

	// epoch is the start time
	epoch int64 = time.Date(2023, 3, 23, 0, 0, 0, 0, time.Now().Location()).UnixMilli() / 10
)

var (
	ErrTimeExceeded     = errors.New("Time Exceed Error")
	ErrSequenceExceeded = errors.New("Sequence Exceed Error")
)

type ID int64

type Generator struct {
	mu                       *sync.Mutex
	startTime                int64
	elapsedTime              int64
	sequence                 int64
	machineIDsequenceShifted int64
}

func newGenerator(startTime int64, machineID uint16) *Generator {
	return &Generator{
		mu:                       new(sync.Mutex),
		startTime:                startTime,
		elapsedTime:              time.Now().UnixMilli() / 10,
		sequence:                 0,
		machineIDsequenceShifted: int64(machineID) << BitShiftMachineID,
	}
}

func createID(currentTime int64, machineIDsequenceShifted int64, sequence int64) (ID, error) {
	if currentTime > MaxTime {
		return 0, ErrTimeExceeded
	} else if sequence > MaxSequence {
		return 0, ErrSequenceExceeded
	} else {
		id := ID((currentTime << BitShiftTime) | machineIDsequenceShifted | sequence)
		// fmt.Println(id, currentTime, machineIDsequenceShifted, sequence)
		return id, nil
	}
}

func (g *Generator) generate() (ID, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	var currentTime int64 = time.Now().UnixMilli() / 10
	if currentTime == g.elapsedTime {
		g.sequence += 1
	} else {
		g.sequence = 0
		g.elapsedTime = currentTime
	}

	id, err := createID(g.elapsedTime-epoch, g.machineIDsequenceShifted, g.sequence)
	if err != nil {
		return 0, err
	} else {
		return id, nil
	}
}
