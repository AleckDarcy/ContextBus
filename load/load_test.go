package load

import (
	"testing"
	"time"
)

func TestSearchHotel(t *testing.T) {
	total := 8192
	speed := 2400
	speedC := 6000 // constant speed

	tasks := []*taskSetting{
		//{total: total, thread: 128, speed: speed},
		//{total: total, thread: 64, speed: speed},
		//{total: total, thread: 32, speed: speed},
		//{total: total, thread: 16, speed: speed},
		//{total: total, thread: 8, speed: speed},
		//{total: total, thread: 4, speed: speed / 2},
		//{total: total, thread: 2, speed: speed / 4},
		//{total: total, thread: 1, speed: speed / 8},

		// constant speed
		{total: total, thread: 128, speed: speedC},
		{total: total, thread: 64, speed: speedC},
		{total: total, thread: 32, speed: speedC},
		{total: total, thread: 16, speed: speedC},
		{total: total, thread: 8, speed: speedC},
		{total: total, thread: 4, speed: speedC},
		{total: total, thread: 2, speed: speedC},
		{total: total, thread: 1, speed: speedC},
	}

	run(searchHotel, tasks)

	_, _ = speed, speedC
}

func TestRecommend(t *testing.T) {
	total := 8192
	speed := 6600
	speedC := 6600 // constant speed

	tasks := []*taskSetting{
		//{total: total, thread: 128, speed: speed},
		//{total: total, thread: 64, speed: speed},
		//{total: total, thread: 32, speed: speed},
		//{total: total, thread: 16, speed: speed},
		//{total: total, thread: 8, speed: speed / 2},
		//{total: total, thread: 4, speed: speed / 3},
		//{total: total, thread: 2, speed: speed / 5},
		//{total: total, thread: 1, speed: speed / 8},

		// constant speed
		{total: total, thread: 128, speed: speedC},
		{total: total, thread: 64, speed: speedC},
		{total: total, thread: 32, speed: speedC},
		{total: total, thread: 16, speed: speedC},
		{total: total, thread: 8, speed: speedC},
		{total: total, thread: 4, speed: speedC},
		{total: total, thread: 2, speed: speedC},
		{total: total, thread: 1, speed: speedC},
	}
	sleep := 5 * time.Second

	for i, task := range tasks {
		recommend(false, task.total, task.thread, task.speed)

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
		//{total: total, thread: 128, speed: speed},
		//{total: total, thread: 64, speed: speed},
		//{total: total, thread: 32, speed: speed},
		//{total: total, thread: 16, speed: speed},
		//{total: total, thread: 8, speed: speed / 2},
		//{total: total, thread: 4, speed: speed / 3},
		//{total: total, thread: 2, speed: speed / 5},
		//{total: total, thread: 1, speed: speed / 8},

		// constant speed
		//{total: total, thread: 128, speed: speedC},
		//{total: total, thread: 64, speed: speedC},
		//{total: total, thread: 32, speed: speedC},
		//{total: total, thread: 16, speed: speedC},
		//{total: total, thread: 8, speed: speedC},
		//{total: total, thread: 4, speed: speedC},
		//{total: total, thread: 2, speed: speedC},
		{total: total, thread: 1, speed: speedC},
	}
	sleep := 5 * time.Second

	for i, task := range tasks {
		reserve(false, task.total, task.thread, task.speed)

		if i != len(tasks)-1 {
			time.Sleep(sleep)
		}
	}

	_, _ = speed, speedC
}
