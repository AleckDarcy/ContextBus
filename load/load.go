package load

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/AleckDarcy/ContextBus/configure"
	cb "github.com/AleckDarcy/ContextBus/proto"
)

const URL = "http://localhost:5001"

type user struct {
	id       string
	password string
}

var users = make([]*user, 500)

func init() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 500; i++ {
		name := fmt.Sprintf("User_%x", strconv.Itoa(i))

		suffix := strconv.Itoa(i)
		password := ""
		for j := 0; j < 10; j++ {
			password += suffix
		}

		users[i] = &user{
			id:       name,
			password: password,
		}
	}
}

type param struct {
	method string
	path   string
}

type params struct {
	size   int
	params []*param
}

func (p *params) get(i int) *param {
	return p.params[i%p.size]
}

type request struct {
	id   int
	para *param
}

type taskSetting struct {
	total        int
	threads      int
	speed        int
	cbcID        int64
	cbTraceRatio int
	cbTraceCount int
}

type result struct {
	latency int64
	err     bool
}

func getUser() *user {
	return users[rand.Int()%500]
}

func searchHotelParaGen(random bool, task *taskSetting) *params {
	if task.cbTraceRatio > 100 {
		task.cbTraceRatio = 100
	} else if task.cbTraceRatio < 0 {
		task.cbTraceRatio = 0
	}

	params := &params{
		size:   task.total,
		params: make([]*param, task.total),
	}

	for i := 0; i < task.total; i++ {
		in_date := 9
		out_date := 10
		lat := 37.7749
		lon := -122.4194

		if random {
			in_date = rand.Int()%14 + 9                        // 9 to 22
			out_date = rand.Int()%(24-in_date-1) + in_date + 1 // in_date to 23

			lat = 38.0235 + float64(rand.Int()%200-100)/1000.0
			lon = -122.095 + float64(rand.Int()%200-100)/1000.0
		}

		cbcID := task.cbcID
		if task.cbcID == configure.CBCID_DEFAULT {
			if rand.Int()%100 < task.cbTraceRatio {
				task.cbTraceCount++
			} else {
				cbcID = configure.CBCID_TRACEBYPASS
			}
		} else if task.cbcID == configure.CBCID_LOGGINGYPASS {
			if rand.Int()%100 < task.cbTraceRatio {
				task.cbTraceCount++
			} else {
				cbcID = configure.CBCID_OBSERVATIONBYPASS
			}
		}

		params.params[i] = &param{
			method: http.MethodGet,
			path: fmt.Sprintf("%s/hotels?cbcID=%d&inDate=2015-04-%02d&outDate=2015-04-%02d&lat=%.4f&lon=%.4f",
				URL, cbcID, in_date, out_date, lat, lon),
		}
	}

	return params
}

func recommendParaGen(random bool, task *taskSetting) *params {
	if !random {

	}

	params := &params{
		size:   task.total,
		params: make([]*param, task.total),
	}

	for i := 0; i < task.total; i++ {
		coin := rand.Float64()
		req_param := ""
		if coin < 0.33 {
			req_param = "dis"
		} else if coin < 0.66 {
			req_param = "rate"
		} else {
			req_param = "price"
		}

		lat := 38.0235 + float64(rand.Int()%200-100)/1000.0
		lon := -122.095 + float64(rand.Int()%200-100)/1000.0

		params.params[i] = &param{
			method: http.MethodGet,
			path: fmt.Sprintf("%s/recommendations?require=%s&lat=%.4f&lon=%.4f",
				URL, req_param, lat, lon),
		}
	}

	return params
}

func reserveParaGen(random bool, task *taskSetting) *params {
	params := &params{
		size:   task.total,
		params: make([]*param, task.total),
	}

	for i := 0; i < task.total; i++ {
		in_date := rand.Int()%14 + 9           // 9 to 22
		out_date := rand.Int()%4 + in_date + 1 // in_date + 1 to in_date + 4

		//lat := 38.0235 + float64(rand.Int()%200-100)/1000.0
		//lon := -122.095 + float64(rand.Int()%200-100)/1000.0

		hotel_id := rand.Int()%79 + 1 // 1 to 79
		user := getUser()
		cust_name := user.id
		num_room := 1

		params.params[i] = &param{
			method: http.MethodPost,
			path: fmt.Sprintf("%s/reservation?inDate=2015-04-%02d&outDate=2015-04-%02d&hotelId=%d&customerName=%s&username=%s&password=%s&number=%d",
				URL, in_date, out_date, hotel_id, cust_name, user.id, user.password, num_room),
		}
	}

	return params
}

