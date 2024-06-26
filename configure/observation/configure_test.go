package observation

import (
	"runtime/debug"

	"github.com/AleckDarcy/ContextBus/helper"
	cb "github.com/AleckDarcy/ContextBus/proto"

	"testing"
	"time"
)

func BenchmarkName(b *testing.B) {
	attrCfgs := []*cb.AttributeConfigure{
		cb.Test_AttributeConfigure_App_Key21,
		cb.Test_AttributeConfigure_App_Message,
		cb.Test_AttributeConfigure_Lib1_Key11,
	}
	cfg := cb.NewLoggingConfigure(nil, nil, attrCfgs, cb.LogOutType_LogOutType_)

	what := new(cb.EventWhat)
	what.WithApplication(nil).
		SetMessage("application message").GetAttributes().SetString("key1", "value1").
		WithAttributes("key2", nil).
		SetString("key21", "value21")
	what.WithLibrary("lib1", nil).GetAttributes().WithAttributes("key1", nil).SetString("key11", "value11")

	er := new(cb.EventRepresentation)
	er.WithWhen(&cb.EventWhen{Time: time.Now().UnixNano()})
	er.WithWhat(what)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		//b.Log((*LoggingConfigure)(cfg).Do(&cb.EventData{Event: er}))
		(*LoggingConfigure)(cfg).Do(&cb.EventData{Event: er})
	}
}

func BenchmarkStackTrace(b *testing.B) {
	a := ""
	fb := func() {
		a = helper.BytesToString(debug.Stack())
		//a = string(debug.Stack())
	}
	fa := func() {
		fb()
	}

	for i := 0; i < b.N; i++ {
		fa()
	}

	_ = a
}
