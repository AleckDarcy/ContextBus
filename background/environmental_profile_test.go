package background

import (
	"runtime"
	"testing"
	"time"
)

func TestEnvironmentalProfile(t *testing.T) {
	t.Log(EnvironmentProfiler.latest)

	signal := make(chan struct{})
	go EnvironmentProfiler.Run(signal)

	<-time.After(11 * time.Second)
	t.Log(EnvironmentProfiler.latest)
}

func TestGetEnvironmentProfile(t *testing.T) {
	time.Sleep(time.Second)

	pf1 := EnvironmentProfiler.GetEnvironmentProfile()
	t.Logf("%+v", pf1)

	go runtime.GC()

	for i := 0; i < 20; i++ {
		pf := EnvironmentProfiler.GetEnvironmentProfile()
		t.Logf("%+v", pf)

		if i == 10 {
			go runtime.GC()
		}
	}
}

func BenchmarkGetEnvironmentProfile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EnvironmentProfiler.GetEnvironmentProfile()
	}
}
