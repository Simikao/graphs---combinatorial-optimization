package main

import (
	"fmt"
	"os"
	"strings"
	// "github.com/jesseduffield/lazygit/pkg/gui/presentation/graph"
)

type Graph struct {
	adjMatrix [][]int
	directed  bool
}

func NewGraph(vertices int, directed bool) *Graph {
	matrix := make([][]int, vertices)
	for i := range matrix {
		matrix[i] = make([]int, vertices)
	}
	return &Graph{
		adjMatrix: matrix,
		directed:  directed,
	}
}

func (g *Graph) String() string {
	var output []string

	// Add the header row
	header := []string{"  "}
	for i := 1; i <= len(g.adjMatrix); i++ {
		header = append(header, fmt.Sprintf("%d", i))
	}
	output = append(output, strings.Join(header, " "))

	// Add each row of the matrix
	for i, row := range g.adjMatrix {
		rowOutput := []string{fmt.Sprintf("%d", i+1), fmt.Sprintf("%v", row)}
		output = append(output, strings.Join(rowOutput, " "))
	}

	return strings.Join(output, "\n")
}

// Eksport grafu do formatu .dot (Graphviz)
func (g *Graph) ToDOT(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Nagłówek grafu
	var graphType string
	if g.directed {
		graphType = "digraph"
	} else {
		graphType = "graph"
	}

	fmt.Fprintf(file, "%s G {\n", graphType)

	// Dodanie krawędzi
	for i := range g.adjMatrix {
		for j := range g.adjMatrix[i] {
			if g.adjMatrix[i][j] == 1 {
				if g.directed {
					fmt.Fprintf(file, "  %d -> %d;\n", i, j)
				} else if i <= j { // aby uniknąć powielania krawędzi w grafie nieskierowanym
					fmt.Fprintf(file, "  %d -- %d;\n", i, j)
				}
			}
		}
	}

	// Zakończenie
	fmt.Fprintf(file, "}\n")
	return nil
}

func main() {
	graph := &Graph{
		adjMatrix: [][]int{
			{1, 1, 1},
			{1, 1, 0},
			{0, 0, 1},
		},
		directed: true,
	}

	graph2 := NewGraph(2, false)
	fmt.Println(graph2)
	fmt.Println(graph)
	graph.ToDOT("test.dot")
}
