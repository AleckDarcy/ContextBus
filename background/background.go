package background

import "github.com/AleckDarcy/ContextBus/configure"

// List of background tasks:
// 1. environmental profiling
// 2. (?) configuration
// Initialized during deployment.

var cfg *configure.ServerConfigure

type signal struct {
	environmentProfiler chan struct{}
	observationBus      chan struct{}
}

var sig = signal{
	environmentProfiler: make(chan struct{}, 1),
	observationBus:      make(chan struct{}, 1),
}

func Run(cfg_ *configure.ServerConfigure) {
	if cfg = cfg_; cfg == nil {
		cfg = &configure.ServerConfigure{}

		return
	}

	if cfg.EnvironmentProfiler {
		go EnvironmentProfiler.Run(sig.environmentProfiler)
	}

	if cfg.ObservationBus {
		go ObservationBus.Run(cfg, sig.observationBus)
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
