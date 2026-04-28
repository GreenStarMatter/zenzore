package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBadSignalLookup(t *testing.T) {
	constantSignalValue := 1
	testSignal := &Signal{Type: "sldkjfI", ExpectedValue: constantSignalValue}
	outputSignal := GenerateSignalSample(testSignal)
	assert.Equal(t, BAD_SAMPLE, outputSignal)
}

func TestSignalGenerateLookup(t *testing.T) {
	constantSignalValue := 1
	testSignal := &Signal{Type: "constant", ExpectedValue: constantSignalValue}
	outputSignal := GenerateSignalSample(testSignal)
	assert.Equal(t, constantSignalValue, outputSignal)
}
