package graph

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type GraphTestSuite struct {
	suite.Suite
}
type IdentifiableMock struct {
	id string
}

func (i IdentifiableMock) ID() string {
	return i.id
}
func TestGraphTestSuite(t *testing.T) {
	suite.Run(t, &GraphTestSuite{})
}
func (s *GraphTestSuite) TestIdVertex() {
	theid := "the id"
	vertex := Vertex[IdentifiableMock, string]{id: VertexID(theid)}
	res := vertex.Id()
	s.Equal(res, VertexID(theid))
}

func (s *GraphTestSuite) TestIdEdge() {
	theid := EdgeId{From: VertexID("from"), To: VertexID("to")}
	vertex := Edge[IdentifiableMock, string]{id: theid}
	res := vertex.Id()
	s.Equal(res, theid)
}

func (s *GraphTestSuite) TestNewGraph() {

	g := NewGraph[IdentifiableMock, string]()

	s.NotNil(g)
	s.NotNil(g.GetNodes())
	s.NotNil(g.GetEdges())
}

func (s *GraphTestSuite) TestNewGraphFrom() {
	vals := []struct {
		name    string
		vertex  []IdentifiableMock
		edges   map[EdgeId]string
		withErr bool
	}{
		{
			name: "New graph from existing vertexes and edges",
			vertex: []IdentifiableMock{
				{
					id: "v1",
				},
				{
					id: "v2",
				},
			},
			edges: map[EdgeId]string{
				EdgeId{From: VertexID("v1"), To: VertexID("v2")}: "e1",
			},
			withErr: false,
		},
		{
			name: "New graph with edge from vertex that does not exists.",
			vertex: []IdentifiableMock{
				{
					id: "v1",
				},
				{
					id: "v2",
				},
			},
			edges: map[EdgeId]string{
				EdgeId{From: VertexID("v7"), To: VertexID("v2")}: "e1",
			},
			withErr: true,
		},
		{
			name: "New graph with edge to vertex that does not exists.",
			vertex: []IdentifiableMock{
				{
					id: "v1",
				},
				{
					id: "v2",
				},
			},
			edges: map[EdgeId]string{
				EdgeId{From: VertexID("v1"), To: VertexID("v5")}: "e1",
			},
			withErr: true,
		},
	}
	for _, val := range vals {
		s.Run(val.name, func() {
			g, err := NewGraphFrom[IdentifiableMock, string](val.vertex, val.edges)
			if val.withErr {
				s.NotNil(err)
			} else {
				s.Equal(g.GetEdges()[EdgeId{From: VertexID("v1"), To: VertexID("v2")}].Data, val.edges[EdgeId{From: VertexID("v1"), To: VertexID("v2")}])
				s.Equal(g.GetNodes()[VertexID("v1")].Data.id, val.vertex[0].id)
				s.Equal(g.GetNodes()[VertexID("v2")].Data.id, val.vertex[1].id)
				s.Nil(err)
			}

		})
	}

}

func (s *GraphTestSuite) TestAddVertex() {
	g := NewGraph[IdentifiableMock, string]()
	idMock := IdentifiableMock{id: "v1"}
	vId := g.AddVertex(idMock)
	s.Equal(vId, VertexID(idMock.id))
}

func (s *GraphTestSuite) TestGetVertexByStringID() {
	g := NewGraph[IdentifiableMock, string]()
	g.AddVertex(IdentifiableMock{
		id: "v1",
	})
	vals := []struct {
		name  string
		id    string
		found bool
	}{
		{

			name:  "get existing vertex",
			id:    "v1",
			found: true,
		},
		{

			name:  "get unexisting vertex",
			id:    "v2",
			found: false,
		},
	}
	for _, val := range vals {
		s.Run(val.name, func() {
			res := g.GetVertexByStringID(val.id)
			if val.found {
				s.Equal(res.Data.id, val.id)
			} else {
				s.Nil(res)
			}

		})
	}
}

func (s *GraphTestSuite) TestGetVertexByID() {
	g := NewGraph[IdentifiableMock, string]()
	g.AddVertex(IdentifiableMock{
		id: "v1",
	})
	vals := []struct {
		name  string
		id    string
		found bool
	}{
		{

			name:  "get existing vertex",
			id:    "v1",
			found: true,
		},
		{

			name:  "get unexisting vertex",
			id:    "v2",
			found: false,
		},
	}
	for _, val := range vals {
		s.Run(val.name, func() {
			res := g.GetVertexByID(VertexID(val.id))
			if val.found {
				s.Equal(res.Data.id, val.id)
			} else {
				s.Nil(res)
			}

		})
	}
}

