// Copyright © Yoshiki Shibata, Inc. All rights reserved.

package server

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/YoshikiShibata/courier/tid"
)

type responseInfo struct {
	res     any
	err     error
	creator ResponseCreator
}

// GRPCServerStub provides basic primitives for a fake gRPC server.
type GRPCServerStub struct {
	mutex     sync.Mutex
	responses map[string]*responseInfo
}

// SetResponse sets specified response and error for a RPC(rcpName).
// tid represents Testing ID.
func (ss *GRPCServerStub) SetResponse(
	tid, rpcName string,
	response any,
	err error,
) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	if ss.responses == nil {
		ss.responses = make(map[string]*responseInfo)
	}

	// Always overwrite.
	ss.responses[tid+":"+rpcName] = &responseInfo{
		res:     response,
		err:     err,
		creator: nil,
	}
}

// ResponseCreator is a function to create response and error.
type ResponseCreator func(ctx context.Context, req any) (response any, err error)

// SetResponseCreator sets a creator for a RPC.
func (ss *GRPCServerStub) SetResponseCreator(
	tid, rpcName string,
	creator ResponseCreator,
) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	if ss.responses == nil {
		ss.responses = make(map[string]*responseInfo)
	}

	// Always overwrite.
	ss.responses[tid+":"+rpcName] = &responseInfo{
		res:     nil,
		err:     nil,
		creator: creator,
	}
}

// ClearAllResponses clear all responsed and errors for tid.
func (ss *GRPCServerStub) ClearAllResponses(
	tid string,
) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()

	prefix := tid + ":"

	for k := range ss.responses {
		if strings.HasPrefix(k, prefix) {
			delete(ss.responses, k)
		}
	}
}

// HandleRequest is a generic function to handle requests.
func HandleRequest[R any](
	ss *GRPCServerStub,
	ctx context.Context,
	rpcName string,
	req any,
) (R, error) {
	var zeroValue R

	res, err := ss.HandleRequest(ctx, rpcName, req)
	if err != nil {
		return zeroValue, err
	}
	return res.(R), nil
}

// HandleRequest is a method to handle requestes.
func (ss *GRPCServerStub) HandleRequest(
	ctx context.Context,
	rpcName string,
	req any,
) (any, error) {
	key := tid.Extract(ctx) + ":" + rpcName

	ss.mutex.Lock()
	resInfo := ss.responses[key]
	ss.mutex.Unlock()

	if resInfo == nil {
		panic(fmt.Sprintf("Response for %s has not been set yet", key))
	}
	if resInfo.creator == nil {
		return resInfo.res, resInfo.err
	}
	return resInfo.creator(ctx, req)
}
