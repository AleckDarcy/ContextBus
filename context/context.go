package context

import (
	"testing"
	"time"

	cb "github.com/AleckDarcy/ContextBus/proto"
	"github.com/AleckDarcy/ContextBus/third-party/github.com/opentracing/opentracing-go"
)

const CB_CONTEXT_NAME = "context_bus"

// RequestContext is inter-service message context
type RequestContext struct {
	lib         string // name of network API
	requestID   uint64
	configureID int64
	event       *cb.EventMessage
	span        *cb.SpanMetadata // parent span
}

func NewRequestContext(lib string, requestID uint64, configureID int64, msg *cb.EventMessage) *RequestContext {
	return &RequestContext{
		lib:         lib,
		requestID:   requestID,
		configureID: configureID,
		event:       msg,
	}
}

func (c *RequestContext) GetLib() string {
	return c.lib
}

func (c *RequestContext) GetRequestID() uint64 {
	return c.requestID
}

func (c *RequestContext) GetConfigureID() int64 {
	return c.configureID
}

func (c *RequestContext) GetEventMessage() *cb.EventMessage {
	return c.event
}

func (c *RequestContext) SetSpanMetadata(sm *cb.SpanMetadata) *RequestContext {
	c.span = sm

	return c
}

func (c *RequestContext) GetSpanMetadata() *cb.SpanMetadata {
	return c.span
}

// EventContext is the context associated with each observation
type EventContext struct {
	codebase         *cb.CodeBaseInfo
	snapshots        *cb.PrerequisiteSnapshots
	offsetSnapshots  *cb.PrerequisiteSnapshots
	prevEventContext *EventContext // todo event-id
	prevEventData    *cb.EventData // todo event-id

	// todo uuid for inter-service communication
}

func NewEventContext(codebase *cb.CodeBaseInfo, snapshots *cb.PrerequisiteSnapshots) *EventContext {
	return &EventContext{
		codebase:  codebase,
		snapshots: snapshots,
	}
}

func (c *EventContext) SetCodeInfoBasic(codebase *cb.CodeBaseInfo) *EventContext {
	c.codebase = codebase

	return c
}

func (c *EventContext) GetCodeInfoBasic() *cb.CodeBaseInfo {
	return c.codebase
}

func (c *EventContext) SetPrerequisiteSnapshots(snapshots *cb.PrerequisiteSnapshots) *EventContext {
	c.snapshots = snapshots

	return c
}

func (c *EventContext) GetPrerequisiteSnapshots() *cb.PrerequisiteSnapshots {
	return c.snapshots
}

func (c *EventContext) SetOffsetSnapshots(snapshots *cb.PrerequisiteSnapshots) *EventContext {
	c.offsetSnapshots = snapshots

	return c
}

func (c *EventContext) GetOffsetSnapshots() *cb.PrerequisiteSnapshots {
	return c.offsetSnapshots
}

func (c *EventContext) SetPrevEvent(eveCtx *EventContext, ed *cb.EventData) *EventContext {
	c.prevEventContext = eveCtx
	c.prevEventData = ed

	return c
}

func (c *EventContext) GetPrevEvent() (*EventContext, *cb.EventData) {
	return c.prevEventContext, c.prevEventData
}

type Context struct {
	reqCtx *RequestContext
	eveCtx *EventContext

	tracer opentracing.Tracer
	span   *cb.SpanMetadata

	timestamp int64
}

func NewContext(reqCtx *RequestContext, eveCtx *EventContext) *Context {
	return &Context{
		reqCtx: reqCtx,
		eveCtx: eveCtx,
	}
}

// Payload generates Payload for grpc messages
// todo: fields
func (c *Context) Payload() *cb.Payload {
	if c == nil {
		return nil
	}

	return &cb.Payload{
		RequestId: c.reqCtx.requestID,
		ConfigId:  c.reqCtx.configureID,
		Snapshots: c.eveCtx.snapshots,
		Addition:  nil,
		Parent:    c.span,
		MType:     cb.MessageType_Message_Request,
		Uuid:      "",
	}
}

func (c *Context) SetTimestamp() {
	c.timestamp = time.Now().UnixNano()
}

func (c *Context) GetTimestamp() int64 {
	return c.timestamp
}

func (c *Context) GetRequestContext() *RequestContext {
	return c.reqCtx
}

// SetRequestContext is written by network APIs on receiving requests. e.g., rest, rpc
// contains the request-wise static values. e.g., session-id, token
func (c *Context) SetRequestContext(reqCtx *RequestContext) *Context {
	c.reqCtx = reqCtx

	return c
}

func (c *Context) GetEventContext() *EventContext {
	return c.eveCtx
}

func (c *Context) PrintPrevEventData(t *testing.T) {
	prev := c.eveCtx.prevEventData
	for prev != nil {
		t.Logf("%+v", prev.Event.Recorder.Name)
		prev = prev.PrevEventData
	}
	t.Log()
}

// SetEventContext is written by generated code when user or library submit their observations
// contains the static code base information
func (c *Context) SetEventContext(eveCtx *EventContext) *Context {
	c.eveCtx = eveCtx

	return c
}

func (c *Context) SetTracer(tracer opentracing.Tracer) *Context {
	c.tracer = tracer

	return c
}

func (c *Context) GetTracer() opentracing.Tracer {
	return c.tracer
}

func (c *Context) SetSpanMetadata(sm *cb.SpanMetadata) *Context {
	c.span = sm

	return c
}

func (c *Context) GetSpanMetadata() *cb.SpanMetadata {
	return c.span
}
