// Code generated by Thrift Compiler (0.14.1). DO NOT EDIT.

package baggage

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/AleckDarcy/ContextBus/third-party/github.com/uber/jaeger-client-go/thrift"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = context.Background
var _ = time.Now
var _ = bytes.Equal

// Attributes:
//  - BaggageKey
//  - MaxValueLength
type BaggageRestriction struct {
	BaggageKey     string `thrift:"baggageKey,1,required" db:"baggageKey" json:"baggageKey"`
	MaxValueLength int32  `thrift:"maxValueLength,2,required" db:"maxValueLength" json:"maxValueLength"`
}

func NewBaggageRestriction() *BaggageRestriction {
	return &BaggageRestriction{}
}

func (p *BaggageRestriction) GetBaggageKey() string {
	return p.BaggageKey
}

func (p *BaggageRestriction) GetMaxValueLength() int32 {
	return p.MaxValueLength
}
func (p *BaggageRestriction) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	var issetBaggageKey bool = false
	var issetMaxValueLength bool = false

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
				issetBaggageKey = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		case 2:
			if fieldTypeId == thrift.I32 {
				if err := p.ReadField2(ctx, iprot); err != nil {
					return err
				}
				issetMaxValueLength = true
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	if !issetBaggageKey {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field BaggageKey is not set"))
	}
	if !issetMaxValueLength {
		return thrift.NewTProtocolExceptionWithType(thrift.INVALID_DATA, fmt.Errorf("Required field MaxValueLength is not set"))
	}
	return nil
}

func (p *BaggageRestriction) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(ctx); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.BaggageKey = v
	}
	return nil
}

func (p *BaggageRestriction) ReadField2(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(ctx); err != nil {
		return thrift.PrependError("error reading field 2: ", err)
	} else {
		p.MaxValueLength = v
	}
	return nil
}

func (p *BaggageRestriction) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "BaggageRestriction"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil {
			return err
		}
		if err := p.writeField2(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *BaggageRestriction) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "baggageKey", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:baggageKey: ", p), err)
	}
	if err := oprot.WriteString(ctx, string(p.BaggageKey)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.baggageKey (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:baggageKey: ", p), err)
	}
	return err
}

func (p *BaggageRestriction) writeField2(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "maxValueLength", thrift.I32, 2); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 2:maxValueLength: ", p), err)
	}
	if err := oprot.WriteI32(ctx, int32(p.MaxValueLength)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.maxValueLength (2) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 2:maxValueLength: ", p), err)
	}
	return err
}

func (p *BaggageRestriction) Equals(other *BaggageRestriction) bool {
	if p == other {
		return true
	} else if p == nil || other == nil {
		return false
	}
	if p.BaggageKey != other.BaggageKey {
		return false
	}
	if p.MaxValueLength != other.MaxValueLength {
		return false
	}
	return true
}

func (p *BaggageRestriction) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BaggageRestriction(%+v)", *p)
}

type BaggageRestrictionManager interface {
	// getBaggageRestrictions retrieves the baggage restrictions for a specific service.
	// Usually, baggageRestrictions apply to all services however there may be situations
	// where a baggageKey might only be allowed to be set by a specific service.
	//
	// Parameters:
	//  - ServiceName
	GetBaggageRestrictions(ctx context.Context, serviceName string) (_r []*BaggageRestriction, _err error)
}

type BaggageRestrictionManagerClient struct {
	c    thrift.TClient
	meta thrift.ResponseMeta
}

func NewBaggageRestrictionManagerClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *BaggageRestrictionManagerClient {
	return &BaggageRestrictionManagerClient{
		c: thrift.NewTStandardClient(f.GetProtocol(t), f.GetProtocol(t)),
	}
}

func NewBaggageRestrictionManagerClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *BaggageRestrictionManagerClient {
	return &BaggageRestrictionManagerClient{
		c: thrift.NewTStandardClient(iprot, oprot),
	}
}

func NewBaggageRestrictionManagerClient(c thrift.TClient) *BaggageRestrictionManagerClient {
	return &BaggageRestrictionManagerClient{
		c: c,
	}
}

func (p *BaggageRestrictionManagerClient) Client_() thrift.TClient {
	return p.c
}

func (p *BaggageRestrictionManagerClient) LastResponseMeta_() thrift.ResponseMeta {
	return p.meta
}

func (p *BaggageRestrictionManagerClient) SetLastResponseMeta_(meta thrift.ResponseMeta) {
	p.meta = meta
}

// getBaggageRestrictions retrieves the baggage restrictions for a specific service.
// Usually, baggageRestrictions apply to all services however there may be situations
// where a baggageKey might only be allowed to be set by a specific service.
//
// Parameters:
//  - ServiceName
func (p *BaggageRestrictionManagerClient) GetBaggageRestrictions(ctx context.Context, serviceName string) (_r []*BaggageRestriction, _err error) {
	var _args0 BaggageRestrictionManagerGetBaggageRestrictionsArgs
	_args0.ServiceName = serviceName
	var _result2 BaggageRestrictionManagerGetBaggageRestrictionsResult
	var _meta1 thrift.ResponseMeta
	_meta1, _err = p.Client_().Call(ctx, "getBaggageRestrictions", &_args0, &_result2)
	p.SetLastResponseMeta_(_meta1)
	if _err != nil {
		return
	}
	return _result2.GetSuccess(), nil
}

