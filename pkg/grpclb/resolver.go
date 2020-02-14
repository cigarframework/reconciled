package grpclb

import (
	"fmt"

	"go.uber.org/atomic"
	"google.golang.org/grpc/resolver"
)

var counter = atomic.NewInt64(0)

const Scheme = "lb"

type ResolverBuilder struct {
	schema  string
	service string
	addrs   []string
}

func NewResolverBuilder(addrs []string) *ResolverBuilder {
	return &ResolverBuilder{
		schema:  fmt.Sprintf("grpc_lb_schema_%d", counter.Inc()),
		service: fmt.Sprintf("grpc_lb_service_%d", counter.Inc()),
		addrs:   addrs,
	}
}

func (b *ResolverBuilder) Addr() string {
	return fmt.Sprintf("%s:///%s", b.schema, b.service)
}

func (b *ResolverBuilder) Scheme() string {
	return b.schema
}

func (b *ResolverBuilder) Service() string {
	return b.service
}

func (b *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &customResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			b.service: b.addrs,
		},
	}
	r.start()
	return r, nil
}

type customResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *customResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*customResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*customResolver) Close()                                  {}
