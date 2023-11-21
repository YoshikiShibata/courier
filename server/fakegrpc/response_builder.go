// Copyright © 2023 Yoshiki Shibata. All rights reserved.

package fakegrpc

import (
	"context"
	"sync/atomic"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// ResponseBuilder builds a ResponseCreator which returns a sequence of response/errors.
type ResponseBuilder[Req, Res protoreflect.ProtoMessage] struct {
	nextIndex uint32

	responses []Res
	errors    []error
}

func (rb *ResponseBuilder[Req, Res]) Creator(
	ctx context.Context,
	req Req,
) (response Res, err error) {
	i := int(atomic.AddUint32(&rb.nextIndex, 1) - 1)

	return rb.responses[i], rb.errors[i]
}

func (rb *ResponseBuilder[Req, Res]) Append(
	res Res,
	err error,
) {
	rb.responses = append(rb.responses, res)
	rb.errors = append(rb.errors, err)
}
