package load

import (
	"testing"
	"time"

	"github.com/AleckDarcy/ContextBus/configure"
)

func TestSearchHotel(t *testing.T) {
	total := 8192
	speed := 3000
	cbcID := configure.CBCID_OBSERVATIONBYPASS
	cbcTraceRatio := 0

	tasks := []*taskSetting{
		//{total: total / 8, threads: 128, speed: speed, cbcID: configure.CBCID_LOGGINGYPASS, cbTraceRatio: 0}, // warm up

		// perf metric on
		{total: total, threads: 256, speed: speed, cbcID: cbcID, cbTraceRatio: cbcTraceRatio},
		//{total: total, threads: 128, speed: speed, cbcID: cbcID, cbTraceRatio: cbcTraceRatio},
		//{total: total, threads: 128, speed: speed, cbcID: cbcID, cbTraceRatio: cbcTraceRatio},
		//{total: total, threads: 128, speed: speed, cbcID: cbcID, cbTraceRatio: cbcTraceRatio},

		//{total: total, threads: 128, speed: speed, cbcID: cbcID, cbTraceRatio: cbcTraceRatio},
		//{total: total, threads: 64, speed: speed, cbcID: cbcID, cbTraceRatio: cbcTraceRatio},
		//{total: total, threads: 32, speed: speed, cbcID: cbcID, cbTraceRatio: cbcTraceRatio},
		//{total: total, threads: 16, speed: speed, cbcID: cbcID, cbTraceRatio: cbcTraceRatio},
		//{total: total, threads: 8, speed: speed, cbcID: cbcID, cbTraceRatio: cbcTraceRatio},
		//{total: total, threads: 4, speed: speed / 2, cbcID: cbcID, cbTraceRatio: cbcTraceRatio},
		//{total: total, threads: 2, speed: speed / 4, cbcID: cbcID, cbTraceRatio: cbcTraceRatio},
		//{total: total, threads: 1, speed: speed / 8, cbcID: cbcID, cbTraceRatio: cbcTraceRatio},

		//{total: total, threads: 128, speed: speed, cbcID: cbcID, cbTraceRatio: 0},
		//{total: total, threads: 128, speed: speed, cbcID: cbcID, cbTraceRatio: 1},
		//{total: total, threads: 128, speed: speed, cbcID: cbcID, cbTraceRatio: 5},
		//{total: total, threads: 128, speed: speed, cbcID: cbcID, cbTraceRatio: 10},
		//{total: total, threads: 128, speed: speed, cbcID: cbcID, cbTraceRatio: 20},
		//{total: total, threads: 128, speed: speed, cbcID: cbcID, cbTraceRatio: 50},
		//{total: total, threads: 128, speed: speed, cbcID: cbcID, cbTraceRatio: 100},
	}

	run(searchHotel, tasks)

	_ = cbcID
	_ = speed
	_ = cbcTraceRatio
}

func TestRecommend(t *testing.T) {
	total := 8192
	speed := 6600
	speedC := 6600 // constant speed

	tasks := []*taskSetting{
		//{total: total, threads: 128, speed: speed},
		//{total: total, threads: 64, speed: speed},
		//{total: total, threads: 32, speed: speed},
		//{total: total, threads: 16, speed: speed},
		//{total: total, threads: 8, speed: speed / 2},
		//{total: total, threads: 4, speed: speed / 3},
		//{total: total, threads: 2, speed: speed / 5},
		//{total: total, threads: 1, speed: speed / 8},

		// constant speed
		{total: total, threads: 128, speed: speedC},
		{total: total, threads: 64, speed: speedC},
		{total: total, threads: 32, speed: speedC},
		{total: total, threads: 16, speed: speedC},
		{total: total, threads: 8, speed: speedC},
		{total: total, threads: 4, speed: speedC},
		{total: total, threads: 2, speed: speedC},
		{total: total, threads: 1, speed: speedC},
	}
	sleep := 5 * time.Second

	for i, task := range tasks {
		recommend(false, task)

		if i != len(tasks)-1 {
			time.Sleep(sleep)
		}
	}

	_, _ = speed, speedC
}

func TestReserve(t *testing.T) {
	total := 8192
	speed := 6600
	speedC := 6000 // constant speed

	tasks := []*taskSetting{
		//{total: total, threads: 128, speed: speed},
		//{total: total, threads: 64, speed: speed},
		//{total: total, threads: 32, speed: speed},
		//{total: total, threads: 16, speed: speed},
		//{total: total, threads: 8, speed: speed / 2},
		//{total: total, threads: 4, speed: speed / 3},
		//{total: total, threads: 2, speed: speed / 5},
		//{total: total, threads: 1, speed: speed / 8},

		// constant speed
		//{total: total, threads: 128, speed: speedC},
		//{total: total, threads: 64, speed: speedC},
		//{total: total, threads: 32, speed: speedC},
		//{total: total, threads: 16, speed: speedC},
		//{total: total, threads: 8, speed: speedC},
		//{total: total, threads: 4, speed: speedC},
		//{total: total, threads: 2, speed: speedC},
		{total: total, threads: 1, speed: speedC},
	}
	sleep := 5 * time.Second

	for i, task := range tasks {
		reserve(false, task)

		if i != len(tasks)-1 {
			time.Sleep(sleep)
		}
	}

	_, _ = speed, speedC
}
