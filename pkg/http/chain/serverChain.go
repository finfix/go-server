package chain

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"

	"pkg/panicRecover"
)

type Chain struct {
	before      []BeforeFunc
	send        SendFunc
	after       []AfterFunc
	encode      EncodeResponseFunc
	errorEncode EncodeErrorFunc
}

func NewChain(send SendFunc, opts ...Option) *Chain {
	chain := &Chain{
		before:      nil,
		send:        send,
		after:       nil,
		encode:      DefaultResponseEncoder,
		errorEncode: DefaultErrorEncoder,
	}
	for _, option := range opts {
		option(chain)
	}
	return chain
}

type Option func(*Chain)

func ResponseEncoder(e EncodeResponseFunc) Option {
	return func(s *Chain) { s.encode = e }
}

func ErrorEncoder(ee EncodeErrorFunc) Option {
	return func(s *Chain) { s.errorEncode = ee }
}

func Before(before ...BeforeFunc) Option {
	return func(s *Chain) { s.before = append(s.before, before...) }
}

func After(after ...AfterFunc) Option {
	return func(s *Chain) { s.after = append(s.after, after...) }
}

type BeforeFunc func(context.Context, *http.Request) (context.Context, error)

type SendFunc func(context.Context, *http.Request) (any, error)

type AfterFunc func(context.Context, http.ResponseWriter) context.Context

type EncodeResponseFunc func(context.Context, http.ResponseWriter, any) error

type EncodeErrorFunc func(context.Context, http.ResponseWriter, error)

type LoggingFunc func(error)

var tracer = otel.Tracer("/pkg/http/chain")

func (s *Chain) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx, span := tracer.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.Path))
	defer span.End()

	var (
		err error
		res any
	)

	defer panicRecover.PanicRecover(func(err error) {
		s.errorEncode(ctx, w, err)
	})

	for _, f := range s.before {
		ctx, err = f(ctx, r)
		if err != nil {
			s.errorEncode(ctx, w, err)
			return
		}
	}

	res, err = s.send(ctx, r)
	if err != nil {
		s.errorEncode(ctx, w, err)
		return
	}

	if err = s.encode(ctx, w, res); err != nil {
		s.errorEncode(ctx, w, err)
		return
	}
}
