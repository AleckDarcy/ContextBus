package observation

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"github.com/uber/jaeger-client-go/config"

	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"time"
)

var (
	defaultSampleRatio float64 = 0.01
)

var jaegerFakeServer = httptest.NewServer(
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		lastBody, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		fmt.Printf("jaeger fake server received payload, len %d, body %v\n", len(lastBody), lastBody)
		w.WriteHeader(http.StatusOK)
	}),
)

// NewTracer returns a newly configured tracer
func NewTracer(serviceName, host string, ratio float64) (opentracing.Tracer, error) {
	if ratio <= 0 {
		ratio = defaultSampleRatio
	}
	if val, ok := os.LookupEnv("JAEGER_SAMPLE_RATIO"); ok {
		ratio, _ = strconv.ParseFloat(val, 64)
	}
	if ratio > 1 {
		ratio = 1.0
	}

	log.Info().Msgf("Jaeger client: adjusted sample ratio %f", ratio)
	tempCfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "probabilistic",
			Param: ratio,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  host,
		},
	}

	log.Info().Msg("Overriding Jaeger config with env variables")
	cfg, err := tempCfg.FromEnv()
	if err != nil {
		return nil, err
	}

	tracer, _, err := cfg.NewTracer()
	if err != nil {
		return nil, err
	}
	return tracer, nil
}
