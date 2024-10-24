package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
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

func (g *Graph) AddEdge(u, v int) {
	u = u - 1
	v = v - 1

	g.adjMatrix[u][v] = 1
	if !g.directed {
		g.adjMatrix[v][u] = 1
	}
}

func (g *Graph) RemoveEdge(u, v int) {
	u -= 1
	v -= 1

	g.adjMatrix[u][v] = 0
	if !g.directed {
		g.adjMatrix[v][u] = 0
	}
}

func (g *Graph) AddVertex() {
	for i := range g.adjMatrix {
		g.adjMatrix[i] = append(g.adjMatrix[i], 0)
	}
	newRow := make([]int, len(g.adjMatrix)+1)
	g.adjMatrix = append(g.adjMatrix, newRow)
}

func (g *Graph) RemoveVertex(v int) {
	v -= 1
	// check if the vertex exists
	if v >= len(g.adjMatrix) {
		return
	}

	// removing a row
	g.adjMatrix = append(g.adjMatrix[:v], g.adjMatrix[v+1:]...)

	// removing a column
	for i := range g.adjMatrix {
		g.adjMatrix[i] = append(g.adjMatrix[i][:v], g.adjMatrix[i][v+1:]...)
	}
}

func (g *Graph) GetOutDegree(v int) int {
	v--
	degree := 0
	for i := range g.adjMatrix[v] {
		degree += g.adjMatrix[v][i]
	}
	return degree
}

func (g *Graph) GetInDegree(v int) int {
	v--
	degree := 0
	for i := range g.adjMatrix {
		degree += g.adjMatrix[i][v]
	}
	return degree
}

func (g *Graph) GetDegree(v int) int {
	// given the fact that we pass this vertex to further function
	// we do not want to lower its value here to prevent doing it twice
	if g.directed {
		return g.GetInDegree(v) + g.GetOutDegree(v)
	}
	return g.GetOutDegree(v)
}

func (g *Graph) GetMinMaxDegree() (minDegree, maxDegree int) {
	minDegree = g.GetDegree(1)
	maxDegree = g.GetDegree(1)
	for i := 1; i <= len(g.adjMatrix); i++ {
		degree := g.GetDegree(i)
		if degree < minDegree {
			minDegree = degree
		}
		if degree > maxDegree {
			maxDegree = degree
		}
	}
	return minDegree, maxDegree
}

func (g *Graph) GetEvenOddDegreeCounts() (evenCount, oddCount int) {
	for i := 1; i <= len(g.adjMatrix); i++ {
		if g.GetDegree(i)%2 == 0 {
			evenCount++
		} else {
			oddCount++
		}
	}
	return evenCount, oddCount
}

func (g *Graph) SortedByDegrees() []int {
	degrees := make([]int, len(g.adjMatrix))
	for i := range g.adjMatrix {
		degrees[i] = g.GetDegree(i + 1)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(degrees)))
	return degrees
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
	path := fmt.Sprintf("../../out/%s", filename)

	file, err := os.Create(path)
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
				indexingFixI := i + 1
				indexingFixJ := j + 1
				if g.directed {
					fmt.Fprintf(file, "  %d -> %d;\n", indexingFixI, indexingFixJ)
				} else if i <= j { // aby uniknąć powielania krawędzi w grafie nieskierowanym
					fmt.Fprintf(file, "  %d -- %d;\n", indexingFixI, indexingFixJ)
				}
			}
		}
	}

	// Zakończenie
	fmt.Fprintf(file, "}\n")
	return nil
}

func LoadGraphFromFile(filename string, directed bool) (*Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
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
			return nil, fmt.Errorf("błąd formatu pliku")
		}
		edges = append(edges, []int{u, v})
		if u > maxVertex {
			maxVertex = u
		}
		if v > maxVertex {
			maxVertex = v
		}
	}

	graph := NewGraph(maxVertex, directed)

	for _, edge := range edges {
		graph.AddEdge(edge[0], edge[1])
	}

	return graph, nil
}

func LoadGraphFromDotFile(filename string) (*Graph, error) {
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

	graph := NewGraph(maxVertex, directed)

	for _, edge := range edges {
		graph.AddEdge(edge[0], edge[1])
	}

	return graph, nil
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

	fmt.Println(graph)
	graph2 := NewGraph(4, false)
	graph2.AddEdge(2, 3)
	graph2.AddEdge(1, 3)
	fmt.Println(graph2)
	graph2.RemoveEdge(2, 3)
	fmt.Println(graph2)

	fmt.Println("Vertexes")
	graph2.AddVertex()
	fmt.Println(graph2)
	graph2.AddEdge(3, 5)
	fmt.Println(graph2)
	graph2.RemoveVertex(5)
	fmt.Println(graph2)
	graph2.RemoveVertex(3)
	fmt.Println(graph2)

	graph.RemoveEdge(1, 3)
	graph.AddEdge(3, 1)
	fmt.Println(graph)
	fmt.Println(graph.GetInDegree(1))
	fmt.Println(graph.GetOutDegree(1))
	graph2.AddEdge(1, 2)
	graph2.AddEdge(2, 3)
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
	fmt.Println(graph3)
	// graph.ToDOT("test.dot")
}
