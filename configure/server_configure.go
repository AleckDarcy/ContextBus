package configure

type ServerConfigure struct {
	ServiceName         string
	JaegerHost          string
	EnvironmentProfiler bool
	ObservationBus      bool
}
