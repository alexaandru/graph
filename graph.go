package utils

import (
    "fmt"
    "strconv"
    "strings"
)

// FIXME: Add tests!

// Node represents a generic Graph node
type Node interface {
    String() string
}

// Nodes consists of a list of Node
type Nodes []Node

// NodeMap consists of a map of Node
type NodeMap map[Node]bool

// Graph represents a directed, weighted graph implementation
type Graph struct {
    Nodes        NodeMap
    Source, Sink Nodes
    Edges        map[Node]Nodes
    Weights      map[Node]intMap
    InDegree     intMap
    OutDegree    intMap
}

// IntNode represents an int Graph node
type IntNode int

type intMap map[Node]int

// GraphNew initializes a new Graph structure and returns a *Graph
func GraphNew() (g *Graph) {
    g = new(Graph)
    g.Nodes = NodeMap{}
    g.Edges = map[Node]Nodes{}
    g.Weights = map[Node]intMap{}
    g.InDegree, g.OutDegree = intMap{}, intMap{}
    g.Source, g.Sink = Nodes{}, Nodes{}

    return
}

// AddEdge adds a directed edge between src and dst, optionally with a weigth.
// The weight defaults to 1 if not given.
func (g *Graph) AddEdge(src, dst Node, opts ...int) {
    weight := 1
    if len(opts) > 0 {
        weight = opts[0]
    }

    g.Nodes[src] = true
    g.Nodes[dst] = true

    g.InDegree[dst]++
    g.OutDegree[src]++

    g.Edges[src] = append(g.Edges[src], dst)
    if g.Weights[src] == nil {
        g.Weights[src] = intMap{dst: weight}
    } else {
        g.Weights[src][dst] = weight
    }
}

// RemoveEdge deletes the edge between src and dst
func (g *Graph) RemoveEdge(src, dst Node) {
    rest := g.Edges[src].Remove(dst)
    if len(rest) == 0 {
        delete(g.Edges, src)
    } else {
        g.Edges[src] = rest
    }

    g.InDegree[dst]--
    g.OutDegree[src]--
}

// HasEdge determines if there is an edge between src and dst
func (g *Graph) HasEdge(src, dst Node) bool {
    for _, node := range g.Edges[src] {
        if node == dst {
            return true
        }
    }

    return false
}

// Clone clones an entire Graph
func (g *Graph) Clone() (h *Graph) {
    h = GraphNew()
    h.Nodes = g.NodesClone()
    h.Edges = g.EdgesClone()
    h.Weights = g.WeightsClone()
    h.InDegree = g.InDegreeClone()
    h.OutDegree = g.OutDegreeClone()
    h.Source = g.SourceClone()
    h.Sink = g.SinkClone()

    return
}

// NodesClone clones the nodes of a Graph
func (g *Graph) NodesClone() (nc NodeMap) {
    nc = NodeMap{}
    for k, v := range g.Nodes {
        nc[k] = v
    }

    return
}

// EdgesClone clones the edges of a Graph
func (g *Graph) EdgesClone() (ec map[Node]Nodes) {
    ec = map[Node]Nodes{}
    for k, v := range g.Edges {
        ec[k] = v.Clone()
    }

    return
}

// WeightsClone clones the weights of a Graph
func (g *Graph) WeightsClone() (wc map[Node]intMap) {
    wc = map[Node]intMap{}
    for k, v := range g.Weights {
        wc[k] = v.Clone()
    }

    return
}

// InDegreeClone clones the indegrees of a Graph
func (g *Graph) InDegreeClone() (ic intMap) {
    ic = intMap{}
    for k, v := range g.InDegree {
        ic[k] = v
    }

    return
}

// OutDegreeClone clones the outdegrees of a Graph
func (g *Graph) OutDegreeClone() (oc intMap) {
    oc = intMap{}
    for k, v := range g.OutDegree {
        oc[k] = v
    }

    return
}

// SourceClone clones the source(s) of a Graph
func (g *Graph) SourceClone() (sc Nodes) {
    sc = make(Nodes, len(g.Source))
    copy(sc, g.Source)

    return
}

// SinkClone clones the sink(s) of a Graph
func (g *Graph) SinkClone() (sc Nodes) {
    sc = make(Nodes, len(g.Sink))
    copy(sc, g.Sink)

    return
}

// ComputeSourceAndSink computes the source(s) and sink(s) of a Graph
func (g *Graph) ComputeSourceAndSink() {
    g.Source = Nodes{}
    g.Sink = Nodes{}

    for n := range g.Nodes {
        if g.InDegree[n] == 0 {
            g.Source = append(g.Source, n)
        }
        if g.OutDegree[n] == 0 {
            g.Sink = append(g.Sink, n)
        }
    }
}

// TopoSort performs a topological sort of a Graph
func (g *Graph) TopoSort() (ts Nodes) {
    ts = Nodes{}
    h := g.Clone()
    g.ComputeSourceAndSink()
    s := g.SourceClone()

    for len(s) > 0 {
        n := s[0]
        s = s[1:]
        ts = append(ts, n)
        for _, m := range g.Edges[n] {
            h.RemoveEdge(n, m)
            if h.InDegree[m] == 0 {
                s = append(s, m)
            }
        }
    }

    // FIXME: Find a way to NOT panic here. Makes testing harder...
    if len(h.Edges) > 0 {
        fmt.Println(h.Edges)
        panic("Not a DAG")
    }

    return
}

// LongestPath computes the longest path between src and dst in a Graph
func (g *Graph) LongestPath(src, dst Node) (l int, path Nodes) {
    path = Nodes{}
    ts := g.TopoSort()

    // Topo sort between src and dst inclusive.
    for k, v := range ts {
        if v == src {
            ts = ts[k:]
            break
        }
    }
    for k, v := range ts {
        if v == dst {
            ts = ts[:k+1]
            break
        }
    }

    dist := intMap{}
    back := map[Node]Node{}
    for i := 1; i < len(ts); i++ {
        max := 0
        for j := 0; j < i; j++ {
            v, w := ts[j], ts[i]
            if !g.HasEdge(v, w) {
                continue
            }
            curr := dist[v] + g.Weights[v][w]
            if curr > max {
                max = curr
                back[w] = v
            }
        }
        dist[ts[i]] = max
    }
    l = dist[dst]

    curr := dst
    path = append(path, curr)
    for len(back) > 0 {
        prev := back[curr]
        path = append(Nodes{prev}, path...)
        if prev == nil || prev == src {
            break
        }
        curr = prev
    }

    return
}

func (nm NodeMap) anyKey() (n Node) {
    for n = range nm {
        return
    }

    return
}

// Clone clones a Nodes list
func (n Nodes) Clone() (c Nodes) {
    c = make(Nodes, len(n))
    copy(c, n)

    return
}

// Remove removes a node from a Nodes list
func (n Nodes) Remove(node Node) (out Nodes) {
    out = make(Nodes, len(n))
    copy(out, n)
    for k, v := range out {
        if v == node {
            out[k], out[0] = out[0], out[k]
            break
        }
    }
    out = out[1:]

    return
}

func (n Nodes) String(opts ...string) (str string) {
    sep := "->"
    if len(opts) > 0 {
        sep = opts[0]
    }

    sx := make([]string, len(n))
    for k, node := range n {
        sx[k] = node.String()
    }

    return strings.Join(sx, sep)
}

func (i IntNode) String() string {
    return strconv.Itoa(int(i))
}

func (wm intMap) Clone() (wo intMap) {
    wo = intMap{}
    for k, v := range wm {
        wo[k] = v
    }

    return
}
