package proto

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync/atomic"
)

const (
	Metric_S = iota

	Metric_Frontend_SearchHandler_Logic_1
	Metric_Frontend_SearchHandler_1
	Metric_Frontend_SearchHandler_Logic_2
	Metric_Frontend_SearchHandler_2
	//Metric_Frontend_SearchHandler_Logic_3
	Metric_Frontend_SearchHandler_3
	Metric_Frontend_SearchHandler_Logic_4
	Metric_Frontend_SearchHandler_4
	Metric_Frontend_SearchHandler_Logic_5
	Metric_Frontend_SearchHandler_5
	Metric_Frontend_SearchHandler_Logic_6
	Metric_Frontend_SearchHandler_6
	Metric_Frontend_SearchHandler_Logic_7

	Metric_Search_NearBy_Observation_1
	Metric_Search_NearBy_Logic_2
	Metric_Search_NearBy_Observation_2
	Metric_Search_NearBy_Logic_3
	Metric_Search_NearBy_Observation_3
	Metric_Search_NearBy_Logic_4

	Metric_E

	CBMetric_Len = Metric_E * 8192
)

var PERF_METRIC = false

var MetricSize = 8192

var PMetric *PerfMetric

func init() {
	tmpInt, err := strconv.Atoi(os.Getenv("CB_PERF_METRIC"))
	if err != nil {
		fmt.Println("lookup PERF_METRIC from env fail:", err)
	} else {
		PERF_METRIC = tmpInt == 1
	}
}

func NewPerfMetric(size int, from, to int) *PerfMetric {
	if !PERF_METRIC {
		return nil
	}

	metric := &PerfMetric{
		CBLatency: &CBLatencyMetric{
			Total:   0,
			Latency: make([]*CBLatency, CBMetric_Len),
		},
		Latency: make([]*LatencyMetric, Metric_E),
	}

	for i := from; i <= to; i++ {
		metric.Latency[i] = &LatencyMetric{
			Total:   0,
			Latency: make([]float64, size),
			Mean:    0,
			Median:  0,
		}
	}

	return metric
}

const (
	CBType_Logging = int64(1) << iota
	CBType_Tracing
	CBType_Metrics
)

func (m *PerfMetric) AddCBLatency(cbType int64, enQueue, deQueue, end int64) {
	if !PERF_METRIC {
		return
	}

	if m == nil || m.CBLatency == nil {
		return
	}

	cbLat := m.CBLatency
	// todo multi-thread
	lat := &CBLatency{
		Type:    cbType,
		Channel: float64(deQueue - enQueue),
		Process: float64(end - deQueue),
	}

	cbLat.Latency = append(cbLat.Latency, lat)
	cbLat.Total++

	if cbLat.FirstEnqueue == 0 {
		cbLat.FirstEnqueue = enQueue
	}

	cbLat.LastFinished = end
}

func (m *PerfMetric) AddLatency(idx int, value float64) {
	if !PERF_METRIC {
		return
	}

	if m == nil || len(m.Latency) <= idx {
		return
	}

	i := atomic.AddInt64(&m.Latency[idx].Total, 1) - 1
	if i >= int64(len(m.Latency[idx].Latency)) {
		return
	}

	m.Latency[idx].Latency[i] = value
}

func (m *PerfMetric) Calculate() *PerfMetric {
	if !PERF_METRIC {
		return nil
	}

	res := &PerfMetric{
		CBLatency: &CBLatencyMetric{
			Total:        m.CBLatency.Total,
			Latency:      make([]*CBLatency, m.CBLatency.Total),
			FirstEnqueue: m.CBLatency.FirstEnqueue,
			LastFinished: m.CBLatency.LastFinished,
		},
		Latency: make([]*LatencyMetric, Metric_E),
	}

	copy(res.CBLatency.Latency, m.CBLatency.Latency)
	m.CBLatency.Total = 0
	m.CBLatency.Latency = m.CBLatency.Latency[0:0]

	for i := 0; i < Metric_E; i++ {
		srcLatI := m.Latency[i]
		if srcLatI == nil {
			continue
		}

		if total := srcLatI.Total; total != 0 {
			dstLatI := &LatencyMetric{}
			res.Latency[i] = dstLatI

			if total > int64(len(srcLatI.Latency)) {
				total = int64(len(srcLatI.Latency))
			}

			dstLatI.Total = total
			//fmt.Println("calculate latency metric (idx total):", i, total)

			sort.Float64s(srcLatI.Latency[:total])

			sum := float64(0)
			for j := 0; j < int(total); j++ {
				sum += srcLatI.Latency[j]
			}
			dstLatI.Mean = sum / float64(total)

			if total%2 == 0 {
				dstLatI.Median = (srcLatI.Latency[total/2-1] + srcLatI.Latency[total/2]) / 2
			} else {
				dstLatI.Median = srcLatI.Latency[total/2]
			}

			dstLatI.Min = srcLatI.Latency[0]
			dstLatI.Max = srcLatI.Latency[total-1]
		}

		srcLatI.Total = 0
	}

	return res
}

func (m *PerfMetric) Merge(src *PerfMetric) {
	if !PERF_METRIC {
		return
	}

	if src == nil {
		return
	}

	for i, srcI := range src.Latency {
		if srcI != nil && srcI.Total != 0 { // serialization makes nils to {}'s
			m.Latency[i] = &LatencyMetric{
				Total:  srcI.Total,
				Mean:   srcI.Mean,
				Median: srcI.Median,
				Min:    srcI.Min,
				Max:    srcI.Max,
			}
		}
	}
}