type BaggageRestrictionManagerProcessor struct {
	processorMap map[string]thrift.TProcessorFunction
	handler      BaggageRestrictionManager
}

func (p *BaggageRestrictionManagerProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

func (p *BaggageRestrictionManagerProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return processor, ok
}

func (p *BaggageRestrictionManagerProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

func NewBaggageRestrictionManagerProcessor(handler BaggageRestrictionManager) *BaggageRestrictionManagerProcessor {

	self3 := &BaggageRestrictionManagerProcessor{handler: handler, processorMap: make(map[string]thrift.TProcessorFunction)}
	self3.processorMap["getBaggageRestrictions"] = &baggageRestrictionManagerProcessorGetBaggageRestrictions{handler: handler}
	return self3
}

func (p *BaggageRestrictionManagerProcessor) Process(ctx context.Context, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	name, _, seqId, err2 := iprot.ReadMessageBegin(ctx)
	if err2 != nil {
		return false, thrift.WrapTException(err2)
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		return processor.Process(ctx, seqId, iprot, oprot)
	}
	iprot.Skip(ctx, thrift.STRUCT)
	iprot.ReadMessageEnd(ctx)
	x4 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	oprot.WriteMessageBegin(ctx, name, thrift.EXCEPTION, seqId)
	x4.Write(ctx, oprot)
	oprot.WriteMessageEnd(ctx)
	oprot.Flush(ctx)
	return false, x4

}

type baggageRestrictionManagerProcessorGetBaggageRestrictions struct {
	handler BaggageRestrictionManager
}

func (p *baggageRestrictionManagerProcessorGetBaggageRestrictions) Process(ctx context.Context, seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := BaggageRestrictionManagerGetBaggageRestrictionsArgs{}
	var err2 error
	if err2 = args.Read(ctx, iprot); err2 != nil {
		iprot.ReadMessageEnd(ctx)
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err2.Error())
		oprot.WriteMessageBegin(ctx, "getBaggageRestrictions", thrift.EXCEPTION, seqId)
		x.Write(ctx, oprot)
		oprot.WriteMessageEnd(ctx)
		oprot.Flush(ctx)
		return false, thrift.WrapTException(err2)
	}
	iprot.ReadMessageEnd(ctx)

	tickerCancel := func() {}
	// Start a goroutine to do server side connectivity check.
	if thrift.ServerConnectivityCheckInterval > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()
		var tickerCtx context.Context
		tickerCtx, tickerCancel = context.WithCancel(context.Background())
		defer tickerCancel()
		go func(ctx context.Context, cancel context.CancelFunc) {
			ticker := time.NewTicker(thrift.ServerConnectivityCheckInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if !iprot.Transport().IsOpen() {
						cancel()
						return
					}
				}
			}
		}(tickerCtx, cancel)
	}

	result := BaggageRestrictionManagerGetBaggageRestrictionsResult{}
	var retval []*BaggageRestriction
	if retval, err2 = p.handler.GetBaggageRestrictions(ctx, args.ServiceName); err2 != nil {
		tickerCancel()
		if err2 == thrift.ErrAbandonRequest {
			return false, thrift.WrapTException(err2)
		}
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing getBaggageRestrictions: "+err2.Error())
		oprot.WriteMessageBegin(ctx, "getBaggageRestrictions", thrift.EXCEPTION, seqId)
		x.Write(ctx, oprot)
		oprot.WriteMessageEnd(ctx)
		oprot.Flush(ctx)
		return true, thrift.WrapTException(err2)
	} else {
		result.Success = retval
	}
	tickerCancel()
	if err2 = oprot.WriteMessageBegin(ctx, "getBaggageRestrictions", thrift.REPLY, seqId); err2 != nil {
		err = thrift.WrapTException(err2)
	}
	if err2 = result.Write(ctx, oprot); err == nil && err2 != nil {
		err = thrift.WrapTException(err2)
	}
	if err2 = oprot.WriteMessageEnd(ctx); err == nil && err2 != nil {
		err = thrift.WrapTException(err2)
	}
	if err2 = oprot.Flush(ctx); err == nil && err2 != nil {
		err = thrift.WrapTException(err2)
	}
	if err != nil {
		return
	}
	return true, err
}

// HELPER FUNCTIONS AND STRUCTURES

// Attributes:
//  - ServiceName
type BaggageRestrictionManagerGetBaggageRestrictionsArgs struct {
	ServiceName string `thrift:"serviceName,1" db:"serviceName" json:"serviceName"`
}

func NewBaggageRestrictionManagerGetBaggageRestrictionsArgs() *BaggageRestrictionManagerGetBaggageRestrictionsArgs {
	return &BaggageRestrictionManagerGetBaggageRestrictionsArgs{}
}

