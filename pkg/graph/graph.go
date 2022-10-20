package graph

import (
	"fmt"
	"sync"
)

type VertexID string

type Hash[N any] func(N) VertexID

type Identifiable interface {
	ID() string
}

type Vertex[N Identifiable, E any] struct {
	id            VertexID
	Data          N
	IncomingEdges map[VertexID]*Edge[N, E]
	OutgoingEdges map[VertexID]*Edge[N, E]
}

func (n Vertex[N, E]) Id() VertexID {
	return n.id
}

type EdgeId struct {
	From VertexID
	To   VertexID
}

type Edge[N Identifiable, E any] struct {
	id   EdgeId
	Data E
	From *Vertex[N, E]
	To   *Vertex[N, E]
}

func (e Edge[N, E]) Id() EdgeId {
	return e.id
}

type Graph[N Identifiable, E any] struct {
	nodes map[VertexID]*Vertex[N, E]
	edges map[EdgeId]*Edge[N, E]
	hash  Hash[N]
	rw    sync.RWMutex
}

func NewGraph[N Identifiable, E any]() Graph[N, E] {
	return Graph[N, E]{
		nodes: map[VertexID]*Vertex[N, E]{},
		edges: map[EdgeId]*Edge[N, E]{},
		rw:    sync.RWMutex{},
	}
}

func NewGraphFrom[N Identifiable, E any](vertex []N, edges map[EdgeId]E) (Graph[N, E], error) {
	d := NewGraph[N, E]()
	for _, n := range vertex {
		d.addVertex(n)
	}
	for id, v := range edges {
		_, err := d.addEdge(id.From, id.To, v)
		if err != nil {
			return d, err
		}
	}
	return d, nil
}

func (d *Graph[N, E]) AddVertex(n N) VertexID {
	d.rw.Lock()
	defer d.rw.Unlock()
	return d.addVertex(n)
}

func (d *Graph[N, E]) addVertex(n N) VertexID {
	id := VertexID(n.ID())
	d.nodes[id] = &Vertex[N, E]{
		id:            id,
		Data:          n,
		IncomingEdges: map[VertexID]*Edge[N, E]{},
		OutgoingEdges: map[VertexID]*Edge[N, E]{},
	}
	return id
}
func (d *Graph[N, E]) GetVertexByStringID(id string) *Vertex[N, E] {
	d.rw.RLock()
	defer d.rw.RUnlock()
	idVertex := VertexID(id)
	return d.getVertexByID(idVertex)
}

func (d *Graph[N, E]) GetVertexByID(id VertexID) *Vertex[N, E] {
	d.rw.RLock()
	defer d.rw.RUnlock()
	return d.getVertexByID(id)
}

func (d *Graph[N, E]) getVertexByID(id VertexID) *Vertex[N, E] {
	if v, ok := d.nodes[id]; ok {
		return v
	}
	return nil
}

func (d *Graph[N, E]) RemoveVertexByID(id VertexID) error {
	d.rw.Lock()
	defer d.rw.Unlock()
	return d.removeVertexByID(id)
}

func (d *Graph[N, E]) RemoveVertex(v Vertex[N, E]) error {
	d.rw.Lock()
	defer d.rw.Unlock()
	id := VertexID(v.Data.ID())
	return d.removeVertexByID(id)
}

func (d *Graph[N, E]) removeVertexByID(id VertexID) error {
	node := d.getVertexByID(id)
	if node == nil {
		return fmt.Errorf("Vertex %v not found to remove", id)
	}
	for edge, _ := range d.edges {
		if edge.From == id || edge.To == id {
			err := d.removeEdge(edge.From, edge.To)
			if err != nil {
				return err
			}
		}
	}
	delete(d.nodes, id)
	return nil
}

func (d *Graph[N, E]) AddEdge(from, to VertexID, value E) (EdgeId, error) {
	d.rw.Lock()
	defer d.rw.Unlock()
	return d.addEdge(from, to, value)
}

func (d *Graph[N, E]) addEdge(from, to VertexID, value E) (EdgeId, error) {
	fromNode, toNode := d.getVertexByID(from), d.getVertexByID(to)

	if fromNode == nil {
		return EdgeId{}, fmt.Errorf("Vertex %v not found", from)
	}
	if toNode == nil {
		return EdgeId{}, fmt.Errorf("Vertex %v not found", to)
	}

	id := EdgeId{from, to}
	if _, ok := d.edges[id]; ok {
		return EdgeId{}, fmt.Errorf("Edge %v -> %v already exists", from, to)
	}
	edge := Edge[N, E]{id, value, fromNode, toNode}

	fromNode.OutgoingEdges[to] = &edge
	toNode.IncomingEdges[from] = &edge
	d.edges[id] = &edge

	return id, nil
}

func (d *Graph[N, E]) GetEdge(from, to VertexID) *Edge[N, E] {
	d.rw.RLock()
	defer d.rw.RUnlock()
	return d.getEdge(from, to)
}

func (d *Graph[N, E]) getEdge(from, to VertexID) *Edge[N, E] {
	if v, ok := d.edges[EdgeId{from, to}]; ok {
		return v
	}
	return nil
}

func (d *Graph[N, E]) RemoveEdge(from, to VertexID) error {
	d.rw.Lock()
	defer d.rw.Unlock()
	return d.removeEdge(from, to)
}

func (d *Graph[N, E]) removeEdge(from, to VertexID) error {
	fromNode, toNode := d.getVertexByID(from), d.getVertexByID(to)
	if fromNode == nil {
		return fmt.Errorf("Vertex %v not found", from)
	}
	if toNode == nil {
		return fmt.Errorf("Vertex %v not found", to)
	}
	delete(d.edges, EdgeId{from, to})
	delete(fromNode.OutgoingEdges, to)
	delete(toNode.IncomingEdges, from)

	return nil
}

func (d *Graph[N, E]) GetEdges() map[EdgeId]*Edge[N, E] {
	return d.edges
}

func (d *Graph[N, E]) GetNodes() map[VertexID]*Vertex[N, E] {
	return d.nodes
}
