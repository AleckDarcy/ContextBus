package ContextBus

import (
	"github.com/AleckDarcy/ContextBus/background"
	"github.com/AleckDarcy/ContextBus/configure"
	"github.com/AleckDarcy/ContextBus/context"
	"github.com/AleckDarcy/ContextBus/helper"
	cb "github.com/AleckDarcy/ContextBus/proto"

	"sync"
	"testing"
	"time"
)

var cfg2 = &cb.Configure{
	Reactions: nil,
	Observations: map[string]*cb.ObservationConfigure{
		"EventA": {
			Logging: &cb.LoggingConfigure{
				Timestamp: &cb.TimestampConfigure{Format: helper.TIME_FORMAT_RFC3339Nano},
				Attrs:     []*cb.AttributeConfigure{cb.Test_AttributeConfigure_Rest_Key},
				Out:       cb.LogOutType_LogOutType_, // omit print
			},
		},
	},
}

var cfg3 = &cb.Configure{
	Reactions: map[string]*cb.ReactionConfigure{
		"EventD": {
			Type: cb.ReactionType_FaultCrash,
			PreTree: &cb.PrerequisiteTree{
				Nodes: []*cb.PrerequisiteNode{
					cb.NewPrerequisiteMessageNode(0, "EventA", &cb.ConditionTree{}, -1, nil),
				},
			}},
		"EventC": {
			Type: cb.ReactionType_FaultCrash,
			PreTree: &cb.PrerequisiteTree{
				Nodes: []*cb.PrerequisiteNode{
					cb.NewPrerequisiteLogicNode(0, cb.LogicType_And_, -1, []int64{1, 2}),
					cb.NewPrerequisiteMessageNode(1, "EventA",
						cb.NewConditionTree([]*cb.ConditionNode{cb.Test_Condition_C_1_0}, nil), 0, nil),
					cb.NewPrerequisiteMessageNode(2, "EventB",
						cb.NewConditionTree([]*cb.ConditionNode{cb.Test_Condition_C_2_0}, nil), 0, nil),
				},
			},
		},
	},
}

func TestObservationBus_Observation(t *testing.T) {
	go background.ObservationBus.Run(nil)

	id := int64(2)
	configure.ConfigureStore.SetConfigure(id, cfg2)

	ctx := context.NewContext(context.NewRequestContext("rest", id, rest), context.NewEventContext(nil, nil))
	app := new(cb.EventMessage).SetMessage("received message from %s").SetPaths([]*cb.Path{path})

	n := 10
	wg := sync.WaitGroup{}
	worker := func() {
		for i := 0; i < n; i++ {
			OnSubmission(ctx, &cb.EventWhere{}, &cb.EventRecorder{Name: "EventA"}, app)

			time.Sleep(time.Millisecond * 500)
		}
		wg.Done()
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker()
	}

	wg.Wait()

	time.Sleep(time.Second * 2)
}

func TestObservationBus_Reaction(t *testing.T) {
	go background.ObservationBus.Run(nil)

	id := int64(3)
	configure.ConfigureStore.SetConfigure(id, cfg3)

	cfg := configure.ConfigureStore.GetConfigure(id)
	t.Log(cfg.ReactionIndex)

	ss := cfg.InitializeSnapshots()
	t.Log(ss)

	cfg.UpdateSnapshots("EventA", ss)
	t.Log(ss)

	cfg.UpdateSnapshots("EventB", ss)
	t.Log(ss)

	cfg.UpdateSnapshots("EventC", ss)
	t.Log(ss)

	reqCtx := context.NewRequestContext("rest", id, rest)
	eveCtx := context.NewEventContext(nil, &cb.PrerequisiteSnapshots{
		Snapshots: map[string]*cb.PrerequisiteSnapshot{
			"EventD": {Value: []int64{0}},
			"EventC": {Value: []int64{0, 0, 0}},
		},
	})
	ctx := context.NewContext(reqCtx, eveCtx)
	app := new(cb.EventMessage).SetMessage("received message from %s").SetPaths([]*cb.Path{path})

	n := 1
	for i := 0; i < n; i++ {
		OnSubmission(ctx, &cb.EventWhere{}, &cb.EventRecorder{Name: "EventA"}, app)

		// initialize PrerequisiteSnapshot for each submission

		OnSubmission(ctx, &cb.EventWhere{}, &cb.EventRecorder{Name: "EventC"}, app)

		time.Sleep(time.Millisecond * 500)
	}

	time.Sleep(time.Second * 2)
}
