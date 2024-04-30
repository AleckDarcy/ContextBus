package background

import (
	// Context Bus
	"github.com/AleckDarcy/ContextBus/configure"
	"github.com/AleckDarcy/ContextBus/configure/observation"
	"github.com/AleckDarcy/ContextBus/context"
	"github.com/AleckDarcy/ContextBus/helper"
	cb "github.com/AleckDarcy/ContextBus/proto"

	// third-party
	"github.com/AleckDarcy/ContextBus/third-party/github.com/opentracing/opentracing-go"
	"github.com/AleckDarcy/ContextBus/third-party/github.com/uber/jaeger-client-go/config"

	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

// EventDataPayload is the package of event data inside LockFreeQueue
type EventDataPayload struct {
	ctx *context.Context
	ed  *cb.EventData
}

type observationBus struct {
	queue  *helper.LockFreeQueue
	tracer opentracing.Tracer
	signal chan struct{}
	eveID  uint64
}

var ObservationBus = &observationBus{
	queue:  helper.NewLockFreeQueue(),
	signal: make(chan struct{}, 1),
}

func (b *observationBus) GetTracer() opentracing.Tracer {
	return b.tracer
}

func (b *observationBus) NewEventID() uint64 {
	return atomic.AddUint64(&b.eveID, 1)
}

func (b *observationBus) OnSubmit(ctx *context.Context, ed *cb.EventData) {
	b.queue.Enqueue(&EventDataPayload{
		ctx: ctx,
		ed:  ed,
	})

	// try to invoke
	select {
	case b.signal <- struct{}{}:
		// fmt.Println("notified")
		// message sent
	default:
		// fmt.Println("failed")
		// message dropped
	}
}

func (b *observationBus) doObservation() (cnt, cntL, cntT, cntM int) {
	for {
		v, ok := b.queue.Dequeue()
		if !ok {
			return
		}

		pay := v.(*EventDataPayload)
		if cfg := configure.Store.GetConfigure(pay.ctx.GetRequestContext().GetConfigureID()); cfg != nil {
			if obs := cfg.GetObservationConfigure(pay.ed.Event.Recorder.Name); obs != nil {
				cntL_, cntT_, cntM_ := obs.Do(pay.ctx, pay.ed)
				cntL += cntL_
				cntT += cntT_
				cntM += cntM_
			}
		}

		cnt++
	}
}

type observationCounter struct {
	payload int
}

func (b *observationBus) Run(sig chan struct{}) {
	var tracerCfg = &config.Configuration{
		ServiceName: "test-service",
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	tracer, closer, err := tracerCfg.NewTracer()
	if err != nil {
		panic(fmt.Sprintf("cannot init tracer: %v", err))
	}

	b.tracer = tracer

	// todo: report ready

	for {
		cnt, cntL, cntT, cntM := 0, 0, 0, 0
		select {
		case <-sig:
			closer.Close()

			return
		case <-b.signal: // triggered by collector notification
			cnt, cntL, cntT, cntM = b.doObservation()
		case <-time.After(helper.BUS_OBSERVATION_QUEUE_INTERVAL): // triggered by timer
			cnt, cntL, cntT, cntM = b.doObservation()
		}

		if cntM != 0 {
			observation.MetricVecStore.Push(observation.PrometheusPusher)
		}

		// fmt.Println("bus processed", cnt, "payloads")

		_, _, _, _ = cnt, cntL, cntT, cntM

		runtime.Gosched()
	}
}