func resetDB() error {
	resp, err := http.Post(fmt.Sprintf("%s/reset", URL), "application/x-www-form-urlencoded", nil)
	if err != nil {
		return err
	} else if b, err := io.ReadAll(resp.Body); err != nil {
		return err
	} else {
		//fmt.Println("Response:", string(b))
		_ = b
	}

	return nil
}

func getMetric() (*cb.PerfMetric, error) {
	resp, err := http.Get(fmt.Sprintf("%s/metric", URL))
	if err != nil {
		return nil, err
	} else if b, err := io.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		m := map[string]interface{}{}

		if err = json.Unmarshal(b, &m); err != nil {
			return nil, err
		}

		ok := false
		if m, ok = m["metric"].(map[string]interface{}); !ok {
			//return nil, errors.New("no metric found")
			return nil, nil
		}

		b, _ = json.Marshal(m)

		res := &cb.PerfMetric{}
		json.Unmarshal(b, res)

		return res, nil
	}
}

func client(reqPool chan *request, resPool chan *result, signal chan struct{}) {
	for {
		select {
		case <-signal:
			return
		case req := <-reqPool:
			res := &result{}

			var resp *http.Response
			var err error

			start := time.Now().UnixNano()
			if req.para.method == http.MethodGet {
				resp, err = http.Get(req.para.path)
			} else if req.para.method == http.MethodPost {
				resp, err = http.Post(req.para.path, "application/x-www-form-urlencoded", nil)
			}
			end := time.Now().UnixNano()
			res.latency = end - start

			if err != nil {
				res.err = true
				if rand.Int()%100 == 1 {
					log.Print(err)
				}
			} else {
				resp.Body.Close()
			}

			resPool <- res
		}
	}
}

func resultPool(results []*result, resPool chan *result, total int, signal chan struct{}) {
	for i := 0; i < total; i++ {
		results[i] = <-resPool
	}

	signal <- struct{}{}
}

func worker(paras *params, random bool, task *taskSetting) {
	results := make([]*result, task.total)
	signal := make(chan struct{}, task.threads)
	reqPool := make(chan *request, task.threads)
	resPool := make(chan *result, task.total)

	for i := 0; i < task.threads; i++ {
		go client(reqPool, resPool, signal)
	}

	resSig := make(chan struct{}, 1)
	go resultPool(results, resPool, task.total, resSig)

	start := time.Now().UnixNano()
	for i := 0; i < task.total; {
		startI := time.Now().UnixNano()
		expEndI := startI + time.Second.Nanoseconds()

		startID := i
		for ; i < startID+task.speed && i < task.total; i++ {
			para := paras.get(i)
			reqPool <- &request{
				id:   i,
				para: para,
			}
		}

		endI := time.Now().UnixNano()
		if endI < expEndI && i != task.total {
			fmt.Println("sleep for", time.Duration(expEndI-endI))
			time.Sleep(time.Duration(expEndI - endI))
		}
	}

	<-resSig
	end := time.Now().UnixNano()

	latency := float64(0)
	errCount := 0
	for i := 0; i < task.total; i++ {
		latency += float64(results[i].latency) / float64(task.total)
		if results[i].err {
			errCount++
		}
	}

	fmt.Printf("%v,%d,%d,%d,%v,%d,%s,%d\n",
		random, task.total, task.threads, task.speed, time.Duration(latency), int(float64(task.total)/(float64(end-start)/float64(time.Second.Nanoseconds()))), fmt.Sprintf("%d%%(%d)", task.cbTraceRatio, task.cbTraceCount), errCount)
}

type api func(random bool, task *taskSetting)

func run(a api, tasks []*taskSetting) {
	sleep := 5 * time.Second

	for i, task := range tasks {
		a(false, task)

		if i != len(tasks)-1 {
			time.Sleep(sleep)
		}
	}
}

var titleOnce = new(sync.Once)

