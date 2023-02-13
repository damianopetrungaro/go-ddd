package instrument

import "github.com/organization/order-service"

// New returns an instrumented order.Repo
func New(base order.Repo, name string) order.Repo {
	return NewRepoWithPrometheus(NewRepoWithTracing(base, name), name)
}
