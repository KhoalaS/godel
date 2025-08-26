package pipeline

type Graph struct {
	Nodes    map[string]*Node   `json:"nodes"`
	Incoming map[string][]*Node `json:"incoming"`
	Outgoing map[string][]*Node `json:"outgoing"`
}

func NewGraph() *Graph {
	return &Graph{
		Nodes:    make(map[string]*Node),
		Incoming: make(map[string][]*Node),
		Outgoing: make(map[string][]*Node),
	}
}

func findFirstNodes(graph *Graph) []*Node {
	nodes := []*Node{}
	for k, v := range graph.Incoming {
		if len(v) == 0 {
			nodes = append(nodes, graph.Nodes[k])
		}
	}
	return nodes
}
