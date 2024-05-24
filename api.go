package ContextBus

import (
	"context"
	"fmt"
	"time"

	"github.com/AleckDarcy/ContextBus/background"
	"github.com/AleckDarcy/ContextBus/configure"
	"github.com/AleckDarcy/ContextBus/configure/reaction"
	cb_context "github.com/AleckDarcy/ContextBus/context"
	cb "github.com/AleckDarcy/ContextBus/proto"
)

// FromContext reads Context from the http.Request.
// returns ContextBus context
func FromContext(ctx context.Context) (*cb_context.Context, bool) {
	if cbCtxItf := ctx.Value(cb_context.CB_CONTEXT_NAME); cbCtxItf != nil {
		// fmt.Printf("retrieved ContextBus context from HTTP: %+v\n", cbCtxItf)

		return cbCtxItf.(*cb_context.Context), true
	}

	return nil, false
}

// FromPayload reads Context from the proto message.
// returns Context Bus context
func FromPayload(ctx context.Context, pay *cb.Payload) (*cb_context.Context, bool) {
	if pay == nil {
		return nil, false
	} else if pay.ConfigId == configure.CBCID_BYPASS {
		return nil, false
	}

	// todo
	reqCtx := cb_context.NewRequestContext("", pay.RequestId, pay.ConfigId, nil).SetSpanMetadata(pay.Parent)
	eveCtx := cb_context.NewEventContext(nil, &cb.PrerequisiteSnapshots{})

	cbCtx := cb_context.NewContext(reqCtx, eveCtx).SetTracer(background.ObservationBus.GetTracer())

	// fmt.Printf("retrieved ContextBus context from Payload: %+v\n", cbCtx)

	return cbCtx, true
}

// OnSubmission user interface
func OnSubmission(ctx *cb_context.Context, where *cb.EventWhere, who *cb.EventRecorder, app *cb.EventMessage) {
	ctx.SetTimestamp()

	er := &cb.EventRepresentation{
		When:     &cb.EventWhen{Time: ctx.GetTimestamp()},
		Where:    where,
		Recorder: who,
		What:     &cb.EventWhat{Application: app},
	}

	reqCtx := ctx.GetRequestContext()
	eveCtx := ctx.GetEventContext()

	// write network API attributes
	if reqCtx != nil {
		er.What.WithLibrary(reqCtx.GetLib(), reqCtx.GetEventMessage())
	}
	// write code base info

	esp := background.EnvironmentProfiler.GetLatest()

	md := &cb.EventMetadata{
		ReqId: reqCtx.GetRequestID(),
		EveId: background.ObservationBus.NewEventID(),
		Pcp:   nil,
		Esp:   esp.Timestamp,
	}

	ed := &cb.EventData{
		Event:    er,
		Metadata: md,
	}

	cfg := configure.Store.GetConfigure(reqCtx.GetConfigureID())
	snapshots := cfg.UpdateSnapshots(who.GetName(), eveCtx.GetPrerequisiteSnapshots())
	offset := eveCtx.GetOffsetSnapshots()
	if offset != nil {
		offset = cfg.UpdateSnapshots(who.GetName(), offset)
	}

	if obs := cfg.GetObservationConfigure(who.GetName()); obs != nil {
		// fmt.Printf("found observation configure for %s\n", who.GetName())

		// todo check stacktrace configure
		// update EventMetadata
		switch obs.Type {
		case cb.ObservationType_ObservationSingle:
			// bypass PrevEvent

		case cb.ObservationType_ObservationStart:
			// initialize event chain
			newEveCtx := new(cb_context.EventContext).SetPrerequisiteSnapshots(snapshots).SetOffsetSnapshots(offset).SetPrevEvent(eveCtx, ed)
			ctx.SetEventContext(newEveCtx)
		case cb.ObservationType_ObservationInter:
			newEveCtx := new(cb_context.EventContext).SetPrerequisiteSnapshots(snapshots).SetOffsetSnapshots(offset).SetPrevEvent(eveCtx, ed)
			ctx.SetEventContext(newEveCtx)
			fallthrough
		case cb.ObservationType_ObservationEnd:
			// finalize event chain
			_, prevED := eveCtx.GetPrevEvent()
			if prevED == nil {
				fmt.Errorf("eveCtx.GetPrevEvent() get nil")
			}
			ed.PrevEventData = prevED
		}

		obs.Prepare(ctx, ed)
		if ed.SpanMetadata != nil {
			ctx.SetSpanMetadata(ed.SpanMetadata)
		}
	}

	if rac := cfg.GetReaction(who.Name); rac != nil {
		if snapshot := snapshots.GetPrerequisiteSnapshot(who.Name); snapshot != nil {
			if ok, err := rac.PreTree.Check((*reaction.PrerequisiteSnapshot)(snapshot)); err != nil {

			} else if !ok {

			} else {
				fmt.Println("prerequisites accomplished")

				switch rac.Type {
				case cb.ReactionType_ReactionFaultDelay:
					params := rac.Params.(*cb.ReactionConfigure_FaultDelay).FaultDelay

					time.Sleep(time.Duration(params.Ms) * time.Millisecond)
					fmt.Println("slept for", params.Ms, "ms")
				}
			}
		}
	}

	// push EventData to bus
	background.ObservationBus.OnSubmit(ctx, ed)
}
