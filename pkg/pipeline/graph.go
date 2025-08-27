package pipeline

type Graph struct {
	Edges    []Edge             `json:"edges"`
	Nodes    map[string]*Node   `json:"nodes"`
	Incoming map[string][]*Node `json:"incoming"`
	Outgoing map[string][]*Node `json:"outgoing"`
}

type Edge struct {
	Target       string `json:"target"`
	Source       string `json:"source"`
	SourceHandle string `json:"sourceHandle"`
	TargetHandle string `json:"targetHandle"`
}

func NewGraph() *Graph {
	return &Graph{
		Edges:    []Edge{},
		Nodes:    make(map[string]*Node),
		Incoming: make(map[string][]*Node),
		Outgoing: make(map[string][]*Node),
	}
}

func ApplyInputs(graph *Graph, node *Node) {
	for _, e := range graph.Edges {
		if e.Target != node.Id {
			continue
		}
		node.Inputs[e.TargetHandle] = graph.Nodes[e.Source].Outputs[e.SourceHandle]
	}
}
