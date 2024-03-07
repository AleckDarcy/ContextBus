package observation

import (
	"github.com/AleckDarcy/ContextBus/third-party/github.com/opentracing/opentracing-go"
	"github.com/AleckDarcy/ContextBus/third-party/github.com/uber/jaeger-client-go"
	"github.com/AleckDarcy/ContextBus/third-party/github.com/uber/jaeger-client-go/config"

	"reflect"
	"testing"
	"time"
)

func TestOpenTracing(t *testing.T) {
	jsOK := jaegerFakeServer
	defer jsOK.Close()

	cfg := &config.Configuration{
		ServiceName: "test-service",
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	tracer, closer, err := cfg.NewTracer()
	defer closer.Close()
	if err != nil {
		t.Errorf("cannot init Jaeger: %v", err)
	}
	t.Log(reflect.TypeOf(tracer))
	t.Logf("%+v", tracer)
	//tracer, err := NewTracer("service", jsOK.URL[7:], 1)

	//tracer, err := NewTracer("test-service", "localhost:14268", 1)
	//if err != nil {
	//	t.Log(err)
	//}

	start1 := time.Now().UnixNano()
	finish1 := start1 + int64(3*time.Second)
	span1 := tracer.StartSpan("span", opentracing.StartTime(time.Unix(0, start1)), opentracing.Tags{"method": "POST"})
	span1.FinishWithOptions(opentracing.FinishOptions{FinishTime: time.Unix(0, finish1)})
	jSpan1 := span1.(*jaeger.Span)
	t.Log(jSpan1)

	//traceID := jSpan1.SpanContext().TraceID()
	//spanID := jSpan1.SpanContext().SpanID()

	sso := &opentracing.SpanReference{
		Type:              opentracing.ChildOfRef,
		ReferencedContext: jSpan1.SpanContext(),
	}

	start2 := start1 + int64(time.Second)
	finish2 := start2 + int64(1*time.Second)
	span2 := tracer.StartSpan("span", sso, opentracing.StartTime(time.Unix(0, start2)), opentracing.Tags{"method": "POST"})
	span2.FinishWithOptions(opentracing.FinishOptions{FinishTime: time.Unix(0, finish2)})
	t.Log(span2.(*jaeger.Span))
}
