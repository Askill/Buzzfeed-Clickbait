package main

import "fmt"

func test() {
	thes := Thesaurus{pntrmap: make(map[string]*Node)}
	thes.addEntry("node1", []string{"node2"}, 1)
	thes.addEntry("node2", []string{"node2", "node1"}, 2)
	thes.addEntry("node3", []string{"node1"}, 2)
	thes.addEntry("node3", []string{"node1"}, 3)
	for i, s := range thes.pntrmap {
		fmt.Println(i, s.edges)
	}
}
