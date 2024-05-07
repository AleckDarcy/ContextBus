package observation

import (
	"github.com/AleckDarcy/ContextBus/context"
	"github.com/AleckDarcy/ContextBus/helper"
	cb "github.com/AleckDarcy/ContextBus/proto"
	"github.com/AleckDarcy/ContextBus/third-party/github.com/opentracing/opentracing-go"
	"github.com/AleckDarcy/ContextBus/third-party/github.com/uber/jaeger-client-go"
	"github.com/rs/zerolog"

	"fmt"
	"os"
	"time"
)

// Prepare pre-allocates SpanID and get stacktrace
func (c *Configure) Prepare(ctx *context.Context, ed *cb.EventData) {
	if c.Tracing == nil { // no tracing required
		return
	} else if c.Tracing.PrevName != "" { //
		return
	}

	if c.IsStacktrace(ctx) { // capture stacktrace
		// todo get stacktrace
	}

	fmt.Printf("request context: %+v\n", ctx.GetRequestContext())
	reqSM := ctx.GetRequestContext().GetSpanMetadata() // parent span from caller
	if reqSM == nil {
		fmt.Println("request span metadata not found")
		return
	} else if !reqSM.Sampled {
		fmt.Println("request not sampled")
		return
	}

	fmt.Println("set span metadata")
	sm := &cb.SpanMetadata{
		Sampled:     true,
		TraceIdHigh: reqSM.TraceIdHigh,
		TraceIdLow:  reqSM.TraceIdLow,
		SpanId:      ctx.GetTracer().RandomID(),
		ParentId:    reqSM.SpanId,
	}

	ed.SpanMetadata = sm
	ctx.SetSpanMetadata(sm)
}

// IsStacktrace returns true if any observation configuration needs stacktrace
func (c *Configure) IsStacktrace(ctx *context.Context) bool {
	if (*LoggingConfigure)(c.Logging).IsStacktrace(ctx) {
		return true
	} else if (*TracingConfigure)(c.Tracing).IsStacktrace(ctx) {
		return true
	}

	return false
}

func (c *StackTraceConfigure) IsStacktrace(ctx *context.Context) bool {
	if c == nil {
		return false
	}

	return c.Switch
}

func (c *LoggingConfigure) IsStacktrace(ctx *context.Context) bool {
	if c == nil {
		return false
	}

	return (*StackTraceConfigure)(c.Stacktrace).IsStacktrace(ctx)
}

func (c *TracingConfigure) IsStacktrace(ctx *context.Context) bool {
	if c == nil {
		return false
	}

	return (*StackTraceConfigure)(c.Stacktrace).IsStacktrace(ctx)
}

// the Do function
// finalize observation

// Do function returns number of log, trace span and metric entries
func (c *Configure) Do(ctx *context.Context, ed *cb.EventData) (cntL, cntT, cntM int) {
	cntL = (*LoggingConfigure)(c.Logging).Do(ed)
	cntT = (*TracingConfigure)(c.Tracing).Do(ctx, ed)
	cntM = len(c.Metrics)
	for _, metric := range c.Metrics {
		(*MetricsConfigure)(metric).Do(ed)
	}

	return
}

func (c *TimestampConfigure) Do() {
	if c != nil {
		// todo default
		return
	}
}

func (c *StackTraceConfigure) Do() {
	if c != nil {
		// todo default
		return
	}
}

func (c *LoggingConfigure) Do(ed *cb.EventData) int {
	if c == nil {
		return 0
	}

	er := ed.Event

	e := newEvent()
	e.buf = helper.JSONEncoder.BeginObject(e.buf)

	e.buf = helper.JSONEncoder.AppendKey(e.buf, "caller")
	e.buf = helper.JSONEncoder.AppendString(e.buf, "test/caller.go")

	e.buf = helper.JSONEncoder.AppendKey(e.buf, "level")
	e.buf = helper.JSONEncoder.AppendString(e.buf, "info")

	e.buf = helper.JSONEncoder.AppendKey(e.buf, "time")

	e.buf = helper.JSONEncoder.BeginString(e.buf)
	if ts := c.Timestamp; ts == nil {
		e.buf = time.Unix(0, er.When.Time).AppendFormat(e.buf, helper.TIME_FORMAT_DEFAULT)
	} else {
		e.buf = time.Unix(0, er.When.Time).AppendFormat(e.buf, ts.Format)
	}
	e.buf = helper.JSONEncoder.EndString(e.buf)

	// do message
	e.buf = helper.JSONEncoder.AppendKey(e.buf, "message")
	msg := er.What.Application.GetMessage()
	paths := er.What.Application.GetPaths()
	values := make([]interface{}, len(paths))
	var err error
	for i, path := range paths {
		values[i], err = er.What.GetValue(path)
		if err != nil {
			values[i] = fmt.Sprintf("!error(%s)", err.Error())
		}
	}
	e.buf = helper.JSONEncoder.AppendString(e.buf, fmt.Sprintf(msg, values...))

	// do tag
	if len(c.Attrs) != 0 {
		e.buf = helper.JSONEncoder.AppendKey(e.buf, "tags")
		e.buf = DoTagFaster(e.buf, c.Attrs, er)

		//tags := DoTag(c.Attrs, er)
		//e.buf = encoder.JSONEncoder.AppendTags(e.buf, tags)
	}

	e.buf = helper.JSONEncoder.EndObject(e.buf)
	str := string(e.buf)
	e.finalize()

	(*TimestampConfigure)(c.Timestamp).Do()
	(*StackTraceConfigure)(c.Stacktrace).Do()

	switch c.Out {
	case cb.LogOutType_LogOutType_:
		// omit print
	case cb.LogOutType_Stdout:
		fmt.Fprintln(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}, str)
	case cb.LogOutType_Stderr:
		fmt.Fprintln(os.Stderr, str)
	case cb.LogOutType_File:

	default:
		fmt.Fprintln(os.Stdout, str)
	}

	return 1
}

