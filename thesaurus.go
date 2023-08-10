package main

type Edge struct {
	name        string
	weight      float64
	target      *Node
	scentenceId uint32
}

type Node struct {
	edges  map[string]map[uint32]Edge
	name   string
	start  int
	end    int
	weight float64
}

func makeNode(name string) *Node {
	node := &Node{name: name, edges: make(map[string]map[uint32]Edge), start: 0, end: 0, weight: 0}
	return node
}

type Thesaurus struct {
	pntrmap map[string]*Node
}

func (thes *Thesaurus) addEntry(start string, targets []string, scentenceId uint32) {
	_, exists := thes.pntrmap[start]
	if !exists {
		thes.pntrmap[start] = makeNode(start)
	}
	thes.pntrmap[start].weight += 1
	thes.addEdges(start, targets, scentenceId)
}

func (thes *Thesaurus) addEdges(start string, targets []string, scentenceId uint32) {
	val := thes.pntrmap[start]
	for _, s := range targets {
		edgeVal, edgeExists := val.edges[s]
		if edgeExists {
			edge := edgeVal[scentenceId]
			edge.weight += 1
			edgeVal[scentenceId] = edge
			val.edges[s] = edgeVal
		} else {
			targetVal, targetExists := thes.pntrmap[s]
			if !targetExists {
				thes.pntrmap[s] = makeNode(s)
				targetVal = thes.pntrmap[s]
			}
			val.edges[s][scentenceId] = Edge{name: s, weight: 1, target: targetVal, scentenceId: scentenceId}
		}
	}
}
