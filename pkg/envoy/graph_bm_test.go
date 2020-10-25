package envoy

import (
	"context"
	"fmt"
	"testing"
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

// Run with:
//  go test -v -benchmem -cpuprofile cpu.out -bench DepRes  ./pkg/envoy/*.go
//
//  go tool pprof cpu.out
//  (pprof) list Next
//  (pprof) ...
func BenchmarkGraph_Next(b *testing.B) {
	b.ReportAllocs()

	var (
		minNodes   = getEnvInt("MIN_NODES", 100)
		maxNodes   = getEnvInt("MAX_NODES", 100000000)
		stepFactor = getEnvInt("STEP_FACTOR", 10)
	)

	for nodes := minNodes; nodes <= maxNodes; nodes *= stepFactor {
		b.Run(fmt.Sprintf("%d", nodes), func(b *testing.B) {
			var g = NewGraph()
			g.Add(&pNode{ii: NodeIdentifiers{"p0"}})
			g.Add(&pNode{ii: NodeIdentifiers{"p1"}})
			g.Add(&pNode{ii: NodeIdentifiers{"p2"}})
			g.Add(&pNode{ii: NodeIdentifiers{"p3"}})

			for n := 0; n < nodes; n++ {
				g.Add(&cNode{
					ii: NodeIdentifiers{fmt.Sprintf("c%d", n)},
					rr: NodeRelationships{pNode{}.Resource(): NodeIdentifiers{fmt.Sprintf("p%d", n%4)}},
				})
			}

			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				b.ReportMetric(float64(len(g.nodes.set)), "nodes")

				if Encode(context.Background(), g) != nil {
					b.FailNow()
				}
			}
		})
	}

}
