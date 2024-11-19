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
func LoadGraphFromFile(filename string, directed, weighted bool) (g.Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return g.Graph{}, err
	}
	defer file.Close()

	var edges [][]int
	var weights [][]float64
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
			return g.Graph{}, fmt.Errorf("invalid weight format")
		}

		edges = append(edges, []int{u, v})
		if u > maxVertex {
			maxVertex = u
		}

		if v > maxVertex {
			maxVertex = v
		}

		if weighted && len(parts) >= 3 {
			weight, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				return g.Graph{}, fmt.Errorf("invalid weight format")
			}
			weights = append(weights, []float64{float64(u), float64(v), weight})
		}
	}

	graph := g.NewGraph(maxVertex, directed, false)
	for _, edge := range edges {
		graph.AddEdge(edge[0], edge[1])
	}

	if weighted {
		for _, weight := range weights {
			err := graph.SetWeight(int(weight[0]), int(weight[1]), weight[2])
			if err != nil {
				log.Fatal(err)
			}
		}
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
	var weights [][]float64
	var maxVertex int
	directed := false
	weighted := false

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Checking if the graph is directed
		if strings.HasPrefix(line, "digraph") {
			directed = true
		}

		if strings.Contains(line, "->") || strings.Contains(line, "--") {
			parts := strings.FieldsFunc(line, func(r rune) bool {
				return r == '-' || r == '>' || r == ';' || r == '[' || r == ']'
			})

			if len(parts) >= 2 {
				u, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
				v, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
				if err1 != nil || err2 != nil {
					return nil, fmt.Errorf("invalid .dot file format")
				}

				edges = append(edges, []int{u, v})
				if u > maxVertex {
					maxVertex = u
				}

				if v > maxVertex {
					maxVertex = v
				}

				if strings.Contains(line, "label=") {
					weighted = true
					start := strings.Index(line, "\"") + 1
					end := strings.LastIndex(line, "\"")
					weight, err := strconv.ParseFloat(line[start:end], 64)
					if err != nil {
						return nil, fmt.Errorf("invalid weight format in .dot file")
					}
					weights = append(weights, []float64{float64(u), float64(v), weight})
				}
			}
		}
	}

	graph := g.NewGraph(maxVertex, directed, false)
	for _, edge := range edges {
		graph.AddEdge(edge[0], edge[1])
	}

	if weighted {
		for _, weight := range weights {
			err := graph.SetWeight(int(weight[0]), int(weight[1]), weight[2])
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return &graph, nil
}

func main() {
	// graph := g.NewGraphWithMatrix([][]int{
	// 	{1, 1, 1},
	// 	{1, 1, 0},
	// 	{0, 0, 1},
	// }, true)
	// graph.UpdateEdges()
	// graph.Inspect().
	// 	RemoveEdge(3, 1)
	// fmt.Println(graph.Edges)
	//
	// graph2 := g.NewGraph(4, false)
	// graph2.AddEdge(2, 3).
	// 	AddEdge(1, 3).
	// 	AddEdge(2, 2).
	// 	AddEdge(4, 1).
	// 	Inspect().
	// 	RemoveEdge(2, 3).
	// 	Inspect()
	// fmt.Println(graph2.Edges)
	// graph2.AddEdge(2, 3)
	// fmt.Println(graph2.Edges)
	// graph2.RemoveEdge(2, 3).
	// 	Inspect()
	// fmt.Println(graph2.Edges)
	// graph2.RemoveEdge(3, 1).
	// 	Inspect()
	// fmt.Println(graph2.Edges)
	//
	// fmt.Println("Vertexes")
	// graph2.AddVertex().
	// 	Inspect().
	// 	AddEdge(3, 5).
	// 	Inspect().
	// 	RemoveVertex(5).
	// 	Inspect().
	// 	RemoveVertex(3).
	// 	Inspect()
	//
	// graph.RemoveEdge(1, 3).AddEdge(3, 1).Inspect()
	// fmt.Println(graph.GetInDegree(1))
	// fmt.Println(graph.GetOutDegree(1))
	//
	// graph2.AddEdge(1, 2).AddEdge(2, 3)
	// fmt.Println(graph2.GetDegree(1))
	// fmt.Println(graph2.GetDegree(2))
	// fmt.Println(graph2.GetDegree(3))
	// fmt.Println(graph2.GetEvenOddDegreeCounts())
	// fmt.Println(graph.GetMinMaxDegree())
	// fmt.Println(graph2.GetMinMaxDegree())
	// fmt.Println(graph.SortedByDegrees())
	//
	// graph3, err := LoadGraphFromDotFile("../../out/test.dot")
	// if err != nil {
	// 	panic(err)
	// }
	// graph3.Inspect()

	// graphTestApproxSmall, err := LoadGraphFromFile("../../in/smallUndirected.txt", false)
	// if err != nil {
	// 	panic(err)
	// }
	// graphTestApproxDir, err := LoadGraphFromFile("../../in/exDirected.txt", true)
	// if err != nil {
	// 	panic(err)
	// }
	// graphTestApproxBig, err := LoadGraphFromFile("../../in/largeUndirected.txt", false)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// var logs string
	// graphTestApproxSmall.Inspect().InspectEdges()
	// result, err := graphTestApproxSmall.ApproximateVertexCover(&logs)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(result)
	// fmt.Println(logs)
	//
	// logs = ""
	// graphTestApproxDir.Inspect().InspectEdges()
	// result, err = graphTestApproxDir.ApproximateVertexCover(&logs)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(result)
	// fmt.Println(logs)
	//
	// logs = ""
	// graphTestApproxBig.Inspect().InspectEdges()
	// result, err = graphTestApproxBig.ApproximateVertexCover(&logs)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(result)
	// fmt.Println(logs)
	// graph.ToDOT("test.dot")

	// Test grafu nieskierowanego, ważonego
	undirectedGraph := g.NewGraph(3, false, true)
	undirectedGraph.AddEdge(1, 2, 4.5)
	undirectedGraph.AddEdge(1, 3, 3.0)
	undirectedGraph.AddEdge(2, 3, 2.5)
	fmt.Println("Undirected Weighted Graph:")
	fmt.Println(undirectedGraph.String())

	// Testowanie metod
	gottenWeight, err := undirectedGraph.GetWeight(1, 2)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println("Weight between 1 and 2:", gottenWeight)
	undirectedGraph.RemoveEdge(1, 3)
	fmt.Println("Graph after removing edge 1 -- 3:")
	fmt.Println(undirectedGraph.String())

	// Eksport do .dot
	err = undirectedGraph.ToDOT("undirected_weighted.dot")
	if err != nil {
		log.Fatalf("Error exporting to DOT: %v", err)
	}

	// Test grafu skierowanego, ważonego
	directedGraph := g.NewGraph(3, true, true)
	directedGraph.AddEdge(1, 2, 5.0)
	directedGraph.AddEdge(2, 3, 1.0)
	directedGraph.AddEdge(3, 1, 2.0)
	fmt.Println("Directed Weighted Graph:")
	fmt.Println(directedGraph.String())

	// Testowanie metod
	gottenWeight, err = directedGraph.GetWeight(3, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Weight between 3 and 1:", gottenWeight)
	err = directedGraph.SetWeight(3, 1, 2.5)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Updated weight between 3 and 1:")
	fmt.Println(directedGraph.GetWeight(3, 1))

	// Eksport do .dot
	err = directedGraph.ToDOT("directed_weighted.dot")
	if err != nil {
		log.Fatalf("Error exporting to DOT: %v", err)
	}
}
