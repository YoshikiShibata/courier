// Copyright © 2023 Yoshiki Shibata. All rights reserved.

package fakegrpc

import (
	"context"
	"sync/atomic"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// ResponseBuilder builds a ResponseCreator which returns a sequence of response/errors.
type ResponseBuilder[RQ, RS protoreflect.ProtoMessage] struct {
	nextIndex uint32

	responses []RS
	errors    []error
}

func (rb *ResponseBuilder[RQ, RS]) Creator(
	ctx context.Context,
	req RQ,
) (RS, error) {
	i := int(atomic.AddUint32(&rb.nextIndex, 1) - 1)

	return rb.responses[i], rb.errors[i]
}

func (rb *ResponseBuilder[RQ, RS]) Append(
	res RS,
	err error,
) {
	rb.responses = append(rb.responses, res)
	rb.errors = append(rb.errors, err)
}
