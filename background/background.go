package background

// List of background tasks:
// 1. environmental profiling
// 2. (?) configuration
// Initialized during deployment.

type Configure struct {
	EnvironmentProfiler bool
	ObservationBus      bool
}

var cfg *Configure

type signal struct {
	environmentProfiler chan struct{}
	observationBus      chan struct{}
}

var sig = signal{
	environmentProfiler: make(chan struct{}, 1),
	observationBus:      make(chan struct{}, 1),
}

func Run(cfg_ *Configure) {
	if cfg = cfg_; cfg == nil {
		cfg = &Configure{}

		return
	}

	if cfg.EnvironmentProfiler {
		go EnvironmentProfiler.Run(sig.environmentProfiler)
	}

	if cfg.ObservationBus {
		go ObservationBus.Run(sig.observationBus)
	}
}

func Stop() {
	if cfg.EnvironmentProfiler {
		sig.environmentProfiler <- struct{}{}
	}

	if cfg.ObservationBus {
		sig.observationBus <- struct{}{}
	}
}
