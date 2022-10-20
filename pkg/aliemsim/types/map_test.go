package types

import (
	"alien-invasion-simulator/pkg/graph"
	"fmt"
	"github.com/stretchr/testify/suite"
	"io"
	"strings"
	"testing"
)

type MapTestSuite struct {
	suite.Suite
}

func TestMapTestSuite(t *testing.T) {
	suite.Run(t, &MapTestSuite{})
}

func (s *MapTestSuite) TestValidTextRow() {
	testVals := []struct {
		name    string
		fields  []string
		wantErr bool
	}{
		{
			name:    "city with 3 directions",
			fields:  []string{"hello", "north=Bar", "west=Baz", "south=Qu-ux"},
			wantErr: false,
		},
		{
			name:    "empty text row",
			fields:  []string{},
			wantErr: true,
		},
		{
			name:    "row with an invalid direction name",
			fields:  []string{"hello", "badir=Bar", "west=Baz", "south=Qu-ux"},
			wantErr: true,
		},
		{
			name:    "row with an invalid direction syntax",
			fields:  []string{"hello", "northBar", "west=Baz", "south=Qu-ux"},
			wantErr: true,
		},
	}
	for _, val := range testVals {
		s.Run(val.name, func() {
			err := validTextRow(val.fields)
			if val.wantErr {
				s.NotNil(err)
			} else {
				s.Nil(err)
			}
		})
	}
}

func (s *MapTestSuite) TestGetCitiesNames() {
	testVals := []struct {
		name     string
		fields   CityStore
		expected []string
	}{
		{
			name: "all cities exists.",
			fields: CityStore{
				"test1": &City{Name: "test1"},
				"test2": &City{Name: "test2"},
				"test3": &City{Name: "test3"},
			},

			expected: []string{"test1", "test2", "test3"},
		},
		{
			name: "2 cities exists.",
			fields: CityStore{
				"foo": &City{Name: "foo"},
				"bar": &City{Name: "bar"},
			},
			expected: []string{"foo", "bar"},
		},
		{
			name:     "No cities exists.",
			fields:   CityStore{},
			expected: []string{},
		},
	}

	for _, val := range testVals {
		s.Run(val.name, func() {
			mapobj := Map[City, Direction]{Cities: val.fields}
			res := mapobj.GetCitiesNames()
			s.ElementsMatch(res, val.expected)
		})

	}
}

func (s *MapTestSuite) TestAddPath() {
	testVals := []struct {
		name     string
		mapBuild func() (*Map[City, Direction], error)
		withErr  bool
		err      error
	}{
		{
			name: "adds 1 path to north",
			mapBuild: func() (*Map[City, Direction], error) {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("test1")
				m.getOrCreateCity("test2")
				m.getOrCreateCity("test2")
				err := m.AddPath("test1", "test2", "north")
				return &m, err
			},
			withErr: false,
		},
		{
			name: "adds 1 path to city that does not exists.",
			mapBuild: func() (*Map[City, Direction], error) {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("test1")
				m.getOrCreateCity("test2")
				m.getOrCreateCity("test3")
				err := m.AddPath("test1", "test4", "north")
				return &m, err
			},
			withErr: true,
			err:     ErrorCityDoesNotExists,
		},
		{
			name: "adds 1 path from city that does not exists.",
			mapBuild: func() (*Map[City, Direction], error) {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("test1")
				m.getOrCreateCity("test2")
				m.getOrCreateCity("test3")
				err := m.AddPath("a bad city", "test3", "south")
				return &m, err
			},
			withErr: true,
			err:     ErrorCityDoesNotExists,
		},
		{
			name: "adds 1 path from city on an invalid direction.",
			mapBuild: func() (*Map[City, Direction], error) {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("test1")
				m.getOrCreateCity("test2")
				m.getOrCreateCity("test3")
				err := m.AddPath("test1", "test3", "wrong dir")
				return &m, err
			},
			withErr: true,
			err:     ErrorInvalidDirection,
		},
	}

	for _, val := range testVals {
		s.Run(val.name, func() {
			_, err := val.mapBuild()
			if val.withErr {
				s.NotNil(err)
				s.EqualError(err, val.err.Error())
			} else {
				s.Nil(err)
			}
		})

	}
}

