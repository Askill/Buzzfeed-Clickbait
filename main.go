package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

type Edge struct {
	name   string
	weight float64
	target *Node
}

type Node struct {
	edges map[string]Edge
	name  string
}

func makeNode(name string) *Node {
	node := &Node{name: name, edges: make(map[string]Edge)}
	return node
}

type Thesaurus struct {
	pntrmap map[string]*Node
}

func (thes *Thesaurus) addEntry(start string, targets []string) {
	_, exists := thes.pntrmap[start]
	if !exists {
		thes.pntrmap[start] = makeNode(start)
	}
	thes.addEdges(start, targets)
}

func (thes *Thesaurus) addEdges(start string, targets []string) {
	val := thes.pntrmap[start]
	for _, s := range targets {
		edgeVal, edgeExists := val.edges[s]
		if edgeExists {
			edgeVal.weight += 1
			val.edges[s] = edgeVal
		} else {
			targetVal, targetExists := thes.pntrmap[s]
			if !targetExists {
				thes.pntrmap[s] = makeNode(s)
				targetVal = thes.pntrmap[s]
			}
			val.edges[s] = Edge{name: s, weight: 1, target: targetVal}
		}
	}
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func test() {
	thes := Thesaurus{pntrmap: make(map[string]*Node)}
	thes.addEntry("node1", []string{"node2"})
	thes.addEntry("node2", []string{"node2", "node1"})
	thes.addEntry("node3", []string{"node1"})
	thes.addEntry("node3", []string{"node1"})
	for i, s := range thes.pntrmap {
		fmt.Println(i, s.edges)
	}
}

func printThesaurus(thes Thesaurus) {
	for i, s := range thes.pntrmap {
		fmt.Println(i, s.edges)
		fmt.Println("####################################################################")
	}
}

func trim(word string) string {
	word1 := strings.Replace(word, "\"", "", -1)
	word1 = strings.Replace(word1, "'", "", -1)

	word1 = strings.Replace(word1, ".", "", -1)
	word1 = strings.Replace(word1, "!", "", -1)
	word1 = strings.Replace(word1, "?", "", -1)
	word1 = strings.Replace(word1, ":", " ", -1)
	word1 = strings.Replace(word1, "#", " ", -1)
	word1 = strings.Replace(word1, ",", " ", -1)
	word1 = strings.Replace(word1, "(", " ", -1)
	word1 = strings.Replace(word1, ")", " ", -1)
	word1 = strings.Replace(word1, "”", " ", -1)
	word1 = strings.Replace(word1, "“", " ", -1)
	word1 = strings.Replace(word1, " ", "", -1)
	return word1
}
func contains(slc *[]string, str string) bool {
	for _, x := range *slc {
		if x == str {
			return true
		}
	}
	return false
}

func drawNode(x graph.Graph[string, string], node *Node, drawn *[]string, limit int) graph.Graph[string, string] {
	if limit <= 0 || len(*drawn) > 100 {
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
		if e.weight < 0.001 {
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

		y = drawNode(y, e.target, drawn, limit-1)
	}
	return y
}

func main() {
	thes := Thesaurus{pntrmap: make(map[string]*Node)}

	records := readCsvFile("./csv_file2.csv")
	for _, record := range records {
		title := record[2]
		words := strings.Split(title, " ")
		for i := 0; i < len(words)-1; i++ {
			thes.addEntry(trim(words[i]), []string{trim(words[i+1])})
		}
	}
	for _, node := range thes.pntrmap {
		sum := 0.0
		for _, edge := range node.edges {
			sum += edge.weight
		}
		for _, edge := range node.edges {
			edge.weight = edge.weight / sum
			//fmt.Println(edge.weight, node.name, edge.name)
		}
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
	g = drawNode(g, thes.pntrmap["10"], &drawn, 13)
	file, _ := os.Create("my-graph.gv")
	_ = draw.DOT(g, file)
}
