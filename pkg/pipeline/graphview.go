package pipeline

type GraphView struct {
	Edges    []Edge      `json:"edges"`
	Nodes    []GraphNode `json:"nodes"`
	Position []float64   `json:"position"`
	Zoom     float64     `json:"zoom"`
	Viewport Viewport    `json:"viewport"`
}

func (gv *GraphView) ToPipelineGraph(nodeRegistry map[string]Node) *Graph {
	g := NewGraph()

	for _, gn := range gv.Nodes {
		node := gn.Data

		// set the node function
		if n, ok := nodeRegistry[node.Type]; ok {
			node.Run = n.Run
		}

		g.Nodes[node.Id] = &node
	}

	for _, e := range gv.Edges {
		g.Edges = append(g.Edges, e)

		inc, ok := g.Incoming[e.Target]
		if !ok {
			g.Incoming[e.Target] = []*Node{
				g.Nodes[e.Source],
			}
		} else {
			found := false
			for _, n := range inc {
				if n.Id == e.Source {
					found = true
					break
				}
			}
			if !found {
				g.Incoming[e.Target] = append(g.Incoming[e.Target], g.Nodes[e.Source])
			}
		}

		out, ok := g.Outgoing[e.Source]
		if !ok {
			g.Outgoing[e.Source] = []*Node{
				g.Nodes[e.Target],
			}
		} else {
			found := false
			for _, n := range out {
				if n.Id == e.Target {
					found = true
					break
				}
			}
			if !found {
				g.Outgoing[e.Source] = append(g.Outgoing[e.Source], g.Nodes[e.Target])
			}
		}
	}

	return g
}

type Viewport struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Zoom float64 `json:"zoom"`
}

type GraphNode struct {
	Id       string
	Type     string
	Position Position
	Data     Node `json:"data"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
