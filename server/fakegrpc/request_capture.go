// Copyright © 2023 Yoshiki Shibata. All rights reserved.

package fakegrpc

import (
	"context"
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type RequestCapture[RQ, RS protoreflect.ProtoMessage] struct {
	lock sync.Mutex
	req  RQ

	res RS
	err error
}

func NewRequestCapture[RQ, RS protoreflect.ProtoMessage](
	res RS,
	err error,
) *RequestCapture[RQ, RS] {
	return &RequestCapture[RQ, RS]{
		res: res,
		err: err,
	}
}

func (rc *RequestCapture[RQ, RS]) Creator(
	ctx context.Context,
	req RQ,
) (RS, error) {
	rc.lock.Lock()
	rc.req = req
	rc.lock.Unlock()
	return rc.res, rc.err
}

func (rc *RequestCapture[RQ, RS]) Get() RQ {
	rc.lock.Lock()
	defer rc.lock.Unlock()

	return rc.req
}
