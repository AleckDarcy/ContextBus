package background

import (
	"github.com/AleckDarcy/ContextBus/helper"
	cb "github.com/AleckDarcy/ContextBus/proto"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"

	"fmt"
	"runtime"
	"sync"
	"time"
)

var MEMSTATS = &runtime.MemStats{}

type environmentProfiler struct {
	latest *cb.EnvironmentalProfile
	lock   sync.RWMutex
	store  map[int64]*cb.EnvironmentalProfile
}

var EnvironmentProfiler = &environmentProfiler{
	latest: &cb.EnvironmentalProfile{
		Hardware: &cb.HardwareProfile{},
	},
	store: map[int64]*cb.EnvironmentalProfile{},
}

func (e *environmentProfiler) GetLatest() *cb.EnvironmentalProfile {
	e.lock.RLock()
	latest := e.latest
	e.lock.RUnlock()

	return latest
}

func (e *environmentProfiler) GetByID(id int64) *cb.EnvironmentalProfile {
	e.lock.RLock()
	ep := e.store[id]
	e.lock.RUnlock()

	return ep
}

var npPrev = EnvironmentProfiler.GetNetProfile()

func (e *environmentProfiler) GetNetProfile() *cb.NetProfile {
	n, err := net.IOCounters(false)

	if err != nil {
		fmt.Println("NetProfile", err)

		return nil
	}

	return &cb.NetProfile{
		BytesSent:   n[0].BytesSent,
		BytesRecv:   n[0].BytesRecv,
		PacketsSent: n[0].PacketsSent,
		PacketsRecv: n[0].PacketsRecv,
		Errin:       n[0].Errin,
		Errout:      n[0].Errout,
		Dropin:      n[0].Dropin,
		Dropout:     n[0].Dropout,
	}
}

func (e *environmentProfiler) GetEnvironmentProfile() *cb.EnvironmentalProfile {
	ep := &cb.EnvironmentalProfile{
		Timestamp: time.Now().UnixNano(),
		Hardware:  &cb.HardwareProfile{},
	}

	signal := make(chan *cb.CPUProfile)
	go func(signal chan *cb.CPUProfile) {
		if c, err := cpu.Percent(helper.CPU_PROFILE_DURATION, false); err == nil {
			signal <- &cb.CPUProfile{
				Percent: c[0],
			}
		} else {
			signal <- nil
		}
	}(signal)

	if m, err := mem.VirtualMemory(); err == nil {
		ep.Hardware.Mem = &cb.MemProfile{
			Total:       m.Total,
			Available:   m.Available,
			Used:        m.Used,
			UsedPercent: m.UsedPercent,
			Free:        m.Free,
		}
	}

	np := e.GetNetProfile()
	ep.Hardware.Net = &cb.NetProfile{
		BytesSent:   np.BytesSent - npPrev.BytesSent,
		BytesRecv:   np.BytesRecv - npPrev.BytesRecv,
		PacketsSent: np.PacketsSent - npPrev.PacketsSent,
		PacketsRecv: np.PacketsRecv - npPrev.PacketsRecv,
		Errin:       np.Errin - npPrev.Errin,
		Errout:      np.Errout - npPrev.Errout,
		Dropin:      np.Dropin - npPrev.Dropin,
		Dropout:     np.Dropout - npPrev.Dropout,
	}
	npPrev = np

	runtime.ReadMemStats(MEMSTATS)
	ep.Language = &cb.LanguageProfile{
		Type: cb.LanguageType_Golang,
		Profile: &cb.LanguageProfile_Go{
			Go: &cb.LanguageGo{
				HeapSys:       MEMSTATS.HeapSys,
				HeapAlloc:     MEMSTATS.HeapAlloc,
				HeapInuse:     MEMSTATS.HeapInuse,
				StackSys:      MEMSTATS.StackSys,
				StackInuse:    MEMSTATS.StackInuse,
				MSpanInuse:    MEMSTATS.MSpanInuse,
				MSpanSys:      MEMSTATS.MSpanSys,
				MCacheInuse:   MEMSTATS.MCacheInuse,
				MCacheSys:     MEMSTATS.MCacheSys,
				LastGC:        MEMSTATS.LastGC,
				NextGC:        MEMSTATS.NextGC,
				GCCPUFraction: MEMSTATS.GCCPUFraction,
			},
		},
	}

	select {
	case c := <-signal:
		ep.Hardware.Cpu = c
	case <-time.After(helper.CPU_PROFILE_DURATION_MAX):
		// todo: report error
	}

	e.lock.Lock()
	e.latest.Next = ep.Timestamp
	ep.Prev = e.latest.Timestamp
	e.latest = ep
	e.store[ep.Timestamp] = ep
	e.lock.Unlock()

	//fmt.Println("GetEnvironmentProfile()")

	return ep
}

func (e *environmentProfiler) Run(sig chan struct{}) {
	e.GetEnvironmentProfile()

	for {
		select {
		case <-sig:

			return
		case <-time.After(helper.ENV_PROFILE_INTERVAL):
			e.GetEnvironmentProfile()
		}
	}
}