func searchHotel(random bool, task *taskSetting) {
	paras := searchHotelParaGen(random, task)
	//fmt.Println("example path:", paras.params[0].path)
	if err := resetDB(); err != nil {
		fmt.Println("reset db fail:", err)
	} else if _, err = getMetric(); err != nil {
		fmt.Println("reset metric fail:", err)
	}

	worker(paras, random, task)

	time.Sleep(10 * time.Second)

	return
	perfMetric, err := getMetric()
	if err != nil {
		fmt.Println("get metric fail:", err)
		return
	} else if perfMetric == nil {
		return
	}

	titleOnce.Do(func() {
		fmt.Println("Type,Frontend,,,,,,,,,,,,,,,Search,,,,,,,,")
		fmt.Println(",L1,O1,L2,O2,O3,L4,O4,L5,O5,L6,O6,L7,LogicTotal,ObservationTotal,Total,O1,L2,O2,L3,O3,L4,LogicTotal,ObservationTotal,Total")
	})

	latency := perfMetric.Latency

	lFrontend := latency[cb.Metric_Frontend_SearchHandler_Logic_1].Median +
		latency[cb.Metric_Frontend_SearchHandler_Logic_2].Median +
		latency[cb.Metric_Frontend_SearchHandler_Logic_4].Median +
		latency[cb.Metric_Frontend_SearchHandler_Logic_5].Median +
		latency[cb.Metric_Frontend_SearchHandler_Logic_6].Median +
		latency[cb.Metric_Frontend_SearchHandler_Logic_7].Median
	oFrontend := latency[cb.Metric_Frontend_SearchHandler_1].Median +
		latency[cb.Metric_Frontend_SearchHandler_2].Median +
		latency[cb.Metric_Frontend_SearchHandler_3].Median +
		latency[cb.Metric_Frontend_SearchHandler_4].Median +
		latency[cb.Metric_Frontend_SearchHandler_5].Median +
		latency[cb.Metric_Frontend_SearchHandler_6].Median
	tFrontend := lFrontend + oFrontend

	lSearch := latency[cb.Metric_Search_NearBy_Logic_2].Median +
		latency[cb.Metric_Search_NearBy_Logic_3].Median +
		latency[cb.Metric_Search_NearBy_Logic_4].Median
	oSearch := latency[cb.Metric_Search_NearBy_Observation_1].Median +
		latency[cb.Metric_Search_NearBy_Observation_2].Median +
		latency[cb.Metric_Search_NearBy_Observation_3].Median
	tSearch := lSearch + oSearch

	fmt.Printf("Median,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f,%.0f\n",
		latency[cb.Metric_Frontend_SearchHandler_Logic_1].Median,
		latency[cb.Metric_Frontend_SearchHandler_1].Median,
		latency[cb.Metric_Frontend_SearchHandler_Logic_2].Median,
		latency[cb.Metric_Frontend_SearchHandler_2].Median,
		latency[cb.Metric_Frontend_SearchHandler_3].Median,
		latency[cb.Metric_Frontend_SearchHandler_Logic_4].Median,
		latency[cb.Metric_Frontend_SearchHandler_4].Median,
		latency[cb.Metric_Frontend_SearchHandler_Logic_5].Median,
		latency[cb.Metric_Frontend_SearchHandler_5].Median,
		latency[cb.Metric_Frontend_SearchHandler_Logic_6].Median,
		latency[cb.Metric_Frontend_SearchHandler_6].Median,
		latency[cb.Metric_Frontend_SearchHandler_Logic_7].Median,
		lFrontend, oFrontend, tFrontend,
		latency[cb.Metric_Search_NearBy_Observation_1].Median,
		latency[cb.Metric_Search_NearBy_Logic_2].Median,
		latency[cb.Metric_Search_NearBy_Observation_2].Median,
		latency[cb.Metric_Search_NearBy_Logic_3].Median,
		latency[cb.Metric_Search_NearBy_Observation_3].Median,
		latency[cb.Metric_Search_NearBy_Logic_4].Median,
		lSearch, oSearch, tSearch,
	)

	b, _ := json.Marshal(perfMetric.CBLatency)
	fmt.Println(string(b))
}

func recommend(random bool, task *taskSetting) {
	paras := recommendParaGen(random, task)
	//fmt.Println("example path:", paras.params[0].path)
	if err := resetDB(); err != nil {
		fmt.Println("reset db fail:", err)
	}

	worker(paras, random, task)
}

func reserve(random bool, task *taskSetting) {
	paras := reserveParaGen(random, task)
	//fmt.Println("example path:", paras.params[0].path)
	if err := resetDB(); err != nil {
		fmt.Println("reset db fail:", err)
	}

	worker(paras, random, task)
}
