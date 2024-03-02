package observation

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"testing"
	"time"
)

func TestOpenTracing(t *testing.T) {
	jsOK := jaegerFakeServer
	defer jsOK.Close()

	tracer, err := NewTracer("service", jsOK.URL[7:], 1)
	if err != nil {
		t.Log(err)
	}

	start1 := time.Now().UnixNano()
	finish1 := start1 + int64(10*time.Second)
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
	start2 := finish1 + int64(time.Second)
	finish2 := start2 + int64(10*time.Second)
	span2 := tracer.StartSpan("span", sso, opentracing.StartTime(time.Unix(0, start2)), opentracing.Tags{"method": "POST"})
	span2.FinishWithOptions(opentracing.FinishOptions{FinishTime: time.Unix(0, finish2)})
	t.Log(span2.(*jaeger.Span))
}
