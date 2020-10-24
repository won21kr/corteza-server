package envoy

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	pNode struct {
		rr NodeRelationships
		ii NodeIdentifiers
	}
	cNode struct {
		rr NodeRelationships
		ii NodeIdentifiers
	}
)

func (n pNode) Resource() string             { return "parent" }
func (n pNode) Relations() NodeRelationships { return n.rr }
func (n pNode) Identifiers() NodeIdentifiers { return n.ii }
func (n *pNode) Update(nn ...Node)           {}
func (n pNode) Matches(r string, ii ...string) bool {
	return r == n.Resource() && n.Identifiers().HasAny(ii...)
}

func (n cNode) Resource() string             { return "child" }
func (n cNode) Relations() NodeRelationships { return n.rr }
func (n cNode) Identifiers() NodeIdentifiers { return n.ii }
func (n *cNode) Update(rr ...Node) {
	for _, r := range rr {
		println("updating child node: ", r)
		if p, is := r.(pNode); is && r.Identifiers().HasAny(p.ii...) {
			n.ii = nil
		}
	}
}
func (n cNode) Matches(r string, ii ...string) bool {
	return r == n.Resource() && n.Identifiers().HasAny(ii...)
}

// Run with:
//  go test -v -benchmem -cpuprofile cpu.out -bench DepRes  ./pkg/envoy/*.go
//
//  go tool pprof cpu.out
//  (pprof) list Next
//  (pprof) ...
func BenchmarkGraph_DepResolution(b *testing.B) {
	var (
		g   = NewGraph()
		req = require.New(b)
	)

	println(b.N)

	g.Add(&pNode{ii: NodeIdentifiers{"p0"}})
	g.Add(&pNode{ii: NodeIdentifiers{"p1"}})
	g.Add(&pNode{ii: NodeIdentifiers{"p2"}})
	g.Add(&pNode{ii: NodeIdentifiers{"p3"}})

	for n := 0; n < b.N; n++ {
		g.Add(&cNode{
			ii: NodeIdentifiers{fmt.Sprintf("c%d", n)},
			rr: NodeRelationships{pNode{}.Resource(): NodeIdentifiers{fmt.Sprintf("p%d", n%4)}},
		})
	}

	b.ResetTimer()
	req.NoError(Encode(context.Background(), g))
}