func (s *GraphTestSuite) TestRemoveVertex() {
	g := NewGraph[IdentifiableMock, string]()
	g2 := NewGraph[IdentifiableMock, string]()
	g.AddVertex(IdentifiableMock{
		id: "v1",
	})
	g.AddVertex(IdentifiableMock{
		id: "v5",
	})
	g2.AddVertex(IdentifiableMock{
		id: "v1",
	})
	g2.AddVertex(IdentifiableMock{
		id: "v5",
	})

	vals := []struct {
		name  string
		id    string
		found bool
		g     Graph[IdentifiableMock, string]
		edges []Edge[IdentifiableMock, string]
	}{
		{

			name:  "remove existing vertex",
			id:    "v1",
			g:     g,
			found: true,
		},
		{

			name:  "remove unexisting vertex",
			id:    "v2",
			g:     g,
			found: false,
		},
		{

			name:  "remove vertex with edges",
			id:    "v1",
			found: true,
			g:     g2,
			edges: []Edge[IdentifiableMock, string]{
				{
					id:   EdgeId{From: VertexID("v1"), To: VertexID("v5")},
					Data: "myedge",
				},
			},
		},
	}
	for _, val := range vals {
		s.Run(val.name, func() {
			for _, edge := range val.edges {
				val.g.AddEdge(edge.id.From, edge.id.To, edge.Data)
			}
			err := val.g.RemoveVertex(Vertex[IdentifiableMock, string]{
				id:   VertexID(val.id),
				Data: IdentifiableMock{id: val.id},
			})
			if val.found {
				s.Nil(err)
				s.Equal(len(val.g.GetEdges()), 0)
			} else {
				s.NotNil(err)

			}

		})
	}
}

func (s *GraphTestSuite) TestRemoveVertexByID() {
	g := NewGraph[IdentifiableMock, string]()
	g.AddVertex(IdentifiableMock{
		id: "v1",
	})
	vals := []struct {
		name  string
		id    string
		found bool
	}{
		{

			name:  "remove existing vertex",
			id:    "v1",
			found: true,
		},
		{

			name:  "remove unexisting vertex",
			id:    "v2",
			found: false,
		},
	}
	for _, val := range vals {
		s.Run(val.name, func() {
			err := g.RemoveVertexByID(VertexID(val.id))
			if val.found {
				s.Nil(err)
			} else {
				s.NotNil(err)
			}

		})
	}
}

func (s *GraphTestSuite) TestAddEdge() {
	g := NewGraph[IdentifiableMock, string]()
	g.AddVertex(IdentifiableMock{
		id: "v1",
	})
	g.AddVertex(IdentifiableMock{
		id: "v2",
	})
	vals := []struct {
		name     string
		from     string
		to       string
		data     string
		withErr  bool
		addTwice bool
	}{
		{
			name: "add edge with existing vertexes",
			from: "v1",
			to:   "v2",
			data: "e1",
		},
		{
			name:    "add edge with non existent vertexes",
			from:    "v1",
			to:      "v7",
			data:    "e1",
			withErr: true,
		},
		{
			name:     "add edge twice",
			from:     "v1",
			to:       "v2",
			data:     "e1",
			withErr:  true,
			addTwice: true,
		},
	}
	for _, val := range vals {
		s.Run(val.name, func() {

			eId, err := g.AddEdge(VertexID(val.from), VertexID(val.to), val.data)
			if val.addTwice {
				eId, err = g.AddEdge(VertexID(val.from), VertexID(val.to), val.data)
			}
			if val.withErr {
				s.NotNil(err)
			} else {
				s.Equal(eId.To, VertexID(val.to))
				s.Equal(eId.From, VertexID(val.from))
			}
		})
	}
}

func (s *GraphTestSuite) TestGetEdge() {
	g := NewGraph[IdentifiableMock, string]()
	g.AddVertex(IdentifiableMock{
		id: "v1",
	})
	g.AddVertex(IdentifiableMock{
		id: "v2",
	})
	g.AddEdge(VertexID("v1"), VertexID("v2"), "e1")
	vals := []struct {
		name      string
		from      string
		to        string
		edgeFound bool
	}{
		{
			name:      "Get edge with existing vertexes",
			from:      "v1",
			to:        "v2",
			edgeFound: true,
		},
		{
			name:      "Get edge with vertex that does not exists.",
			from:      "v1",
			to:        "v8",
			edgeFound: false,
		},
	}

	for _, val := range vals {
		s.Run(val.name, func() {
			res := g.GetEdge(VertexID(val.from), VertexID(val.to))
			if val.edgeFound {
				s.NotNil(res)
				s.Equal(res.To.id, VertexID(val.to))
				s.Equal(res.From.id, VertexID(val.from))
			} else {
				s.Nil(res)
			}
		})
	}
}

func (s *GraphTestSuite) TestRemoveEdge() {
	g := NewGraph[IdentifiableMock, string]()
	g.AddVertex(IdentifiableMock{
		id: "v1",
	})
	g.AddVertex(IdentifiableMock{
		id: "v2",
	})
	g.AddEdge(VertexID("v1"), VertexID("v2"), "e1")
	vals := []struct {
		name   string
		from   string
		to     string
		witErr bool
	}{
		{
			name:   "Remove edge with existing vertexes",
			from:   "v1",
			to:     "v2",
			witErr: true,
		},
		{
			name:   "Remove with vertex that does not exists [To].",
			from:   "v1",
			to:     "v8",
			witErr: false,
		},
		{
			name:   "Remove with vertex that does not exists [From].",
			from:   "v8",
			to:     "v2",
			witErr: false,
		},
	}

	for _, val := range vals {
		s.Run(val.name, func() {
			err := g.RemoveEdge(VertexID(val.from), VertexID(val.to))
			if val.witErr {
				s.Nil(err)
			} else {
				s.NotNil(err)
			}
		})
	}
}
