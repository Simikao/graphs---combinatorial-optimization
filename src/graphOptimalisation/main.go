package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	g "github.com/Simikao/graphOptimalisation/internal/graph"
)

// Eksport grafu do formatu .dot (Graphviz)
func LoadGraphFromFile(filename string, directed bool) (g.Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return g.Graph{}, err
	}
	defer file.Close()

	var edges [][]int
	var maxVertex int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		u, err1 := strconv.Atoi(parts[0])
		v, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return g.Graph{}, fmt.Errorf("błąd formatu pliku")
		}

		edges = append(edges, []int{u, v})
		if u > maxVertex {
			maxVertex = u
		}

		if v > maxVertex {
			maxVertex = v
		}
	}

	graph := g.NewGraph(maxVertex, directed)
	for _, edge := range edges {
		graph.AddEdge(edge[0], edge[1])
	}

	return graph, nil
}

func LoadGraphFromDotFile(filename string) (*g.Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var edges [][]int
	var maxVertex int
	directed := false
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Checking if the graph is directed
		if strings.HasPrefix(line, "digraph") {
			directed = true
		}

		if strings.Contains(line, "->") || strings.Contains(line, "--") {
			parts := strings.FieldsFunc(line, func(r rune) bool {
				return r == '-' || r == '>' || r == ';'
			})

			if len(parts) >= 2 {
				u, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
				v, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
				if err1 != nil || err2 != nil {
					return nil, fmt.Errorf("błąd formatu pliku .dot")
				}

				edges = append(edges, []int{u, v})
				if u > maxVertex {
					maxVertex = u
				}

				if v > maxVertex {
					maxVertex = v
				}
			}
		}
	}

	graph := g.NewGraph(maxVertex, directed)
	for _, edge := range edges {
		graph.AddEdge(edge[0], edge[1])
	}

	return &graph, nil
}

func main() {
	graph := g.NewGraphWithMatrix([][]int{
		{1, 1, 1},
		{1, 1, 0},
		{0, 0, 1},
	}, true)
	graph.UpdateEdges()
	graph.Inspect().
		RemoveEdge(3, 1)
	fmt.Println(graph.Edges)

	graph2 := g.NewGraph(4, false)
	graph2.AddEdge(2, 3).
		AddEdge(1, 3).
		AddEdge(2, 2).
		AddEdge(4, 1).
		Inspect().
		RemoveEdge(2, 3).
		Inspect()
	fmt.Println(graph2.Edges)
	graph2.AddEdge(2, 3)
	fmt.Println(graph2.Edges)
	graph2.RemoveEdge(2, 3).
		Inspect()
	fmt.Println(graph2.Edges)
	graph2.RemoveEdge(3, 1).
		Inspect()
	fmt.Println(graph2.Edges)

	fmt.Println("Vertexes")
	graph2.AddVertex().
		Inspect().
		AddEdge(3, 5).
		Inspect().
		RemoveVertex(5).
		Inspect().
		RemoveVertex(3).
		Inspect()

	graph.RemoveEdge(1, 3).AddEdge(3, 1).Inspect()
	fmt.Println(graph.GetInDegree(1))
	fmt.Println(graph.GetOutDegree(1))

	graph2.AddEdge(1, 2).AddEdge(2, 3)
	fmt.Println(graph2.GetDegree(1))
	fmt.Println(graph2.GetDegree(2))
	fmt.Println(graph2.GetDegree(3))
	fmt.Println(graph2.GetEvenOddDegreeCounts())
	fmt.Println(graph.GetMinMaxDegree())
	fmt.Println(graph2.GetMinMaxDegree())
	fmt.Println(graph.SortedByDegrees())

	graph3, err := LoadGraphFromDotFile("../../out/test.dot")
	if err != nil {
		panic(err)
	}
	graph3.Inspect()

	graphTestApproxSmall, err := LoadGraphFromFile("../../in/smallUndirected.txt", false)
	if err != nil {
		panic(err)
	}
	graphTestApproxDir, err := LoadGraphFromFile("../../in/exDirected.txt", true)
	if err != nil {
		panic(err)
	}
	graphTestApproxBig, err := LoadGraphFromFile("../../in/largeUndirected.txt", false)
	if err != nil {
		panic(err)
	}

	var logs string
	graphTestApproxSmall.Inspect().InspectEdges()
	result, err := graphTestApproxSmall.ApproximateVertexCover(&logs)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result)
	fmt.Println(logs)

	logs = ""
	graphTestApproxDir.Inspect().InspectEdges()
	result, err = graphTestApproxDir.ApproximateVertexCover(&logs)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result)
	fmt.Println(logs)

	logs = ""
	graphTestApproxBig.Inspect().InspectEdges()
	result, err = graphTestApproxBig.ApproximateVertexCover(&logs)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(result)
	fmt.Println(logs)
	// graph.ToDOT("test.dot")
}
