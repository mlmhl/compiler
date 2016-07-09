package regex

type dfaNode struct {
	groupIds []int
	next     map[int]*dfaNode
}

type dfaGraph struct {
	startNode *dfaNode
}

// return the last position and group ids
func (graph *dfaGraph) match(reText string) (int, []int) {
	node := graph.startNode
	for i, c := range reText {
		if nextNode, ok := node.next[(int)(c)]; ok {
			node = nextNode
		} else {

			if arbitraryNode, ok := node.next[ARBITRARY_ID]; ok {
				node = arbitraryNode
			} else {
				return i, node.groupIds
			}
		}
	}
	return len(reText), node.groupIds
}
