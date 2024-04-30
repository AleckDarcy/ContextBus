package configure

import (
	cb "github.com/AleckDarcy/ContextBus/proto"

	"testing"
)

func TestDefault(t *testing.T) {
	cfg := &cb.Configure{
		Reactions: map[string]*cb.ReactionConfigure{
			"test": {
				Type: 1,
				Params: &cb.ReactionConfigure_FaultDelay{
					FaultDelay: &cb.FaultDelayParam{Ms: 999},
				},
				PreTree: &cb.PrerequisiteTree{
					LeafIDs: []int64{1, 2, 3, 4, 5},
				},
			},
		},
	}

	t.Log(cfg)
	t.Log(Store.defaultConfigure)

	Store.SetDefault(cfg)
	t.Log(Store.defaultConfigure)
}
