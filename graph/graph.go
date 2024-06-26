package graph

import (
	"fmt"

	model "github.com/akshaynanavare/zomato-mock/models"
	utils "github.com/akshaynanavare/zomato-mock/utils"
)

// Edge represents an edge in the graph
type Edge struct {
	Node   Node
	Weight float64
}

type Node struct {
	ID       string
	Location *model.Location
}

// Graph represents a graph using an adjacency list
type Graph struct {
	AdjacencyList map[string][]Edge
}

func (g *Graph) AddEdge(source, destination *Node, dist float64) {
	g.AdjacencyList[source.ID] = append(g.AdjacencyList[source.ID], Edge{
		Node:   *destination,
		Weight: dist,
	})

	g.AdjacencyList[destination.ID] = append(g.AdjacencyList[destination.ID], Edge{
		Node:   *source,
		Weight: dist,
	})
}

func (g *Graph) GetEdges(node string) []Edge {
	return g.AdjacencyList[node]
}

func (g *Graph) AddEdgeToUnvistitedNodes(nodes map[string]*Node, visited map[string]bool, source *Node) {
	fmt.Println("visited map : ", visited)
	for k, v := range nodes {
		if !visited[k] {
			for _, nbr := range g.GetEdges(v.ID) {
				if nbr.Node.ID == source.ID {
					continue
				}
			}
			g.AddEdge(source, v, utils.CalculateDistance(source.Location, v.Location))
		}
	}

	nodes[source.ID] = source
}

func (g *Graph) PrintGraph() {
	for src, node := range g.AdjacencyList {
		fmt.Print("source : ", src, " Nbr : ")
		for _, nbr := range node {
			fmt.Print(",", nbr)
		}

		fmt.Println()
	}
}
