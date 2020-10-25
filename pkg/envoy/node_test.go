package envoy

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func makeTestNode(ident string, rr NodeRelationships) *TestNode {
	return &TestNode{
		rr: rr,
		ii: NodeIdentifiers{ident},
	}
}

func TestNodeSet_Has(t *testing.T) {
	var (
		req = require.New(t)

		set = NodeSet{}
		n   = makeTestNode("tn1", nil)
	)

	req.False(set.Has(n))
	set.Add(n)
	req.True(set.Has(n))
	set.Remove(n)
	req.False(set.Has(n))
}

func TestNodeSet_Remove(t *testing.T) {
	var (
		req = require.New(t)

		set = NodeSet{}
		n1  = makeTestNode("tn1", nil)
		n2  = makeTestNode("tn2", nil)
		n3  = makeTestNode("tn3", nil)
	)

	req.Len(set.set, 0)
	set.Add(n1, n2, n3)
	req.Len(set.set, 3)
	set.Remove(n1, n2)
	req.Len(set.set, 3)
	set.Add(n3)
	req.Len(set.set, 3)
}

func BenchmarkNodeSet_indexOf(b *testing.B) {
	b.ReportAllocs()

	var (
		pNodes     = 4
		minNodes   = getEnvInt("MIN_NODES", 100)
		maxNodes   = getEnvInt("MAX_NODES", 100000000)
		stepFactor = getEnvInt("STEP_FACTOR", 10)
	)

	for nodes := minNodes; nodes <= maxNodes; nodes *= stepFactor {
		b.Run(fmt.Sprintf("%d", nodes), func(b *testing.B) {
			var (
				set    = NodeSet{}
				parent = make([]Node, pNodes)
			)

			for f := 0; f < pNodes; f++ {
				parent[f] = makeTestNode(fmt.Sprintf("parent_%d", f), nil)

				for n := 0; n < nodes; n++ {
					set.Add(makeTestNode(fmt.Sprintf("n%d", n), nil))
				}
			}

			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				b.ReportMetric(float64(len(set.set)), "nodes")
				set.indexOf(parent[n%pNodes])
			}
		})
	}
}
