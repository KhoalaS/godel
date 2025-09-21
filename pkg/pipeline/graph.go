package pipeline

type Graph struct {
	Edges    []Edge             `json:"edges"`
	Nodes    map[string]*Node   `json:"nodes"`
	Incoming map[string][]*Node `json:"incoming"`
	Outgoing map[string][]*Node `json:"outgoing"`
}

type Edge struct {
	Id           string `json:"id"`
	Target       string `json:"target"`
	Source       string `json:"source"`
	SourceHandle string `json:"sourceHandle"`
	TargetHandle string `json:"targetHandle"`
	Label        string `json:"label"`
}

func NewGraph() *Graph {
	return &Graph{
		Edges:    []Edge{},
		Nodes:    make(map[string]*Node),
		Incoming: make(map[string][]*Node),
		Outgoing: make(map[string][]*Node),
	}
}

func (g *Graph) ApplyInputs(node *Node) {
	for _, e := range g.Edges {
		if e.Target != node.Id {
			continue
		}
		node.Io[e.TargetHandle].Value = g.Nodes[e.Source].Io[e.SourceHandle].Value
	}
}

func HasCycle(g *Graph) bool {
	inDegree := map[string]int{}

	for id := range g.Nodes {
		inDegree[id] = len(g.Incoming[id])
	}

	q := []string{}

	done := 0

	for id, degree := range inDegree {
		if degree == 0 {
			q = append(q, id)
		}
	}

	for len(q) > 0 {
		n := q[0]
		q = q[1:]
		done++

		for _, out := range g.Outgoing[n] {
			inDegree[out.Id]--
			if inDegree[out.Id] == 0 {
				q = append(q, out.Id)
			}
		}
	}

	return done != len(g.Nodes)
}
