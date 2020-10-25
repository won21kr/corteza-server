package envoy

import (
	"github.com/cortezaproject/corteza-server/pkg/slice"
)

type (
	// Node defines the signature of any valid graph node
	Node interface {
		// Identifiers returns a set of values that identify the node.
		//
		// The identifiers may **not** be unique across all resources, but
		// they **must** be unique inside a given resource.
		Identifiers() NodeIdentifiers

		// Resource returns the Corteza resource identifier that this node handles
		Resource() string

		// Relations returns a set of NodeRelationships regarding this node
		//
		// The graph layer **must** be able to handle dynamic relationships (changed in runtime).
		Relations() NodeRelationships
	}

	// NodeUpdater defines a node that can update its state based on the given set of Nodes
	//
	// For example, a ComposeRecordNode should know how to update the referenced ComposeModule resource.
	NodeUpdater interface {
		Node

		// Update receives a set of nodes that should be used when updating the given node n
		//
		// The caller **must** only provide nodes that the given node n is dependent of (it's parent nodes).
		Update(...Node)
	}

	nodeMatcher interface {
		// Matches checks if the node matches the given resources and **any** of the identifiers.
		Matches(resource string, identifiers ...string) bool
	}

	state int8

	// NodeSet is a set of Nodes
	NodeSet struct {
		set []Node

		// node index per resource types and identifiers
		index map[string]map[string]int

		// rr - node relationship index
		rIndex map[string]map[string][]int

		// node state index
		sIndex map[int]state

		// count per state
		sCount map[state]int

		// record empty spots
		holes []int
	}

	// NodeRelationships holds relationships for a specific node
	NodeRelationships map[string]NodeIdentifiers

	// NodeIdentifiers represents a set of node identifiers
	NodeIdentifiers []string
)

const (
	nodeInitCap = 2 << 16
)

const (
	pending state = 1 << iota
	conflicting
	conflictingParents
	processed

	all         = processed & pending & conflicting & conflictingParents
	unprocessed = pending & conflicting & conflictingParents
)

// Add adds a new identifier for the given resource
func (n NodeRelationships) Add(resource string, ii ...string) {
	if _, has := n[resource]; !has {
		n[resource] = make(NodeIdentifiers, 0, 1)
	}

	n[resource] = n[resource].Add(ii...)
}

// Add adds a new identifiers
//
// Identifier can be string, Stringer or uint64 and should not be empty (zero)
func (ii NodeIdentifiers) Add(idents ...string) NodeIdentifiers {
	m := slice.ToStringBoolMap(ii)

	for _, i := range idents {
		if len(i) == 0 || m[i] {
			// skip all empty and existing identifiers
			continue
		}

		ii = append(ii, i)
	}

	return ii
}

// HasAny checks if any of the provided identifiers appear in the given set of identifiers
func (ii NodeIdentifiers) HasAny(jj ...string) bool {
	for _, i := range ii {
		for _, j := range jj {
			if i == j {
				return true
			}
		}
	}

	return false
}

func NewNodeSet() *NodeSet {
	var (
		secCap = 128
	)

	set := &NodeSet{
		set:    make([]Node, nodeInitCap),
		index:  make(map[string]map[string]int),
		rIndex: make(map[string]map[string][]int),
		sIndex: make(map[int]state),
		sCount: make(map[state]int),
		holes:  make([]int, secCap),
	}

	return set
}

func (ss *NodeSet) Add(nn ...Node) {
	var index int
	for _, n := range nn {
		if index = ss.indexOf(n); index != -1 {
			// skip existing
			continue
		}

		if len(ss.holes) == 0 {
			index = len(ss.set)
			ss.set = append(ss.set, n)
		} else {
			// filling holes created by Remove fn.
			index, ss.holes = ss.holes[0], ss.holes[1:]
			ss.set[index] = n
		}

		ss.sIndex[index] = pending
		ss.sCount[pending]++

		ss.reindex(n, index, nil, nil)

		// @todo update statuses of all referenced nodes to conflicting
	}
}

func (ss *NodeSet) Remove(nn ...Node) {
	var index int
	var oldState state
	for _, n := range nn {
		if index = ss.indexOf(n); index == -1 {
			// skip non existing
			continue
		}
		oldState = ss.sIndex[index]
		ss.sCount[oldState]--

		ss.set[index] = nil
		ss.sIndex[index] = pending

		ss.reindex(n, -1, n.Identifiers(), n.Relations())
		ss.holes = append(ss.holes, index)
	}
}

// Changes node state, adjusts counters and returns old state
func (ss *NodeSet) SetState(index int, s state) state {
	var (
		old = ss.sIndex[index]
	)

	if old != s {
		ss.sCount[old]--
		ss.sCount[pending]++
		ss.sIndex[index] = s
	}

	// @todo if new status is processed, update to all
	//       conflicting referenced nodes to pending

	return old
}

