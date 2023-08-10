package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

func drawNode(x graph.Graph[string, string], node *Node, drawn *[]string, limit int, weightLimit float64) graph.Graph[string, string] {
	if limit <= 0 || len(*drawn) > 1000 {
		return x
	}
	//fmt.Println(limit)
	y := x
	//fmt.Println(node.name)
	if !contains(drawn, node.name) {
		_ = x.AddVertex(node.name, graph.VertexAttribute("label", node.name))
		*drawn = append(*drawn, node.name)
	}
	edgesCounter := 0
	for _, e := range node.edges {
		if e.target.weight < weightLimit {
			continue
		}
		edgesCounter++
		if edgesCounter >= 10 {
			break
		}
		edgeIsDrawn := contains(drawn, e.target.name)
		if !edgeIsDrawn {
			_ = x.AddVertex(e.target.name, graph.VertexAttribute("label", e.target.name))
			*drawn = append(*drawn, e.target.name)
		}
		_ = x.AddEdge(node.name, e.target.name)

		y = drawNode(y, e.target, drawn, limit-1, weightLimit)
	}
	return y
}

var thes = Thesaurus{pntrmap: make(map[string]*Node)}

func main() {

	records := readCsvFile("./csv_file2.csv")
	for _, record := range records {
		title := strings.ToLower(record[2])
		words := strings.Split(title, " ")
		scentenceId := hash(title)
		for i := 0; i < len(words)-1; i++ {
			thes.addEntry(trim(words[i]), []string{trim(words[i+1])}, scentenceId)
		}
	}
	for _, node := range thes.pntrmap {
		sum := 0.0
		node.weight /= float64(len(thes.pntrmap))
		for _, edge := range node.edges {
			sum += edge.weight
		}
		for _, edge := range node.edges {
			edge.weight = edge.weight / sum
			//fmt.Println(edge.weight, node.name, edge.name)
		}
		//fmt.Println(node.weight)
	}

	g := graph.New(graph.StringHash, graph.Directed())
	drawn := []string{}
	//ctr := 0
	//for _, node := range thes.pntrmap {
	//	if ctr >= 4 {
	//		break
	//	}
	//	ctr++
	//
	//	fmt.Println(node.name)
	//	g = drawNode(g, node, &drawn, 4)
	//}
	g = drawNode(g, thes.pntrmap["the"], &drawn, 4, 0.3)
	file, _ := os.Create("my-graph.gv")
	_ = draw.DOT(g, file)

	http.ListenAndServe("127.0.0.1:8080", http.HandlerFunc(Serve))
}
