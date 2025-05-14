package ContextBus

import (
	"encoding/json"
	"github.com/AleckDarcy/ContextBus/proto"
	"io/ioutil"
	"os"
	"sort"
	"testing"
)

func static(nums []float64) (min, max, avg, med float64) {
	total := float64(len(nums))

	if total == 0 {
		return
	}

	sort.Float64s(nums)

	for _, num := range nums {
		avg += num / total
	}

	return nums[0], nums[len(nums)-1], avg, nums[len(nums)/2]
}

func TestName(t *testing.T) {
	f, _ := os.Open("metric.txt")
	b, _ := ioutil.ReadAll(f)

	m := &proto.Payload{}
	json.Unmarshal(b, &m)

	l_c, lt_c := []float64{}, []float64{}
	l_p, lt_p := []float64{}, []float64{}
	for _, l := range m.Metric.CBLatency.Latency {
		if l.Type == proto.CBType_Logging {
			l_c = append(l_c, l.Channel)
			l_p = append(l_p, l.Process)
		} else if l.Type == proto.CBType_Logging|proto.CBType_Tracing {
			lt_c = append(lt_c, l.Channel)
			lt_p = append(lt_p, l.Process)
		}
	}

	t.Log(static(l_c))
	t.Log(static(l_p))
	t.Log(static(lt_c))
	t.Log(static(lt_p))
}