func (s *MapTestSuite) TestBuildPathFromToken() {
	testVals := []struct {
		name    string
		fields  string
		dir     string
		err     error
		withErr bool
	}{
		{
			name:    "builds path to north dir",
			fields:  "north=mycity",
			dir:     "north",
			withErr: false,
		},
		{
			name:    "builds path to south dir",
			fields:  "south=mycity",
			dir:     "south",
			withErr: false,
		},
		{
			name:    "builds path to east dir",
			fields:  "east=mycity",
			dir:     "east",
			withErr: false,
		},
		{
			name:    "builds path to west dir",
			fields:  "west=mycity",
			dir:     "west",
			withErr: false,
		},
		{
			name:    "builds path to same city",
			fields:  "west=existingcity",
			dir:     "west",
			err:     ErrorPathToSameCity,
			withErr: true,
		},
		{
			name:    "builds path to an invalid direction",
			fields:  "invaliddir=mycity",
			dir:     "",
			err:     ErrorInvalidDirection,
			withErr: true,
		},
		{
			name:    "builds path with an invalid text syntax.",
			fields:  "invaliddirmycity",
			dir:     "",
			err:     InvalidToken,
			withErr: true,
		},
	}

	for _, values := range testVals {
		s.Run(values.name, func() {
			mapobj := Map[City, Direction]{
				Cities: CityStore{},
				Graph:  graph.NewGraph[City, Direction](),
			}
			mapobj.getOrCreateCity("existingcity")
			err := mapobj.buildPathFromToken(mapobj.Cities["existingcity"], values.fields)
			if values.withErr {
				fmt.Println(values)
				s.NotNil(err)
				s.EqualError(err, values.err.Error())
			} else {
				s.Nil(err)
			}
		})

	}
}

func (s *MapTestSuite) TestGetOrCreateCity() {
	mapobj := Map[City, Direction]{
		Cities: CityStore{},
		Graph:  graph.NewGraph[City, Direction](),
	}
	mapobj.getOrCreateCity("existingcity")
	mapobj2 := Map[City, Direction]{
		Cities: CityStore{},
		Graph:  graph.NewGraph[City, Direction](),
	}
	mapobj2.getOrCreateCity("existingcity")
	testVals := []struct {
		name      string
		cityName  string
		mapobj    Map[City, Direction]
		err       error
		withErr   bool
		numCities int
	}{
		{
			name:      "tests city that already exists",
			cityName:  "existingcity",
			mapobj:    mapobj,
			withErr:   false,
			numCities: 1,
		},
		{
			name:      "tests city that does not exists",
			cityName:  "mycitynew",
			mapobj:    mapobj2,
			withErr:   false,
			numCities: 2,
		},
	}

	for _, values := range testVals {
		s.Run(values.name, func() {
			city := values.mapobj.getOrCreateCity(values.cityName)
			s.Equal(city.Name, values.cityName)
			s.Equal(len(values.mapobj.Cities), values.numCities)

		})
	}
}

func (s *MapTestSuite) TestToString() {

	testVals := []struct {
		name        string
		mapobjBuild func() Map[City, Direction]
		result      string
	}{
		{
			name: "string convert 3 cities and 1 path to north",
			mapobjBuild: func() Map[City, Direction] {
				mapobj := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				mapobj.getOrCreateCity("city1")
				mapobj.getOrCreateCity("city2")
				mapobj.getOrCreateCity("city3")
				mapobj.AddPath("city1", "city2", "north")
				return mapobj
			},
			result: "city1 north=city2\n",
		},
		{
			name: "string convert 3 cities and paths on all cities",
			mapobjBuild: func() Map[City, Direction] {
				mapobj2 := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				mapobj2.getOrCreateCity("city1")
				mapobj2.getOrCreateCity("city2")
				mapobj2.getOrCreateCity("city3")
				mapobj2.AddPath("city1", "city3", "south")
				mapobj2.AddPath("city1", "city2", "west")
				mapobj2.AddPath("city2", "city3", "east")
				mapobj2.AddPath("city3", "city1", "north")
				return mapobj2
			},
			result: "city1 west=city2 south=city3\ncity2 east=city3\ncity3 north=city1\n",
		},
		{
			name: "string convert 5 cities and paths on all directions in 1 city",
			mapobjBuild: func() Map[City, Direction] {
				mapobj3 := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				mapobj3.getOrCreateCity("city1")
				mapobj3.getOrCreateCity("city2")
				mapobj3.getOrCreateCity("city3")
				mapobj3.getOrCreateCity("city4")
				mapobj3.getOrCreateCity("city5")
				mapobj3.AddPath("city1", "city2", "east")
				mapobj3.AddPath("city1", "city3", "north")
				mapobj3.AddPath("city1", "city4", "south")
				mapobj3.AddPath("city1", "city5", "west")
				return mapobj3
			},
			result: "city1 east=city2 north=city3 south=city4 west=city5\n",
		},
	}

	for _, values := range testVals {
		s.Run(values.name, func() {
			m := values.mapobjBuild()
			s.Equal(m.ToString(), values.result)
		})
	}
}