func (c *TracingConfigure) Do(ctx *context.Context, ed *cb.EventData) int {
	if c == nil {
		return 0
	} else if c.PrevName == "" { // skip
		return 0
	}

	prev := ed.GetPreviousEventData(c.PrevName)
	if prev == nil {
		fmt.Println("previous event not found", c.PrevName)
		return 0
	}

	sm := prev.SpanMetadata
	if sm == nil {
		fmt.Println("previous span metadata not found", c.PrevName)
		return 0
	}

	tags := DoTraceTag(c.Attrs, ed.Event)
	// todo: tags for testing
	tags["method"] = "POST"
	tags["key"] = "value"

	sr := opentracing.SpanReference{
		Type:              opentracing.ChildOfRef,
		ReferencedContext: sm.ParentSpanContext(),
	}

	fmt.Println("span reference", sr)

	span := ctx.GetTracer().StartSpan(c.Name, sr, opentracing.SpanID(sm.SpanId), opentracing.StartTime(time.Unix(0, prev.Event.When.Time)), opentracing.Tags(tags))
	span.FinishWithOptions(opentracing.FinishOptions{FinishTime: time.Unix(0, ed.Event.When.Time)})

	fmt.Printf("todo tracing span=%s (from %s to %s)\n", span.Context().(jaeger.SpanContext).ToString(), prev.Event.Recorder.Name, ed.Event.Recorder.Name)

	return 1
}

func (c *MetricsConfigure) Do(ed *cb.EventData) int {
	if c == nil {
		return 0
	}

	labels := DoTag(c.Attrs, ed.Event)

	switch c.Type {
	case cb.MetricType_Counter:
		vec := MetricVecStore.getCounter(c.OptsId)
		if vec == nil { // todo report error
			fmt.Println("counter vec not found for Opts", c.OptsId)

			return 0
		}

		fmt.Println(c)
		vec.With(labels).Inc()

		if len(labels) != 0 {
			fmt.Printf("todo metrics Counter(\"%s\", %s)\n", c.Name, labels)
		} else {
			fmt.Printf("todo metrics Counter(\"%s\")\n", c.Name)
		}
	case cb.MetricType_Gauge:
		fmt.Println("todo MetricsConfigure do", c.Name)
	case cb.MetricType_Histogram:
		if prev := ed.GetPreviousEventData(c.PrevName); prev != nil {
			vec := MetricVecStore.getHistogram(c.OptsId)
			if vec == nil { // todo report error
				fmt.Println("histogram vec not found for Opts", c.OptsId)

				return 0
			}

			vec.With(labels).Observe(float64(ed.Event.When.Time-prev.Event.When.Time) / float64(time.Millisecond))

			if len(labels) != 0 {
				fmt.Printf("todo metrics Histogram(\"%s\", %s)=%d (from %s to %s)\n",
					c.Name, labels, ed.Event.When.Time-prev.Event.When.Time, prev.Event.Recorder.Name, ed.Event.Recorder.Name)
			} else {
				fmt.Printf("todo metrics Histogram(\"%s\", {})=%d (from %s to %s)\n",
					c.Name, ed.Event.When.Time-prev.Event.When.Time, prev.Event.Recorder.Name, ed.Event.Recorder.Name)
			}
		} else {
			fmt.Println("previous event not found", c.PrevName)
		}
	case cb.MetricType_Summary:
		fmt.Println("todo MetricsConfigure do", c.Name)
	}

	return 1
}

func DoTraceTag(cfg []*cb.AttributeConfigure, er *cb.EventRepresentation) map[string]interface{} {
	tags := map[string]interface{}{}

	for _, path := range cfg {
		str, err := er.What.GetValue(path.Path)

		if err == nil {
			tags[path.Name] = str
		}
	}

	return tags
}

func DoTag(cfg []*cb.AttributeConfigure, er *cb.EventRepresentation) map[string]string {
	tags := map[string]string{}

	for _, path := range cfg {
		str, err := er.What.GetValue(path.Path)

		if err == nil {
			tags[path.Name] = str
		}
	}

	return tags
}

func DoTagFaster(dst []byte, cfg []*cb.AttributeConfigure, er *cb.EventRepresentation) []byte {
	dst = helper.JSONEncoder.BeginObject(dst)

	for _, path := range cfg {
		str, err := er.What.GetValue(path.Path)

		if err == nil {
			dst = helper.JSONEncoder.AppendKey(dst, path.Name)
			dst = helper.JSONEncoder.AppendString(dst, str)
		}
	}

	return dst
}
