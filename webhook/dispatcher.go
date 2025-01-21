package webhook

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

// Dispatcher is a dispatcher for webhook events.
type Dispatcher struct {
	secret            string
	followListeners   []FollowListener
	unFollowListeners []UnFollowListener
}

func (d *Dispatcher) Registers(listeners ...any) {
	for _, listener := range listeners {
		if l, ok := listener.(FollowListener); ok {
			d.FollowListener(l)
		}
		if l, ok := listener.(UnFollowListener); ok {
			d.UnFollowListener(l)
		}
	}
}

type Option func(*Dispatcher)

func WithRegisters(listeners ...any) Option {
	return func(d *Dispatcher) {
		d.Registers(listeners...)
	}
}

func WithSecret(secret string) Option {
	return func(d *Dispatcher) {
		d.secret = secret
	}
}

type dispatchRequestOptions struct {
	ctx context.Context
}

type DispatchRequestOption func(*dispatchRequestOptions)

func (d *Dispatcher) DispatchRequest(req *http.Request, opts ...DispatchRequestOption) error {
	o := &dispatchRequestOptions{
		ctx: req.Context(),
	}
	for _, opt := range opts {
		opt(o)
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Println("Error closing request body:", err)
		}
	}(req.Body)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return errors.New("line: webhook request body is null")
	}

	signature := req.Header.Get("x-line-signature")
	return d.DispatchBody(o.ctx, signature, body)
}

func (d *Dispatcher) DispatchBody(ctx context.Context, signature string, body []byte) error {
	events, err := ParseWebhookEvent(d.secret, signature, body)
	if err != nil {
		return err
	}
	return d.Dispatch(ctx, events)
}

// NewDispatcher returns a new Dispatcher instance.
func NewDispatcher(opts ...Option) *Dispatcher {
	dispatcher := &Dispatcher{}
	for _, opt := range opts {
		opt(dispatcher)
	}
	return dispatcher
}

func (d *Dispatcher) Dispatch(ctx context.Context, events []any) error {
	for _, event := range events {
		switch e := event.(type) {
		case *FollowEvent:
			return d.registerFollow(ctx, e)
		case *UnFollowEvent:
			return d.registerUnFollow(ctx, e)
		default:
			return errors.New("line: webhook dispatcher unsupported event")
		}
	}
	return nil
}

func (d *Dispatcher) registerFollow(ctx context.Context, event *FollowEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.followListeners {
		eg.Go(func() error {
			return listener.OnFollow(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) registerUnFollow(ctx context.Context, event *UnFollowEvent) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, listener := range d.unFollowListeners {
		eg.Go(func() error {
			return listener.OnUnFollow(ctx, event)
		})
	}
	return eg.Wait()
}

func (d *Dispatcher) FollowListener(listeners ...FollowListener) {
	d.followListeners = append(d.followListeners, listeners...)
}

func (d *Dispatcher) UnFollowListener(listeners ...UnFollowListener) {
	d.unFollowListeners = append(d.unFollowListeners, listeners...)
}
