package types

import (
	"alien-invasion-simulator/pkg/graph"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AlienTestSuite struct {
	suite.Suite
}

func TestAlienTestSuite(t *testing.T) {
	suite.Run(t, &AlienTestSuite{})
}

func (s *AlienTestSuite) TestNewAlien() {
	id := 123
	name := "test alien"
	city := "an alien city"
	alient := NewAlien(id, name, city, &Map[City, Direction]{})
	s.Equal(id, alient.ID)
	s.Equal(name, alient.Name)
	s.Equal(city, alient.CurrentCityName)
}

func (s *AlienTestSuite) TestIsTrapped() {
	id := 123
	name := "test alien"
	city1 := &City{Name: "city1"}
	city2 := &City{Name: "city2"}
	mapobj := Map[City, Direction]{
		Cities: CityStore{},
		Graph:  graph.NewGraph[City, Direction](),
	}
	_ = mapobj.getOrCreateCity("city1")
	_ = mapobj.getOrCreateCity("city2")
	_ = mapobj.AddPath("city1", "city2", "north")
	vals := []struct {
		name      string
		alien     Alien
		m         Map[City, Direction]
		hasErr    bool
		isTrapped bool
		err       error
	}{
		{
			name:      "alien with 1 path on current city",
			alien:     NewAlien(id, name, city1.Name, &mapobj),
			m:         mapobj,
			hasErr:    false,
			isTrapped: false,
		},
		{
			name:      "alien with 0 paths on current city",
			alien:     NewAlien(id, name, city2.Name, &mapobj),
			m:         mapobj,
			hasErr:    false,
			isTrapped: true,
		},
		{
			name:      "alien with in an invalid city name",
			alien:     NewAlien(id, name, "unknowncity", &mapobj),
			m:         mapobj,
			hasErr:    true,
			isTrapped: true,
			err:       ErrorCityDoesNotExists,
		},
	}
	for _, val := range vals {
		s.Run(val.name, func() {
			trapped, err := val.alien.isTrapped()
			if !val.hasErr {
				s.Nil(err)
			} else {
				s.NotNil(err)
				s.EqualError(err, val.err.Error())
			}
			s.Equal(trapped, val.isTrapped)

		})
	}
}

func (s *AlienTestSuite) TestInvadeCity() {
	citytest := &City{Name: "citytest"}
	id := 123
	name := "citytest"
	alien := NewAlien(id, "An alien", name, &Map[City, Direction]{})
	cityAfterInvasion := alien.invadeCity(citytest)
	s.Equal(alien.NumMovements, 1)
	s.Equal(cityAfterInvasion.Aliens[0], &alien)
	s.Equal(alien.CurrentCityName, cityAfterInvasion.Name)

}
