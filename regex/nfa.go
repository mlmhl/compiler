package regex

import (
	"sort"

	"github.com/mlmhl/goutil/container"
)

type nfaNode struct {
	c       int
	groupId int

	isEnd bool
}

// Set of nfa nodes
type identifier struct {
	set []int
}

func newIdentifier(set []int) *identifier {
	sort.Ints(set)
	return &identifier{set: set}
}

func (id *identifier) HashCode() uint32 {
	var code uint32 = 0
	for _, i := range id.set {
		code ^= (uint32)(i)
		code <<= 1
	}
	return code
}

func (id *identifier) Equal(other container.Hashable) bool {
	if otherId, ok := other.(*identifier); !ok {
		return false
	} else {
		if len(id.set) != len(otherId.set) {
			return false
		}
		for i, v := range id.set {
			if v != otherId.set[i] {
				return false
			}
		}
		return true
	}
}

func newNfaNode(c, groupId int, isEnd bool) *nfaNode {
	return &nfaNode{
		c:       c,
		groupId: groupId,

		isEnd: isEnd,
	}
}

type nfaGraph struct {
	nodes []*nfaNode
	edges map[int][]int
}

func newNfaGraph() *nfaGraph {
	return &nfaGraph{
		nodes: []*nfaNode{&nfaNode{}}, // Add a common start node.
		edges: map[int][]int{},
	}
}

func (graph *nfaGraph) size() int {
	return len(graph.nodes)
}

func (graph *nfaGraph) addRegexExpression(reText string, groupId int) error {
	start := graph.size()
	if err := graph.parseRegexExpression(reText, groupId); err != nil {
		return err
	}
	graph.addEdge(0, start)
	graph.compile(start, graph.size())
	return nil
}

func (graph *nfaGraph) parseRegexExpression(reText string, groupId int) error {
	if len(reText) == 0 {
		return nil
	}

	normalized, err := normalizeRegex(reText)
	if err != nil {
		return err
	}

	for i := 0; i < len(normalized)-1; i++ {
		graph.addNode(normalized[i], groupId, false)
	}
	graph.addNode(normalized[len(normalized)-1], groupId, true)

	return nil
}

func (graph *nfaGraph) compile(start, end int) {
	nodes := graph.nodes
	ops := container.NewStack()

	for i := start; i < end; i++ {
		lsp := i
		node := nodes[i]

		if node.c == LSP_ID || node.c == CHOICE_ID {
			ops.Push(i)
		} else if node.c == RSP_ID {
			choices := []int{}
			for {
				op := ops.Peek().(int)
				ops.Pop()
				if nodes[op].c == CHOICE_ID {
					choices = append(choices, op)
				} else {
					lsp = op
					break
				}
			}
			for _, choice := range choices {
				graph.addEdge(lsp, choice+1)
				graph.addEdge(choice, i)
			}
		}

		if i < graph.size()-1 {
			c := nodes[i+1].c
			if c == REPETITION_ID {
				graph.addEdge(lsp, i+1)
				graph.addEdge(i+1, lsp)
			} else if c == ZERO_OR_ONE_ID {
				graph.addEdge(lsp, i+1)
			} else if c == ONE_OR_MORE_ID {
				graph.addEdge(i+1, lsp)
			}
		}

		if (node.c == LSP_ID || node.c == REPETITION_ID || node.c == ZERO_OR_ONE_ID ||
			node.c == ONE_OR_MORE_ID || node.c == RSP_ID) && i < graph.size()-1 {
			graph.addEdge(i, i+1)
		}
	}
}

func (graph *nfaGraph) addNode(c, groupId int, isEnd bool) {
	graph.nodes = append(graph.nodes, newNfaNode(c, groupId, isEnd))
}

func (graph *nfaGraph) addEdge(from, to int) {
	graph.edges[from] = append(graph.edges[from], to)
}

func (graph *nfaGraph) toDfa() *dfaGraph {
	nfaNodes := graph.nodes
	unmarked := container.NewQueue()
	dfaNodes := container.NewHashMap()

	start := newIdentifier(graph.getClosure([]int{0}))
	startNode := graph.newDfaNode(start)
	dfaNodes.Put(start, startNode)
	unmarked.Push(start)

	for unmarked.Len() > 0 {
		id := unmarked.Front().(*identifier)
		unmarked.Pop()

		next := map[int][]int{}
		for _, i := range id.set {
			if i == len(nfaNodes) - 1 {
				continue
			}
			c := nfaNodes[i].c
			if !isMataSymbol(c) || c == ARBITRARY_ID {
				if v, ok := next[c]; ok {
					next[c] = append(v, i+1)
				} else {
					next[c] = []int{i+1}
				}
			}
		}

		node := dfaNodes.Get(id).(*dfaNode)
		for c, v := range next {
			nId := newIdentifier(graph.getClosure(v))
			var nNode *dfaNode
			v := dfaNodes.Get(nId)
			if v == nil {
				nNode = graph.newDfaNode(nId)
				unmarked.Push(nId)
				dfaNodes.Put(nId, nNode)
			} else {
				nNode = v.(*dfaNode)
			}
			node.next[c] = nNode
		}
	}

	return &dfaGraph{startNode}
}

func (graph *nfaGraph) newDfaNode(id *identifier) *dfaNode {
	groupIds := []int{}
	for _, i := range id.set {
		if graph.nodes[i].isEnd {
			groupIds = append(groupIds, graph.nodes[i].groupId)
		}
	}
	sort.Ints(groupIds)
	return &dfaNode{groupIds, map[int]*dfaNode{}}
}

func (graph *nfaGraph) getClosure(candidates []int) []int {
	marked := map[int]bool{}

	for _, i := range candidates {
		if _, ok := marked[i]; !ok {
			graph.dfs(i, marked)
		}
	}

	res := []int{}
	for i, _ := range marked {
		res = append(res, i)
	}
	return res
}

func (graph *nfaGraph) dfs(i int, marked map[int]bool) {
	marked[i] = true
	for _, next := range graph.edges[i] {
		if _, ok := marked[next]; !ok {
			graph.dfs(next, marked)
		}
	}
}

// for test
// return the last position and group ids
func (graph *nfaGraph) match(text string) (int, []int) {
	if graph.size() == 0 {
		return 0, []int{}
	}

	nodes := graph.nodes

	pos := len(text)
	candidates := graph.getClosure([]int{0})

	for p, c := range text {
		matched := []int{}
		for _, i := range candidates {
			if nodes[i].c == (int)(c) || nodes[i].c == ARBITRARY_ID {
				matched = append(matched, i+1)
			}
		}

		if len(matched) == 0 {
			pos = p
			break
		} else {
			candidates = graph.getClosure(matched)
		}
	}

	groups := []int{}
	for _, i := range candidates {
		if nodes[i].isEnd {
			groups = append(groups, nodes[i].groupId)
		}
	}

	return pos, groups
}