func (ss NodeSet) Count(s state) int {
	return ss.sCount[s]
}

// returns next node with compatible state and flips state of that node to processed
func (ss NodeSet) Shift(s state) Node {
	for i, n := range ss.set {
		if n == nil {
			// skip deleted
			continue
		}

		if ss.sIndex[i] != s {
			// skip nodes in different state
			continue
		}

		ss.SetState(i, processed)
		return n
	}

	return nil
}

// EachValid calls fn() for every valid node
func (ss NodeSet) Each(s state, fn func(Node, int)) {
	for i, n := range ss.set {
		if n == nil {
			continue
		}

		if ss.sIndex[i] != s {
			continue
		}

		fn(n, i)
	}
}

func (ss *NodeSet) reindex(n Node, index int, oldIdentifiers []string, oldRelations NodeRelationships) {
	var (
		r = n.Resource()
	)

	if _, has := ss.index[r]; !has {
		ss.index[r] = make(map[string]int)
	}

	if ss.rIndex == nil {
		ss.rIndex = make(map[string]map[string][]int)
	}

	for _, i := range oldIdentifiers {
		delete(ss.index[r], i)
	}

	// @todo cleanup oldRelations

	if index > -1 {

		for _, identifier := range n.Identifiers() {
			ss.index[r][identifier] = index
		}

		for r, ii := range n.Relations() {
			if _, has := ss.rIndex[r]; !has {
				ss.rIndex[r] = make(map[string][]int)
			}

			for _, i := range ii {

				has := false
				for _, e := range ss.rIndex[r][i] {
					if e == index {
						has = true
						break
					}
				}

				if !has {
					ss.rIndex[r][i] = append(ss.rIndex[r][i], index)
				}
			}
		}
	}
}

// Has checks if the given NodeSet contains a specific Node
func (ss NodeSet) Has(n Node) bool {
	if n == nil {
		return false
	}

	if i := ss.indexOf(n); i > -1 {
		return ss.set[i] != nil
		//m := ss.set[i]
		// @todo if it does not match, we'll probably have to recalculate
		//return match(n, m.Resource(), m.Identifiers()...)
	}

	return false
}

// Finds all matching nodes
//
// 1st check is done on index directly and then we recheck it directly on the node
func (ss NodeSet) filter(mask state, r string, ii ...string) []Node {
	nn := make([]Node, 0, len(ii))

	if len(ss.index) == 0 && len(ss.index[r]) == 0 {
		return nn
	}

	for _, i := range ii {
		index, has := ss.index[r][i]
		if !has {
			continue
		}

		if ss.set[index] == nil {
			// index pointing to removed node
			continue
		}

		if ss.sIndex[index]&mask > 0 {
			// index does not match
			continue
		}

		// recheck the node
		if match(ss.set[index], r, ii...) {
			nn = append(nn, ss.set[index])
		}
	}

	return nn
}

// Find all patent nodes
func (ss NodeSet) ParentsOf(n Node, mask state) []Node {
	parents := make([]Node, 0)

	if n == nil {
		panic("filtering with empty node")
	}

	if len(ss.rIndex) == 0 || len(ss.rIndex[n.Resource()]) == 0 {
		return parents
	}

	// iterate through all node's identifiers
	for _, fi := range n.Identifiers() {
		if len(ss.rIndex[n.Resource()][fi]) == 0 {
			continue
		}

		for _, index := range ss.rIndex[n.Resource()][fi] {
			if ss.sIndex[index]&mask > 0 {
				// index does not match
				continue
			}

			if ss.set[index] == nil {
				// removed
				continue
			}

			// match
			parents = append(parents, ss.set[index])
		}
	}

	return parents
}

// Finds all child nodes
func (ss NodeSet) ChildrenOf(n Node, mask state) []Node {
	children := make([]Node, 0)

	for r, ii := range n.Relations() {
		children = append(children, ss.filter(mask, r, ii...)...)
	}

	return children
}

// Finds node index
func (ss NodeSet) indexOf(n Node) int {
	var r = n.Resource()
	if len(ss.index) > 0 && len(ss.index[r]) > 0 {
		for _, i := range n.Identifiers() {
			if index, has := ss.index[r][i]; has {
				return index
			}
		}
	}

	return -1
}

func match(n Node, res string, ii ...string) bool {
	if matcher, is := n.(nodeMatcher); is {
		// Handle custom matchers
		return matcher.Matches(res, ii...)
	}

	return n.Resource() == res && n.Identifiers().HasAny(ii...)
}
