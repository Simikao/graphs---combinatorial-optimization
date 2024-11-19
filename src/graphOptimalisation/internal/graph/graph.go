package graph

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

var ErrDirectedGraph = errors.New("Cannot use a directed graph in this algorithm")

type Graph struct {
	AdjMatrix    [][]int
	WeightMatrix [][]float64
	Directed     bool
	Weighted     bool
	Edges        [][2]int
}

func getEdges(vertices [][]int, directed bool) [][2]int {
	var edges [][2]int
	for i := range vertices {
		for j := range vertices[i] {
			if vertices[i][j] == 1 {
				if directed {
					edges = append(edges, [2]int{i + 1, j + 1})
				} else {
					if i <= j {
						edges = append(edges, [2]int{i + 1, j + 1})
					}
				}
			}
		}
	}
	return edges
}

func (g *Graph) SetWeight(u, v int, weight float64) error {
	if !g.Weighted {
		return fmt.Errorf("cannot set weight on an unweighted graph")
	}

	u, v = u-1, v-1
	g.WeightMatrix[u][v] = weight
	if !g.Directed {
		g.WeightMatrix[v][u] = weight
	}
	return nil
}

func (g *Graph) GetWeight(u, v int) (float64, error) {
	if !g.Weighted {
		return 0, fmt.Errorf("graph is unweighted")
	}

	u, v = u-1, v-1
	return g.WeightMatrix[u][v], nil
}

func NewGraphWithMatrix(vertices [][]int, directed bool) Graph {
	edges := getEdges(vertices, directed)
	return Graph{
		AdjMatrix: vertices,
		Directed:  directed,
		Edges:     edges,
	}
}

func NewGraph(vertices int, directed, weighted bool) Graph {
	matrix := make([][]int, vertices)
	for i := range matrix {
		matrix[i] = make([]int, vertices)
	}
	edges := getEdges(matrix, directed)

	var weights [][]float64
	if weighted {
		weights = make([][]float64, vertices)
		for i := range weights {
			weights[i] = make([]float64, vertices)
		}
	}

	return Graph{
		AdjMatrix:    matrix,
		WeightMatrix: weights,
		Directed:     directed,
		Weighted:     weighted,
		Edges:        edges,
	}
}

func internalNewGraph(matrix [][]int, directed bool, edges [][2]int) Graph {
	return Graph{
		AdjMatrix: matrix,
		Directed:  directed,
		Edges:     edges,
	}
}

func (g *Graph) UpdateEdges() {
	g.Edges = getEdges(g.AdjMatrix, g.Directed)
}

func (g *Graph) AddEdge(u, v int, weight ...float64) *Graph {
	u = u - 1
	v = v - 1

	g.AdjMatrix[u][v] = 1
	if !g.Directed {
		g.AdjMatrix[v][u] = 1
	}

	if g.Weighted && len(weight) > 0 {
		g.WeightMatrix[u][v] = weight[0]
		if !g.Directed {
			g.WeightMatrix[v][u] = weight[0]
		}
	}

	// Update the list of edges with new edge
	g.Edges = append(g.Edges, [2]int{u + 1, v + 1})
	return g
}

func (g *Graph) RemoveEdge(u, v int) *Graph {
	u -= 1
	v -= 1

	g.AdjMatrix[u][v] = 0
	if !g.Directed {
		g.AdjMatrix[v][u] = 0
	}
	// Remove the edge from edge slice as well
	for i, edge := range g.Edges {
		if (edge[0] == u+1 && edge[1] == v+1) || (!g.Directed && edge[0] == v+1 && edge[1] == u+1) {
			g.Edges = append(g.Edges[:i], g.Edges[i+1:]...)
			break
		}
	}
	return g
}

func (g *Graph) AddVertex() *Graph {
	for i := range g.AdjMatrix {
		g.AdjMatrix[i] = append(g.AdjMatrix[i], 0)
	}
	newRow := make([]int, len(g.AdjMatrix)+1)
	g.AdjMatrix = append(g.AdjMatrix, newRow)
	return g
}

func (g *Graph) RemoveVertex(v int) *Graph {
	v -= 1
	// check if the vertex exists
	if v >= len(g.AdjMatrix) {
		return g
	}

	// removing a row
	g.AdjMatrix = append(g.AdjMatrix[:v], g.AdjMatrix[v+1:]...)

	// removing a column
	for i := range g.AdjMatrix {
		g.AdjMatrix[i] = append(g.AdjMatrix[i][:v], g.AdjMatrix[i][v+1:]...)
	}

	// Update edges of the graph by removing all edges with the removed vertex from the list
	var updatedEdges [][2]int
	for _, edge := range g.Edges {
		if edge[0] != v+1 && edge[1] != v+1 {
			updatedEdges = append(updatedEdges, edge)
		}
	}
	g.Edges = updatedEdges

	return g
}

func (g *Graph) GetOutDegree(v int) int {
	v--
	degree := 0
	for i := range g.AdjMatrix[v] {
		degree += g.AdjMatrix[v][i]
	}
	return degree
}

func (g *Graph) GetInDegree(v int) int {
	v--
	degree := 0
	for i := range g.AdjMatrix {
		degree += g.AdjMatrix[i][v]
	}
	return degree
}

func (g *Graph) GetDegree(v int) int {
	// given the fact that we pass this vertex to further function
	// we do not want to lower its value here to prevent doing it twice
	if g.Directed {
		return g.GetInDegree(v) + g.GetOutDegree(v)
	}
	return g.GetOutDegree(v)
}

