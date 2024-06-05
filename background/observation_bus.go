package background

import (
	// Context Bus
	"os"

	"github.com/AleckDarcy/ContextBus/configure"
	"github.com/AleckDarcy/ContextBus/configure/observation"
	"github.com/AleckDarcy/ContextBus/context"
	"github.com/AleckDarcy/ContextBus/helper"
	cb "github.com/AleckDarcy/ContextBus/proto"
	"github.com/rs/zerolog"

	// third-party
	"github.com/AleckDarcy/ContextBus/third-party/github.com/opentracing/opentracing-go"
	"github.com/AleckDarcy/ContextBus/third-party/github.com/uber/jaeger-client-go"
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

var todoLoggingConfigure = &observation.LoggingConfigure{
	Out: cb.LogOutType_Stdout,
}

func (b *observationBus) doObservation() (cnt, cntL, cntT, cntM int) {
	for {
		v, ok := b.queue.Dequeue()
		if !ok {
			return
		}

		var start int64
		if cb.PERF_METRIC {
			start = time.Now().UnixNano()
		}

		pay := v.(*EventDataPayload)
		var cntL_, cntT_, cntM_ int
		if cfg := configure.Store.GetConfigure(pay.ctx.GetRequestContext().GetConfigureID()); cfg != nil {
			if obs := cfg.GetObservationConfigure(pay.ed.Event.Recorder.Name); obs != nil {
				cntL_, cntT_, cntM_ = obs.Do(pay.ctx, pay.ed)
				cntL += cntL_
				cntT += cntT_
				cntM += cntM_
			}

			// check after observation alerts
			if rac := cfg.GetReaction(pay.ed.Event.Recorder.Name); rac != nil {
				if rac.Type == cb.ReactionType_ReactionPrintLog { // todo
					prevName := rac.PreTree.Nodes[0].PrevEvent.GetName()
					prevED := pay.ed.GetPreviousEventData(prevName)
					if prevED != nil {
						latency := (pay.ed.Event.When.Time - prevED.Event.When.Time) / int64(time.Millisecond)
						if latency > rac.PreTree.Nodes[0].PrevEvent.GetLatency() {
							pay.ctx.GetRequestContext()
							tags := map[string]interface{}{
								"RequestID": pay.ed.GetMetadata().ReqId,
								"EventID":   pay.ed.GetMetadata().EveId,
							}

							span := pay.ctx.GetTracer().StartSpan(prevName, opentracing.StartTime(time.Unix(0, prevED.Event.When.Time)), opentracing.Tags(tags))
							span.FinishWithOptions(opentracing.FinishOptions{FinishTime: time.Unix(0, pay.ed.Event.When.Time)})

							// print logs in reversed order
							fmt.Printf("report high latency %d ms, from %s to %s, tags %v\n", latency, prevName, pay.ed.Event.Recorder.Name, tags)
							todoLoggingConfigure.Do(pay.ed)

							esp := int64(-1)

							buf := make([]byte, 512)
							for prevED != nil {
								buf = buf[0:0]
								if prevED.Metadata.Esp != esp {
									es := EnvironmentProfiler.GetByID(prevED.Metadata.Esp)
									esp = prevED.Metadata.Esp

									if es != nil {
										encoder := helper.JSONEncoder
										buf = encoder.BeginObject(buf)

										buf = helper.JSONEncoder.AppendKey(buf, "caller")
										buf = helper.JSONEncoder.AppendString(buf, "environmental profile")

										buf = helper.JSONEncoder.AppendKey(buf, "level")
										buf = helper.JSONEncoder.AppendString(buf, "warn")

										buf = helper.JSONEncoder.AppendKey(buf, "time")

										buf = helper.JSONEncoder.BeginString(buf)
										buf = time.Unix(0, esp).AppendFormat(buf, helper.TIME_FORMAT_DEFAULT)
										buf = helper.JSONEncoder.EndString(buf)

										buf = encoder.AppendKey(buf, "message")
										buf = encoder.AppendString(buf, fmt.Sprintf("%v", es))
										buf = helper.JSONEncoder.AppendKey(buf, "ID")
										buf = helper.JSONEncoder.AppendIDs(buf, pay.ed.GetMetadata().ReqId, pay.ed.GetMetadata().EveId)
										buf = encoder.EndObject(buf)

										str := helper.BytesToString(buf)
										fmt.Fprintln(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}, str)
									}
								}

								todoLoggingConfigure.Do(prevED)
								prevED = prevED.PrevEventData
							}
						}
					}
				}
			}
		}

		if cb.PERF_METRIC {
			end := time.Now().UnixNano()

			cbType := int64(0)
			if cntL_ != 0 {
				cbType |= cb.CBType_Logging
			}
			if cntT_ != 0 {
				cbType |= cb.CBType_Tracing
			}
			if cntM_ != 0 {
				cbType |= cb.CBType_Metrics
			}

			cb.PMetric.AddCBLatency(cbType, pay.ed.Event.When.Time, start, end)
		}

		cnt++
	}
}

type observationCounter struct {
	payload int
}

func (b *observationBus) Run(cfg *configure.ServerConfigure, sig chan struct{}) {
	var tracerCfg = &config.Configuration{
		ServiceName: cfg.ServiceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		}, Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  cfg.JaegerHost,
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
