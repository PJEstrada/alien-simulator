package types

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type CityTestSuite struct {
	suite.Suite
}

func TestCityTestSuite(t *testing.T) {
	suite.Run(t, &CityTestSuite{})
}

func (s *CityTestSuite) TestDestroy() {
	vals := []struct {
		name          string
		city          City
		destroyed     bool
		numDeadAliens int
	}{
		{
			name:          "city with no aliens",
			city:          City{Name: "mycity", isDestroyed: false},
			destroyed:     true,
			numDeadAliens: 0,
		},
		{
			name: "city with 1 alien",
			city: City{Name: "mycity", isDestroyed: false, Aliens: []*Alien{
				{Name: "alien 1", IsDead: false},
			}},
			destroyed:     true,
			numDeadAliens: 1,
		},
		{
			name: "city with 1 alien",
			city: City{Name: "mycity", isDestroyed: false, Aliens: []*Alien{
				{Name: "alien 1", IsDead: false},
				{Name: "alien 2", IsDead: false},
			}},
			destroyed:     true,
			numDeadAliens: 2,
		},
	}

	for _, val := range vals {
		s.Run(val.name, func() {
			val.city.Destroy()
			s.Equal(val.city.isDestroyed, val.destroyed)

			s.Equal(len(val.city.Aliens), 0)
		})
	}
}