func (g *Graph) GetMinMaxDegree() (minDegree, maxDegree int) {
	minDegree = g.GetDegree(1)
	maxDegree = g.GetDegree(1)
	for i := 1; i <= len(g.AdjMatrix); i++ {
		degree := g.GetDegree(i)
		if degree < minDegree {
			minDegree = degree
		}
		if degree > maxDegree {
			maxDegree = degree
		}
	}
	return
}

func (g *Graph) GetEvenOddDegreeCounts() (evenCount, oddCount int) {
	for i := 1; i <= len(g.AdjMatrix); i++ {
		if g.GetDegree(i)%2 == 0 {
			evenCount++
		} else {
			oddCount++
		}
	}
	return
}

func (g *Graph) SortedByDegrees() []int {
	degrees := make([]int, len(g.AdjMatrix))
	for i := range g.AdjMatrix {
		degrees[i] = g.GetDegree(i + 1)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(degrees)))
	return degrees
}

func (g *Graph) ApproximateVertexCover(logs *string) ([]int, error) {
	errWrap := func(err error) error {
		return fmt.Errorf("ApproximateVertexCover: %w", err)
	}
	var doLogs bool
	if logs != nil {
		doLogs = true
	}

	if g.Directed {
		return nil, errWrap(ErrDirectedGraph)
	}
	var log strings.Builder
	edgesInternal := g.Edges
	cover := make(map[int]struct{})

	// Step 1: while there are edges in graph...
	for len(edgesInternal) > 0 {
		// pick an edge
		edge := edgesInternal[0]
		u, v := edge[0], edge[1]
		// Step 2: save two vertices of the chosen edge
		cover[u] = struct{}{}
		cover[v] = struct{}{}

		if doLogs {
			log.WriteString(fmt.Sprintf("Adding vertices %d and %d to the cover and removing them from the clone.\n", u, v))
		}

		// Step 3: remove both endpoints and their adjecent edges
		var updatedEdges [][2]int
		for _, e := range edgesInternal {
			if e[0] != u && e[0] != v && e[1] != u && e[1] != v {
				updatedEdges = append(updatedEdges, e)
			}
		}

		edgesInternal = updatedEdges
		if doLogs {
			log.WriteString(fmt.Sprintf("Current cover set: %v\n", cover))
			log.WriteString(fmt.Sprintf("Remaining edges after removal: %v\n", edgesInternal))
		}

	}

	result := make([]int, 0, len(cover))
	for vertex := range cover {
		result = append(result, vertex)
	}
	sort.Ints(result)

	if doLogs {
		log.WriteString(fmt.Sprintf("Approximate Vertex Cover: %v\n", result))
		*logs = log.String()
	}
	return result, nil
}

func (g *Graph) String() string {
	var sb strings.Builder

	// Nagłówek z indeksami kolumn
	sb.WriteString("  ") // Puste miejsce dla indeksów wierszy
	for i := 1; i <= len(g.AdjMatrix); i++ {
		sb.WriteString(fmt.Sprintf(" %d", i))
	}
	sb.WriteString("\n")

	// Macierz sąsiedztwa z indeksami wierszy
	for i, row := range g.AdjMatrix {
		sb.WriteString(fmt.Sprintf("%d ", i+1)) // Indeks wiersza
		for _, val := range row {
			sb.WriteString(fmt.Sprintf(" %d", val))
		}
		sb.WriteString("\n")
	}

	// Jeśli graf jest ważony, dodaj macierz wag
	if g.Weighted {
		sb.WriteString("\nWeight Matrix:\n")
		sb.WriteString(" ") // Puste miejsce dla indeksów wierszy
		for i := 1; i <= len(g.WeightMatrix); i++ {
			sb.WriteString(fmt.Sprintf("%5d", i))
		}
		sb.WriteString("\n")

		for i, row := range g.WeightMatrix {
			sb.WriteString(fmt.Sprintf("%d ", i+1)) // Indeks wiersza
			for _, val := range row {
				sb.WriteString(fmt.Sprintf("%5.1f", val)) // Wagi w formacie dziesiętnym
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (g *Graph) ToDOT(filename string) error {
	path := fmt.Sprintf("../../out/%s", filename)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Ustaw nagłówek grafu w zależności od typu (skierowany/nieskierowany)
	var graphType string
	var edgeConnector string
	if g.Directed {
		graphType = "digraph"
		edgeConnector = "->"
	} else {
		graphType = "graph"
		edgeConnector = "--"
	}

	fmt.Fprintf(file, "%s G {\n", graphType)

	// Dodanie krawędzi
	for i := range g.AdjMatrix {
		for j := range g.AdjMatrix[i] {
			if g.AdjMatrix[i][j] == 1 {
				indexingFixI := i + 1
				indexingFixJ := j + 1
				if g.Directed || i <= j {
					if g.Weighted {
						fmt.Fprintf(file, "  %d %s %d [label=\"%.2f\"];\n", indexingFixI, edgeConnector, indexingFixJ, g.WeightMatrix[i][j])
					} else {
						fmt.Fprintf(file, "  %d %s %d;\n", indexingFixI, edgeConnector, indexingFixJ)
					}
				}
			}
		}
	}

	// Zakończenie
	fmt.Fprintf(file, "}\n")
	return nil
}

func (g *Graph) Inspect() *Graph {
	fmt.Println(g)
	return g
}

func (g *Graph) InspectEdges() *Graph {
	fmt.Println(g.Edges)
	return g
}
