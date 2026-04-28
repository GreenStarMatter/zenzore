package main

const BAD_SAMPLE = Sample(-9999)

type Sample int

type Signal struct {
	Type          string
	ExpectedValue int
	RandomValue   int
}

func GenerateConstantSignal(expected int) Sample {
	return Sample(expected)
}

func GenerateSignalSample(incomingSignal *Signal) Sample {
	switch incomingSignal.Type {
	case "constant":
		return GenerateConstantSignal(incomingSignal.ExpectedValue)
	default:
		return BAD_SAMPLE
	}
}
