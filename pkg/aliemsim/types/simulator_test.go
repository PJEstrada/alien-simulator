package types

import (
	"alien-invasion-simulator/pkg/graph"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SimulatorTestSuite struct {
	suite.Suite
}

func TestSimulatorTestSuite(t *testing.T) {
	suite.Run(t, &SimulatorTestSuite{})
}

func (s *SimulatorTestSuite) TestNewAlienSimulator() {
	m := Map[City, Direction]{}
	aliens := []*Alien{{Name: "alien"}}
	maxIters := 55
	verbose := true
	sim := NewAlienSimulator(&m, aliens, maxIters, verbose)
	s.Equal(sim.Map, &m)
	s.Equal(sim.Aliens, aliens)
	s.Equal(sim.MaxIterations, maxIters)
}

func (s *SimulatorTestSuite) TestGetStats() {
	m := Map[City, Direction]{
		Cities: CityStore{},
	}
	aliens := []*Alien{{Name: "alien"}}
	maxIters := 55
	verbose := true
	sim := NewAlienSimulator(&m, aliens, maxIters, verbose)
	stats := sim.getStats()

	s.Equal(stats,
		fmt.Sprintf("Iteration # %d - Total Aliens: %d - Dead Aliens: %d - Cities Left: %d - Num Aliens Reached Max Moves: %d",
			sim.CurrentIteration, len(sim.Aliens), sim.NumDeadAliens, len(sim.Map.Cities), sim.NumAliensReachedMaxMoves))
}

func (s *SimulatorTestSuite) TestSimulateInvasion() {
	vals := []struct {
		name                     string
		simBuilder               func() AlienSimulator
		withErr                  bool
		err                      error
		NumDeadAliens            int
		CurrentIteration         int
		NumAliensCannotMove      int
		NumAliensReachedMaxMoves int
	}{
		{
			name: "Invasion with no cities",
			simBuilder: func() AlienSimulator {
				m := Map[City, Direction]{
					Cities: CityStore{},
				}
				aliens := []*Alien{{Name: "alien"}}
				maxIters := 55
				verbose := true
				sim := NewAlienSimulator(&m, aliens, maxIters, verbose)
				return sim
			},
			withErr:                  false,
			err:                      nil,
			NumDeadAliens:            0,
			CurrentIteration:         0,
			NumAliensCannotMove:      0,
			NumAliensReachedMaxMoves: 0,
		},
		{
			name: "Invasion with 3 cities",
			simBuilder: func() AlienSimulator {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("a")
				m.getOrCreateCity("b")
				m.getOrCreateCity("c")
				m.AddPath("a", "b", "north")
				m.AddPath("b", "c", "west")
				aliens := []*Alien{{Name: "alien1", ID: 0, CurrentCityName: "a", Map: &m}, {Name: "alien2", ID: 1, CurrentCityName: "b", Map: &m}}
				maxIters := 1
				sim := NewAlienSimulator(&m, aliens, maxIters, false)
				return sim
			},
			NumDeadAliens:            2,
			CurrentIteration:         1,
			NumAliensCannotMove:      2,
			NumAliensReachedMaxMoves: 0,
			withErr:                  false,
			err:                      nil,
		},
		{
			name: "Invasion with some dead aliens",
			simBuilder: func() AlienSimulator {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("a")
				m.getOrCreateCity("b")
				m.getOrCreateCity("c")
				aliens := []*Alien{{Name: "alien1", ID: 0, CurrentCityName: "a", Map: &m}, {Name: "alien2", ID: 1, CurrentCityName: "b", Map: &m, IsDead: true}}
				maxIters := 10
				sim := NewAlienSimulator(&m, aliens, maxIters, false)
				sim.NumDeadAliens = 1
				sim.NumAliensReachedMaxMoves = 2
				sim.NumAliensCannotMove = 2
				return sim
			},
			NumDeadAliens:            1,
			CurrentIteration:         0,
			NumAliensCannotMove:      2,
			NumAliensReachedMaxMoves: 2,
			withErr:                  false,
			err:                      nil,
		},
		{
			name: "Invasion with aliens that reached max iterations",
			simBuilder: func() AlienSimulator {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("a")
				m.getOrCreateCity("b")
				m.getOrCreateCity("c")
				maxIters := 10
				aliens := []*Alien{{Name: "alien1", ID: 0, CurrentCityName: "a", Map: &m, NumMovements: maxIters, CanMove: true}}

				sim := NewAlienSimulator(&m, aliens, maxIters, true)
				return sim
			},
			NumDeadAliens:            0,
			CurrentIteration:         0,
			NumAliensCannotMove:      1,
			NumAliensReachedMaxMoves: 1,
			withErr:                  false,
			err:                      nil,
		},
		{
			name: "Invasion with alien move error",
			simBuilder: func() AlienSimulator {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("a")
				m.getOrCreateCity("b")
				m.getOrCreateCity("c")
				maxIters := 10
				aliens := []*Alien{{Name: "alien1", ID: 0, CurrentCityName: "bad city", Map: &m, NumMovements: 0, CanMove: true}}

				sim := NewAlienSimulator(&m, aliens, maxIters, true)
				return sim
			},
			NumDeadAliens:            0,
			CurrentIteration:         0,
			NumAliensCannotMove:      0,
			NumAliensReachedMaxMoves: 0,
			withErr:                  true,
			err:                      ErrorCityDoesNotExists,
		},
		{
			name: "Invasion with no more aliens to move",
			simBuilder: func() AlienSimulator {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("a")
				m.getOrCreateCity("b")
				m.getOrCreateCity("c")
				maxIters := 10
				aliens := []*Alien{{Name: "alien1", ID: 0, CurrentCityName: "a", Map: &m, NumMovements: 0, CanMove: false}}

				sim := NewAlienSimulator(&m, aliens, maxIters, true)
				sim.NumAliensCannotMove = 1
				return sim
			},
			NumDeadAliens:            0,
			CurrentIteration:         0,
			NumAliensCannotMove:      1,
			NumAliensReachedMaxMoves: 0,
			withErr:                  false,
			err:                      nil,
		},
	}

	for _, val := range vals {
		s.Run(val.name, func() {
			sim := val.simBuilder()
			err := sim.SimulateInvasion()
			if val.withErr {
				s.NotNil(err)
				s.EqualError(err, val.err.Error())
			} else {
				s.Nil(err)
				s.Equal(sim.NumDeadAliens, val.NumDeadAliens)
				s.True(sim.CurrentIteration >= val.CurrentIteration)
				s.Equal(sim.NumAliensCannotMove, val.NumAliensCannotMove)
				s.Equal(sim.NumAliensReachedMaxMoves, val.NumAliensReachedMaxMoves)
			}

		})
	}
}

func (s *SimulatorTestSuite) TestAlienMove() {
	vals := []struct {
		name                     string
		simBuilder               func() AlienSimulator
		withErr                  bool
		deadAliens               bool
		err                      error
		NumDeadAliens            int
		movedTo                  string
		NumAliensCannotMove      int
		isTrapped                bool
		NumAliensReachedMaxMoves int
	}{
		{
			name: "Alien move from a city that does not exists.",
			simBuilder: func() AlienSimulator {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				aliens := []*Alien{{Name: "alien", CurrentCityName: "invalidciy", Map: &m}, {Name: "alien2", CurrentCityName: "b", Map: &m}}
				maxIters := 55
				sim := NewAlienSimulator(&m, aliens, maxIters, true)
				return sim
			},
			withErr:                  true,
			deadAliens:               false,
			err:                      ErrorCityDoesNotExists,
			NumDeadAliens:            0,
			movedTo:                  "",
			NumAliensCannotMove:      0,
			NumAliensReachedMaxMoves: 0,
		},
		{
			name: "Alien move from a city that exist.",
			simBuilder: func() AlienSimulator {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("a")
				m.getOrCreateCity("b")
				m.AddPath("a", "b", "east")
				aliens := []*Alien{{Name: "alien", CurrentCityName: "a", Map: &m, ID: 0}, {Name: "alien2", CurrentCityName: "b", Map: &m, ID: 1}}
				maxIters := 55
				sim := NewAlienSimulator(&m, aliens, maxIters, true)
				return sim
			},
			withErr:                  false,
			deadAliens:               false,
			err:                      nil,
			NumDeadAliens:            0,
			movedTo:                  "b",
			NumAliensCannotMove:      0,
			NumAliensReachedMaxMoves: 0,
		},
		{
			name: "Alien move from a city that has an alien.",
			simBuilder: func() AlienSimulator {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				aliens := []*Alien{{Name: "alien", CurrentCityName: "a", Map: &m, ID: 0}, {Name: "alien2", CurrentCityName: "b", Map: &m, ID: 1}}
				m.getOrCreateCity("a")
				m.getOrCreateCity("b")
				m.Cities["b"].Aliens = []*Alien{aliens[1]}
				v := m.Graph.GetVertexByStringID("b")
				v.Data.Aliens = []*Alien{aliens[1]}
				m.AddPath("a", "b", "east")

				maxIters := 55
				sim := NewAlienSimulator(&m, aliens, maxIters, true)
				return sim
			},
			withErr:                  false,
			deadAliens:               true,
			err:                      nil,
			NumDeadAliens:            2,
			movedTo:                  "b",
			NumAliensCannotMove:      0,
			NumAliensReachedMaxMoves: 0,
		},
		{
			name: "Try to move trapped alien.",
			simBuilder: func() AlienSimulator {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				aliens := []*Alien{{Name: "alien", CurrentCityName: "a", Map: &m, ID: 0}}
				m.getOrCreateCity("a")
				maxIters := 55
				sim := NewAlienSimulator(&m, aliens, maxIters, true)
				return sim
			},
			withErr:                  false,
			deadAliens:               false,
			isTrapped:                true,
			err:                      nil,
			NumDeadAliens:            0,
			movedTo:                  "a",
			NumAliensCannotMove:      0,
			NumAliensReachedMaxMoves: 0,
		},
	}
	for _, val := range vals {
		s.Run(val.name, func() {
			sim := val.simBuilder()
			_, err := sim.alienMove(sim.Aliens[0])
			if val.withErr {
				s.NotNil(err)
				s.EqualError(err, val.err.Error())
			} else {

				for _, al := range sim.Aliens {
					trapped, _ := al.isTrapped()
					s.Equal(al.CurrentCityName, val.movedTo)
					s.Equal(sim.NumDeadAliens, val.NumDeadAliens)
					s.Equal(al.IsDead, val.deadAliens)
					if val.isTrapped {
						s.Equal(trapped, val.isTrapped)
					}

				}
			}
		})
	}
}
