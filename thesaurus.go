package main

type Edge struct {
	name   string
	weight float64
	target *Node
}

type Node struct {
	edges       map[string]Edge
	name        string
	start       int
	end         int
	weight      float64
	scentenceId uint32
}

func makeNode(name string, scentenceId uint32) *Node {
	node := &Node{name: name, edges: make(map[string]Edge), start: 0, end: 0, weight: 0, scentenceId: scentenceId}
	return node
}

type Thesaurus struct {
	pntrmap map[string]*Node
}

func (thes *Thesaurus) addEntry(start string, targets []string, scentenceId uint32) {
	_, exists := thes.pntrmap[start]
	if !exists {
		thes.pntrmap[start] = makeNode(start, scentenceId)
	}
	thes.pntrmap[start].weight += 1
	thes.addEdges(start, targets, scentenceId)
}

func (thes *Thesaurus) addEdges(start string, targets []string, scentenceId uint32) {
	val := thes.pntrmap[start]
	for _, s := range targets {
		edgeVal, edgeExists := val.edges[s]
		if edgeExists {
			edgeVal.weight += 1
			val.edges[s] = edgeVal
		} else {
			targetVal, targetExists := thes.pntrmap[s]
			if !targetExists {
				thes.pntrmap[s] = makeNode(s, scentenceId)
				targetVal = thes.pntrmap[s]
			}
			val.edges[s] = Edge{name: s, weight: 1, target: targetVal}
		}
	}
}
