package helper

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/rs/zerolog/pkgerrors"

	"github.com/pkg/errors"

	"github.com/rs/zerolog"
)

const SIZE_OF_UINTPTR = int(unsafe.Sizeof(uintptr(0)))
const MAX_CALLERS_LEN = 50

type callerFrame struct {
	frame *runtime.Frame
	str   string
}

var newCount = 0

type callerFrameStore struct {
	store map[uintptr]*callerFrame
	lock  sync.RWMutex
}

var CallerFrameStore = &callerFrameStore{store: map[uintptr]*callerFrame{}}

func (s *callerFrameStore) jsonMarshaller(frame *runtime.Frame) string {
	function, source := frame.Function, frame.File
	if id := strings.LastIndex(function, "("); id >= 0 {
		function = function[id:]
	} else {
		function = function[strings.LastIndex(function, ".")+1:]
	}
	source = source[strings.LastIndex(source, "/")+1:]

	return fmt.Sprintf("{\"func\":\"%s\",\"line\":\"%d\",\"source\":\"%s\"}", function, frame.Line, source)
}

func (s *callerFrameStore) goMarshaller(frame *runtime.Frame) string {
	return fmt.Sprintf("%s()\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
}

func (s *callerFrameStore) Get(pcs []uintptr) []*callerFrame {
	callers := make([]*callerFrame, len(pcs))
	ok, newPCSFlag := false, false
	s.lock.RLock()
	for i, pc := range pcs {
		callers[i], ok = s.store[pc]
		newPCSFlag = newPCSFlag || !ok
	}
	s.lock.RUnlock()

	if newPCSFlag {
		newPCS := make([]int, 0, len(pcs))
		for i, pc := range pcs {
			if callers[i] == nil {
				frames := runtime.CallersFrames([]uintptr{pc})
				frame, _ := frames.Next()
				callers[i] = &callerFrame{
					frame: &frame,
					str:   s.jsonMarshaller(&frame),
				}
				newPCS = append(newPCS, i)
			}
		}

		s.lock.Lock()
		for i := range newPCS {
			s.store[pcs[i]] = callers[i]
		}
		s.lock.Unlock()
	}

	return callers
}

type stacktraceStore struct {
	store map[string][]byte
	lock  sync.RWMutex
	pool  chan []uintptr
}

var StacktraceStore = &stacktraceStore{
	store: map[string][]byte{},
	pool:  make(chan []uintptr, 1000),
}

func (s *stacktraceStore) jsonMarshaller(callers []*callerFrame) []byte {
	bfLen := len(callers) + 1
	for _, caller := range callers {
		bfLen += len(caller.str)
	}
	bytes := make([]byte, bfLen)
	bytes[0] = '['
	head := 1
	for _, caller := range callers {
		head += copy(bytes[head:], caller.str)
		bytes[head] = ','
		head++
	}
	bytes[head-1] = ']'

	return bytes
}

func (s *stacktraceStore) goMarshaller(callers []*callerFrame) []byte {
	bfLen := 0
	for _, caller := range callers {
		bfLen += len(caller.str)
	}
	bytes := make([]byte, bfLen)
	head := 0
	for _, caller := range callers {
		head += copy(bytes[head:], caller.str)
	}

	return bytes
}

func (s *stacktraceStore) Get(pcs []uintptr) []byte {
	pcsBytes := pcsToBytes(pcs)

	s.lock.RLock()
	bytes, ok := s.store[BytesToString(pcsBytes)]
	s.lock.RUnlock()

	if !ok {
		callers := CallerFrameStore.Get(pcs)
		bytes = s.jsonMarshaller(callers)

		s.lock.Lock()
		s.store[string(pcsBytes)] = bytes
		s.lock.Unlock()
	}

	return bytes
}

func pcsToBytes(pcs []uintptr) []byte {
	pcsHeader := (*reflect.StringHeader)(unsafe.Pointer(&pcs))
	b := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: pcsHeader.Data,
		Len:  pcsHeader.Len * SIZE_OF_UINTPTR,
		Cap:  pcsHeader.Len * SIZE_OF_UINTPTR,
	}))

	runtime.KeepAlive(pcs)

	return b
}

func GetStacktrace() []byte {
	var pcs []uintptr
	select {
	case pcs = <-StacktraceStore.pool:
	default:
		newCount++
		pcs = make([]uintptr, MAX_CALLERS_LEN)
	}

	num := runtime.Callers(2, pcs)
	pcs = pcs[0:num]
	bytes := StacktraceStore.Get(pcs)

	select {
	case StacktraceStore.pool <- pcs[:MAX_CALLERS_LEN]:
		// recycled
	default:
		// dropped
	}

	return bytes
}

func f1() {
	GetStacktrace()
}

func f2() {
	f1()
}

func TestName(t *testing.T) {
	print(string(GetStacktrace()))
	////
	f1()
	////
	f2()
}

func BenchmarkDebugStack(b *testing.B) {
	var bytes []byte
	for i := 0; i < b.N; i++ {
		bytes = debug.Stack()
	}
	b.StopTimer()
	_ = bytes
}

var N = 500000

func TestCachedStack(t *testing.T) {
	var bytes []byte
	var str string

	start := time.Now().UnixNano()
	for i := 0; i < N; i++ {
		bytes = GetStacktrace()
		str = BytesToString(bytes)
		err := errors.Wrap(errors.New("error message"), "from error")
		str = fmt.Sprintf("{\"stack\":%s,\"error\":\"%s\"}", str, err)

		//fmt.Println(str)
	}
	end := time.Now().UnixNano()

	t.Log((end - start) / int64(N))
	_ = bytes
}

func BenchmarkPrint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Println(BytesToString(make([]byte, 292)))
	}
}

func TestZeroLog(t *testing.T) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	var str string
	out := &bytes.Buffer{}

	start := time.Now().UnixNano()
	for i := 0; i < N; i++ {
		log := zerolog.New(out).With().Stack().Logger()
		err := errors.Wrap(errors.New("error message"), "from error")
		log.Log().Err(err).Msg("") // not explicitly calling Stack()
		str = out.String()
		out.Reset()

		//fmt.Print(str)
	}
	end := time.Now().UnixNano()
	_ = str
	t.Log((end - start) / int64(N))
}

func TestA(t *testing.T) {

}
