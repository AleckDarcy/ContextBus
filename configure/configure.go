package configure

import (
	"github.com/AleckDarcy/ContextBus/configure/observation"
	"github.com/AleckDarcy/ContextBus/configure/reaction"
	"github.com/AleckDarcy/ContextBus/helper"
	cb "github.com/AleckDarcy/ContextBus/proto"

	"sync"
	"sync/atomic"
	"unsafe"
)

type store struct {
	lock sync.RWMutex

	defaultConfigure *Configure
	configures       map[int64]*Configure // int64: configure_id
}

var Store = store{configures: map[int64]*Configure{}}

// DefaultObservation converts event data to a single log entry
var DefaultObservation = &cb.ObservationConfigure{
	Type: cb.ObservationType_ObservationSingle,
	Logging: &cb.LoggingConfigure{
		Timestamp: &cb.TimestampConfigure{Format: helper.TIME_FORMAT_DEFAULT},
		Out:       cb.LogOutType_Stdout,
	},
}

func (s *store) convertConfigure(cfg *cb.Configure) *Configure {
	racs := map[string]*reaction.Configure(nil)
	racMapMap := map[string]map[string]*reaction.Configure{} // <event (observation), <event (reaction), cfg>>

	if reactions := cfg.Reactions; reactions != nil {
		racs = make(map[string]*reaction.Configure, len(reactions))
		for name, reaction_ := range reactions {
			rac := &reaction.Configure{
				Name:    name,
				Type:    reaction_.Type,
				Params:  reaction_.Params,
				PreTree: reaction.NewPrerequisiteTree(reaction_.PreTree),
			}
			racs[name] = rac

			for _, node := range reaction_.PreTree.Nodes {
				if node.Type == cb.PrerequisiteNodeType_PrerequisiteMessage_ {
					racMap, ok := racMapMap[node.Message.Name]
					if !ok {
						racMap = map[string]*reaction.Configure{}
						racMapMap[node.Message.Name] = racMap
					}
					racMap[name] = rac
				}
			}
		}
	}

	racIndex := map[string][]*reaction.Configure{}
	for name, racMap := range racMapMap {
		racList := make([]*reaction.Configure, 0, len(racMap))
		for _, rac := range racMap {
			racList = append(racList, rac)
		}

		racIndex[name] = racList
	}

	return &Configure{
		Reactions:     racs,
		Observations:  cfg.Observations,
		ReactionIndex: racIndex,
	}
}

func (s *store) SetDefault(configure *cb.Configure) {
	cfg := s.convertConfigure(configure)

	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&s.defaultConfigure)), unsafe.Pointer(cfg))
}

func (s *store) GetDefault() *cb.Configure {
	return (*cb.Configure)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s.defaultConfigure))))
}

func (s *store) SetConfigure(id int64, configure *cb.Configure) {
	cfg := s.convertConfigure(configure)

	s.lock.Lock()
	s.configures[id] = cfg
	s.lock.Unlock()
}

func (s *store) GetConfigure(id int64) *Configure {
	s.lock.RLock()
	cfg := s.configures[id]
	s.lock.RUnlock()

	if cfg != nil {
		return cfg
	}

	// return default Configure
	return &Configure{}
}

func (c *Configure) InitializeSnapshots() *cb.PrerequisiteSnapshots {
	ss := make(map[string]*cb.PrerequisiteSnapshot, len(c.Reactions))

	for name, rac := range c.Reactions {
		ss[name] = (*cb.PrerequisiteSnapshot)(rac.PreTree.InitializeSnapshot())
	}

	return &cb.PrerequisiteSnapshots{Snapshots: ss}
}

func (c *Configure) UpdateSnapshots(name string, ss *cb.PrerequisiteSnapshots) *cb.PrerequisiteSnapshots {
	racs, ok := c.ReactionIndex[name]
	if ok {
		for _, rac := range racs {
			rac.PreTree.UpdateSnapshot(name, (*reaction.PrerequisiteSnapshot)(ss.Snapshots[rac.Name]))
		}
	}

	return ss
}

func (c *Configure) UpdateBothSnapshots(name string, ss, offset *cb.PrerequisiteSnapshots) (*cb.PrerequisiteSnapshots, *cb.PrerequisiteSnapshots) {
	racs, ok := c.ReactionIndex[name]
	if ok {
		for _, rac := range racs {
			rac.PreTree.UpdateSnapshot(name, (*reaction.PrerequisiteSnapshot)(ss.Snapshots[rac.Name]))
			rac.PreTree.UpdateSnapshot(name, (*reaction.PrerequisiteSnapshot)(offset.Snapshots[rac.Name]))
		}
	}

	return ss, offset
}

func (c *Configure) GetObservationConfigure(name string) *observation.Configure {
	if c.Observations == nil {
		return (*observation.Configure)(DefaultObservation)
	}

	return (*observation.Configure)(c.Observations[name])
}

func (c *Configure) GetReaction(name string) *reaction.Configure {
	if c.Reactions == nil {
		return nil
	}

	return c.Reactions[name]
}
