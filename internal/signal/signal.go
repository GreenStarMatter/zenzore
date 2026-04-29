package signal

const BAD_SAMPLE = Sample(-9999)

type Sample int

type Signal struct {
	Type          string
	ExpectedValue int
	RandomValue   int
}

type Device struct {
	SN        string `json:"SN"`
	PN        string `json:"PN"`
	Reading_1 int    `json:"Reading_1"`
	Reading_2 int    `json:"Reading_2"`
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
