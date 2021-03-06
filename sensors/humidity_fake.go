package sensors

import "time"

const (
	fakeHygrometerFrq = 2 * time.Second
	fakeHygrometerMin = 40
	fakeHygrometerMax = 90
)

type fakeHygrometer struct {
	results chan Humidity
}

func NewFakeHygrometer(frq time.Duration) Hygrometer {
	fake := &fakeHygrometer{
		results: make(chan Humidity),
	}

	go fake.start(frq)

	return fake
}

func (f *fakeHygrometer) start(frq time.Duration) {
	for {
		f.results <- f.nextTemp()
		<-time.After(frq)
	}
}

func (f *fakeHygrometer) nextTemp() Humidity {
	return Humidity(theRand.Float64()*(fakeHygrometerMax-fakeHygrometerMin) + fakeHygrometerMin)
}

func (f *fakeHygrometer) Read() <-chan Humidity {
	return f.results
}
