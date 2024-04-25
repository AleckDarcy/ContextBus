package ContextBus

import (
	"fmt"
	"github.com/AleckDarcy/ContextBus/background"
	"github.com/AleckDarcy/ContextBus/configure"
	"github.com/AleckDarcy/ContextBus/configure/reaction"
	"github.com/AleckDarcy/ContextBus/context"
	cb "github.com/AleckDarcy/ContextBus/proto"
	"os"
	"time"
)

var HOSTNAME = os.Getenv("HOSTNAME")
var GOLANG_VERSION = os.Getenv("GOLANG_VERSION")

type signal struct {
	environmentalProfile chan struct{}
}

func init() {
	if len(HOSTNAME) == 0 {
		fmt.Println("lookup HOSTNAME from env fail")
		return
	} else if len(GOLANG_VERSION) == 0 {
		fmt.Println("lookup GOLANG_VERSION from env fail")
		return
	}

	sig := &signal{
		environmentalProfile: make(chan struct{}),
	}

	go background.EnvironmentProfiler.Run(sig.environmentalProfile)

	fmt.Printf("Initialize ContextBus(HOSTNAME=%s, GOLANG_VERSION=%s)\n", HOSTNAME, GOLANG_VERSION)
}

// OnSubmission user interface
func OnSubmission(ctx *context.Context, where *cb.EventWhere, who *cb.EventRecorder, app *cb.EventMessage) {
	er := &cb.EventRepresentation{
		When:     &cb.EventWhen{Time: time.Now().UnixNano()},
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
		Id:  background.ObservationBus.NewEventID(),
		Pcp: nil,
		Esp: esp.Timestamp,
	}

	ed := &cb.EventData{
		Event:    er,
		Metadata: md,
	}

	cfg := configure.ConfigureStore.GetConfigure(reqCtx.GetConfigureID())
	snapshots := cfg.UpdateSnapshots(who.GetName(), eveCtx.GetPrerequisiteSnapshots())
	offset := eveCtx.GetOffsetSnapshots()
	if offset != nil {
		offset = cfg.UpdateSnapshots(who.GetName(), offset)
	}

	if obs := cfg.GetObservationConfigure(who.GetName()); obs != nil {
		// todo check stacktrace configure
		// update EventMetadata
		switch obs.Type {
		case cb.ObservationType_ObservationSingle:
			// by pass PrevEvent

		case cb.ObservationType_ObservationStart:
			// initialize event pair
			newEveCtx := new(context.EventContext).SetPrerequisiteSnapshots(snapshots).SetOffsetSnapshots(offset).SetPrevEvent(eveCtx, ed)
			ctx.SetEventContext(newEveCtx)
		case cb.ObservationType_ObservationInter:
			newEveCtx := new(context.EventContext).SetPrerequisiteSnapshots(snapshots).SetOffsetSnapshots(offset).SetPrevEvent(eveCtx, ed)
			ctx.SetEventContext(newEveCtx)
			fallthrough
		case cb.ObservationType_ObservationEnd:
			// finalize event pair
			_, prevED := eveCtx.GetPrevEvent()
			if prevED == nil {
				fmt.Errorf("eveCtx.GetPrevEvent() get nil")
			}
			ed.PrevEventData = prevED
		}

		obs.Prepare(ctx, ed)
		if ed.SpanMetadata != nil {
			ctx.GetRequestContext().SetSpanMetadata(ed.SpanMetadata)
		}
	}

	if rac := cfg.GetReaction(who.Name); rac != nil {
		if snapshot := snapshots.GetPrerequisiteSnapshot(who.Name); snapshot != nil {
			if ok, err := rac.PreTree.Check((*reaction.PrerequisiteSnapshot)(snapshot)); err != nil {

			} else if !ok {

			} else {
				fmt.Println("prerequisites accomplished")

				switch rac.Type {
				case cb.ReactionType_FaultDelay:
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