func (s *MapTestSuite) TestNewMapFromReader() {
	vals := []struct {
		name             string
		reader           io.Reader
		hasErr           bool
		err              error
		buildExpectedMap func() Map[City, Direction]
	}{
		{
			name:   "reader with 1 city",
			reader: strings.NewReader("city1\n"),
			hasErr: false,
			err:    nil,
			buildExpectedMap: func() Map[City, Direction] {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("city1")
				return m
			},
		},
		{
			name:   "reader with 1 city and 4 dirs",
			reader: strings.NewReader("city1 north=a south=b east=c west=d\n"),
			hasErr: false,
			err:    nil,
			buildExpectedMap: func() Map[City, Direction] {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("city1")
				m.getOrCreateCity("a")
				m.getOrCreateCity("b")
				m.getOrCreateCity("c")
				m.getOrCreateCity("d")
				m.AddPath("city1", "a", "north")
				m.AddPath("city1", "b", "south")
				m.AddPath("city1", "c", "east")
				m.AddPath("city1", "d", "west")
				return m
			},
		},
		{
			name:   "reader with invalid cities",
			reader: strings.NewReader("      "),
			hasErr: true,
			buildExpectedMap: func() Map[City, Direction] {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				return m
			},
		},
		{
			name:   "reader with invalid dirs",
			reader: strings.NewReader("city1 baddir=a"),
			hasErr: true,
			buildExpectedMap: func() Map[City, Direction] {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("city1")
				m.getOrCreateCity("a")
				return m
			},
		},
		{
			name:   "reader with invalid syntax for paths",
			reader: strings.NewReader("city1 baddir==a"),
			hasErr: true,
			buildExpectedMap: func() Map[City, Direction] {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("city1")
				m.getOrCreateCity("a")
				return m
			},
		},
	}
	for _, val := range vals {
		s.Run(val.name, func() {
			mapFromReader, err := NewMapFromReader(val.reader)
			if val.hasErr {
				s.NotNil(err)
			} else {
				s.Nil(err)
				m := val.buildExpectedMap()
				s.Equal(mapFromReader.Cities, m.Cities)
				s.Equal(mapFromReader.Graph.GetNodes(), m.Graph.GetNodes())
				s.Equal(mapFromReader.Graph.GetEdges(), m.Graph.GetEdges())
			}
		})
	}
}

func (s *MapTestSuite) TestInverseMapper() {
	vals := []struct {
		name   string
		result string
	}{
		{
			name:   "north",
			result: "south",
		},
		{
			name:   "south",
			result: "north",
		},
		{
			name:   "east",
			result: "west",
		},
		{
			name:   "west",
			result: "east",
		},
	}
	for _, val := range vals {
		s.Run(val.name, func() {
			res := InverseMapper(Direction(val.name))
			s.Equal(Direction(val.result), res)
		})
	}
}

func (s *MapTestSuite) TestDestroyCity() {
	vals := []struct {
		name          string
		destroyedCity string
		withErr       bool
		mapBuild      func() Map[City, Direction]
		mapExpected   func() Map[City, Direction]
	}{
		{
			name:          "destroy city with no paths",
			destroyedCity: "my city1",
			withErr:       false,
			mapBuild: func() Map[City, Direction] {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("my city1")
				m.getOrCreateCity("my city2")
				return m
			},
			mapExpected: func() Map[City, Direction] {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				m.getOrCreateCity("my city2")
				return m
			},
		},
		{
			name:          "destroy city that does not exists",
			destroyedCity: "myunexistentcity",
			withErr:       true,
			mapBuild: func() Map[City, Direction] {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				return m
			},
			mapExpected: func() Map[City, Direction] {
				m := Map[City, Direction]{
					Cities: CityStore{},
					Graph:  graph.NewGraph[City, Direction](),
				}
				return m
			},
		},
	}
	for _, val := range vals {
		s.Run(val.name, func() {
			mBefore := val.mapBuild()
			mExpected := val.mapExpected()

			if val.withErr {
				err := mBefore.DestroyCity(&City{Name: val.destroyedCity})
				s.NotNil(err)
			} else {
				err := mBefore.DestroyCity(mBefore.Cities[val.destroyedCity])
				s.Nil(err)
				s.Equal(mBefore.Cities, mExpected.Cities)
				s.Equal(mBefore.Graph.GetNodes(), mExpected.Graph.GetNodes())
				s.Equal(mBefore.Graph.GetEdges(), mExpected.Graph.GetEdges())
			}

		})
	}
}
