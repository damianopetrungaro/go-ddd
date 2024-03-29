// Code generated by gowrap. DO NOT EDIT.
// template: https://raw.githubusercontent.com/hexdigest/gowrap/6c8f05695fec23df85903a8da0af66ac414e2a63/templates/opentelemetry
// gowrap: http://github.com/hexdigest/gowrap

package instrument

//go:generate gowrap gen -p github.com/organization/order-service -i Repo -t https://raw.githubusercontent.com/hexdigest/gowrap/6c8f05695fec23df85903a8da0af66ac414e2a63/templates/opentelemetry -o trace.go -l ""

import (
	"context"

	"github.com/organization/order-service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// RepoWithTracing implements order.Repo interface instrumented with opentracing spans
type RepoWithTracing struct {
	order.Repo
	_instance      string
	_spanDecorator func(span trace.Span, params, results map[string]interface{})
}

// NewRepoWithTracing returns RepoWithTracing
func NewRepoWithTracing(base order.Repo, instance string, spanDecorator ...func(span trace.Span, params, results map[string]interface{})) RepoWithTracing {
	d := RepoWithTracing{
		Repo:      base,
		_instance: instance,
	}

	if len(spanDecorator) > 0 && spanDecorator[0] != nil {
		d._spanDecorator = spanDecorator[0]
	}

	return d
}

// Add implements order.Repo
func (_d RepoWithTracing) Add(ctx context.Context, order *order.Order) (err error) {
	ctx, _span := otel.Tracer(_d._instance).Start(ctx, "order.Repo.Add")
	defer func() {
		if _d._spanDecorator != nil {
			_d._spanDecorator(_span, map[string]interface{}{
				"ctx":   ctx,
				"order": order}, map[string]interface{}{
				"err": err})
		} else if err != nil {
			_span.RecordError(err)
			_span.SetAttributes(
				attribute.String("event", "error"),
				attribute.String("message", err.Error()),
			)
		}

		_span.End()
	}()
	return _d.Repo.Add(ctx, order)
}

// Get implements order.Repo
func (_d RepoWithTracing) Get(ctx context.Context, id order.ID) (op1 *order.Order, err error) {
	ctx, _span := otel.Tracer(_d._instance).Start(ctx, "order.Repo.Get")
	defer func() {
		if _d._spanDecorator != nil {
			_d._spanDecorator(_span, map[string]interface{}{
				"ctx": ctx,
				"id":  id}, map[string]interface{}{
				"op1": op1,
				"err": err})
		} else if err != nil {
			_span.RecordError(err)
			_span.SetAttributes(
				attribute.String("event", "error"),
				attribute.String("message", err.Error()),
			)
		}

		_span.End()
	}()
	return _d.Repo.Get(ctx, id)
}
