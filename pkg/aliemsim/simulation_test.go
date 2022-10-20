package aliemsim

import (
	"alien-invasion-simulator/pkg/aliemsim/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SimulationTestSuite struct {
	suite.Suite
}

func TestSimulationSuite(t *testing.T) {
	suite.Run(t, &SimulationTestSuite{})
}

func (s *SimulationTestSuite) TestSpawnAliens() {
	numAliens := 5
	mapobj := types.Map[types.City, types.Direction]{Cities: map[string]*types.City{
		"city1": &types.City{Name: "city1"},
		"city2": &types.City{Name: "city2"},
		"city3": &types.City{Name: "city3"},
	}}
	aliens := spawnAliens(numAliens, &mapobj)
	s.Equal(len(aliens), numAliens)
	for _, al := range aliens {
		s.NotEqual(al.Name, "")
	}
}

// TestStartSimulationSuccess tests a successful simulation start
func (s *SimulationTestSuite) TestStartSimulationSuccess() {
	mockFs := fakeFS{}
	newMapFromReader = types.NewMapFromReaderMock
	err := StartSimulation("/home/cities.txt", 10, mockFs, 10, true)
	s.Nil(err)
	newMapFromReader = types.NewMapFromReader

}

// TestStartSimulationWrongFile tests an incorrect path being sent
func (s *SimulationTestSuite) TestStartSimulationWrongFile() {
	mockFs := fakeFSErr{}
	err := StartSimulation("some wrong path", 100, mockFs, 10, true)
	s.NotNil(err)
	s.EqualError(err, FileOpenErrMock.Error())

}

// TestStartSimulationErrorMap tests an error building the map during simulation
func (s *SimulationTestSuite) TestStartSimulationErrorMap() {
	mockFs := fakeFS{}
	newMapFromReader = types.NewMapFromReaderMockErr
	err := StartSimulation("/home/cities.txt", 100, mockFs, 10, true)
	s.NotNil(err)
	s.EqualError(err, types.ErrorNewMapMock.Error())
	newMapFromReader = types.NewMapFromReader

}