func (p *BaggageRestrictionManagerGetBaggageRestrictionsArgs) GetServiceName() string {
	return p.ServiceName
}
func (p *BaggageRestrictionManagerGetBaggageRestrictionsArgs) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if fieldTypeId == thrift.STRING {
				if err := p.ReadField1(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *BaggageRestrictionManagerGetBaggageRestrictionsArgs) ReadField1(ctx context.Context, iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(ctx); err != nil {
		return thrift.PrependError("error reading field 1: ", err)
	} else {
		p.ServiceName = v
	}
	return nil
}

func (p *BaggageRestrictionManagerGetBaggageRestrictionsArgs) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "getBaggageRestrictions_args"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField1(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *BaggageRestrictionManagerGetBaggageRestrictionsArgs) writeField1(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin(ctx, "serviceName", thrift.STRING, 1); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field begin error 1:serviceName: ", p), err)
	}
	if err := oprot.WriteString(ctx, string(p.ServiceName)); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T.serviceName (1) field write error: ", p), err)
	}
	if err := oprot.WriteFieldEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write field end error 1:serviceName: ", p), err)
	}
	return err
}

func (p *BaggageRestrictionManagerGetBaggageRestrictionsArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BaggageRestrictionManagerGetBaggageRestrictionsArgs(%+v)", *p)
}

// Attributes:
//  - Success
type BaggageRestrictionManagerGetBaggageRestrictionsResult struct {
	Success []*BaggageRestriction `thrift:"success,0" db:"success" json:"success,omitempty"`
}

func NewBaggageRestrictionManagerGetBaggageRestrictionsResult() *BaggageRestrictionManagerGetBaggageRestrictionsResult {
	return &BaggageRestrictionManagerGetBaggageRestrictionsResult{}
}

var BaggageRestrictionManagerGetBaggageRestrictionsResult_Success_DEFAULT []*BaggageRestriction

func (p *BaggageRestrictionManagerGetBaggageRestrictionsResult) GetSuccess() []*BaggageRestriction {
	return p.Success
}
func (p *BaggageRestrictionManagerGetBaggageRestrictionsResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *BaggageRestrictionManagerGetBaggageRestrictionsResult) Read(ctx context.Context, iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read error: ", p), err)
	}

	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin(ctx)
		if err != nil {
			return thrift.PrependError(fmt.Sprintf("%T field %d read error: ", p, fieldId), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if fieldTypeId == thrift.LIST {
				if err := p.ReadField0(ctx, iprot); err != nil {
					return err
				}
			} else {
				if err := iprot.Skip(ctx, fieldTypeId); err != nil {
					return err
				}
			}
		default:
			if err := iprot.Skip(ctx, fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(ctx); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T read struct end error: ", p), err)
	}
	return nil
}

func (p *BaggageRestrictionManagerGetBaggageRestrictionsResult) ReadField0(ctx context.Context, iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin(ctx)
	if err != nil {
		return thrift.PrependError("error reading list begin: ", err)
	}
	tSlice := make([]*BaggageRestriction, 0, size)
	p.Success = tSlice
	for i := 0; i < size; i++ {
		_elem5 := &BaggageRestriction{}
		if err := _elem5.Read(ctx, iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", _elem5), err)
		}
		p.Success = append(p.Success, _elem5)
	}
	if err := iprot.ReadListEnd(ctx); err != nil {
		return thrift.PrependError("error reading list end: ", err)
	}
	return nil
}

func (p *BaggageRestrictionManagerGetBaggageRestrictionsResult) Write(ctx context.Context, oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin(ctx, "getBaggageRestrictions_result"); err != nil {
		return thrift.PrependError(fmt.Sprintf("%T write struct begin error: ", p), err)
	}
	if p != nil {
		if err := p.writeField0(ctx, oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(ctx); err != nil {
		return thrift.PrependError("write field stop error: ", err)
	}
	if err := oprot.WriteStructEnd(ctx); err != nil {
		return thrift.PrependError("write struct stop error: ", err)
	}
	return nil
}

func (p *BaggageRestrictionManagerGetBaggageRestrictionsResult) writeField0(ctx context.Context, oprot thrift.TProtocol) (err error) {
	if p.IsSetSuccess() {
		if err := oprot.WriteFieldBegin(ctx, "success", thrift.LIST, 0); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field begin error 0:success: ", p), err)
		}
		if err := oprot.WriteListBegin(ctx, thrift.STRUCT, len(p.Success)); err != nil {
			return thrift.PrependError("error writing list begin: ", err)
		}
		for _, v := range p.Success {
			if err := v.Write(ctx, oprot); err != nil {
				return thrift.PrependError(fmt.Sprintf("%T error writing struct: ", v), err)
			}
		}
		if err := oprot.WriteListEnd(ctx); err != nil {
			return thrift.PrependError("error writing list end: ", err)
		}
		if err := oprot.WriteFieldEnd(ctx); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T write field end error 0:success: ", p), err)
		}
	}
	return err
}

func (p *BaggageRestrictionManagerGetBaggageRestrictionsResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BaggageRestrictionManagerGetBaggageRestrictionsResult(%+v)", *p)
}
