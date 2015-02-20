package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var sccSize int
var leader int64 = -1
var finishingTimes []int64

type graph struct {
	sccs  []int
	edges map[int64][]int64
	nodes map[int64]*vertex
}

type vertex struct {
	explored bool
}

func (g *graph) addEdge(tail int64, head int64) {
	g.edges[tail] = append(g.edges[tail], head)

	g.nodes[tail] = &vertex{explored: false}
	g.nodes[head] = &vertex{explored: false}
}

func loadGraph(filename string, g *graph, ginv *graph) error {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		numbers := strings.Split(scanner.Text(), " ")

		tail, err := strconv.ParseInt(numbers[0], 10, 0)
		if err != nil {
			log.Fatal(err)
		}

		head, err := strconv.ParseInt(numbers[1], 10, 0)
		if err != nil {
			log.Fatal(err)
		}

		g.addEdge(tail, head)
		ginv.addEdge(head, tail)
	}

	return scanner.Err()
}

func depthFirstSearchLoop(g *graph) {
	for index, node := range g.nodes {
		if node.explored == false {
			depthFirstSearch(index, g)
		}
	}
}

func depthFirstSearch(tail int64, g *graph) {
	g.nodes[tail].explored = true

	if leader >= 0 {
		sccSize++
	}

	for _, head := range g.edges[tail] {
		if g.nodes[head].explored == false {
			depthFirstSearch(head, g)
		}
	}

	finishingTimes = append(finishingTimes, tail)
}

func kosaraju(g *graph) {
	for i := len(finishingTimes) - 1; i >= 0; i-- {
		node := finishingTimes[i]
		if g.nodes[node].explored == false {
			leader = node
			depthFirstSearch(node, g)
		}
		g.sccs = append(g.sccs, sccSize)
		sccSize = 0
	}
}

func main() {

	var g, ginv graph

	// Initialize maps
	g.nodes = make(map[int64]*vertex, 1000000)
	g.edges = make(map[int64][]int64, 1000000)

	ginv.edges = make(map[int64][]int64, 1000000)
	ginv.nodes = make(map[int64]*vertex, 1000000)

	log.Println("Loading graphs...")
	loadGraph("SCC.txt", &g, &ginv)

	log.Println("First pass...")
	depthFirstSearchLoop(&ginv)

	log.Println("Second pass...")
	kosaraju(&g)

	log.Println("Sorting strongly connected component sizes...")
	sort.Sort(sort.Reverse(sort.IntSlice(g.sccs)))

	fmt.Printf("SCCS: %v", g.sccs[:5])
}
