package types

import (
	"alien-invasion-simulator/pkg/graph"
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

type CityStore map[string]*City
type Direction string

// StringInSlice return string representation of a Path
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

var ValidDirections = []string{"north", "south", "east", "west"}
var ErrorCityDoesNotExists = errors.New("City does not exists.")
var ErrorInvalidDirection = errors.New("Invalid Direction.")
var ErrorPathToSameCity = errors.New("Cannot create a path to the same city.")

// Map represents a set of cities
type Map[N City, E Direction] struct {
	Cities                 CityStore
	DirectionInverseMapper func(Direction) Direction
	Graph                  graph.Graph[City, Direction]
}

// InverseMapper mapper of the possible directions to is opposite direction
func InverseMapper(d Direction) Direction {
	inverseDir := map[Direction]Direction{
		"north": "south",
		"south": "north",
		"east":  "west",
		"west":  "east",
	}
	return inverseDir[d]
}

// NewMapFromReader create a Map object from the given file reader. Reader should have format: 'city dir=city' per line.
func NewMapFromReader(reader io.Reader) (*Map[City, Direction], error) {
	scanner := bufio.NewScanner(reader)
	mapObj := &Map[City, Direction]{
		Cities:                 CityStore{},
		DirectionInverseMapper: InverseMapper,
		Graph:                  graph.NewGraph[City, Direction](),
	}

	for scanner.Scan() {
		textLine := scanner.Text()
		fields := strings.Fields(textLine)
		err := validTextRow(fields)
		if err != nil {
			return nil, err
		}
		var fromCity *City
		for i, strToken := range fields {
			if i == 0 {
				fromCity = mapObj.getOrCreateCity(strToken)
			} else {
				err := mapObj.buildPathFromToken(fromCity, strToken)
				if err != nil {
					return nil, err
				}
			}
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return mapObj, nil
}

// EdgeToString returns the string representation of an edge
func (m *Map[N, E]) EdgeToString(edge *graph.Edge[City, Direction]) string {
	return fmt.Sprintf("%s=%s", string(edge.Data), edge.To.Data.Name)
}

// ToString returns the map representation as a string
func (m *Map[N, E]) ToString() string {
	result := ""
	keys := make([]string, 0, len(m.Cities))
	for k := range m.Cities {
		keys = append(keys, k)
	}
	sortedKeys := sort.StringSlice(keys)
	sort.Sort(sortedKeys)

	for _, cityName := range sortedKeys {
		city := m.Cities[cityName]
		paths, _ := m.GetPaths(city)
		if len(paths) == 0 {
			continue
		}
		result += city.Name

		keysPaths := make([]string, 0, len(paths))
		for k := range paths {
			keysPaths = append(keysPaths, string(k))
		}
		sortedP := sort.StringSlice(keysPaths)
		sort.Sort(sortedP)
		for _, k := range sortedP {
			val := paths[graph.VertexID(k)]
			result += fmt.Sprintf(" %s", m.EdgeToString(val))
		}
		result += "\n"
	}
	return result
}

// AddPath creates a path on the map between 2 cities on the given direction.
func (m *Map[N, E]) AddPath(fromCityName string, toCityName string, dir Direction) error {
	// TODO ENFORCE ONLY ONE DIRECTION EXISTS PER CITY
	if _, ok := m.Cities[fromCityName]; !ok {
		return ErrorCityDoesNotExists
	}
	if _, ok := m.Cities[toCityName]; !ok {
		return ErrorCityDoesNotExists
	}
	if !StringInSlice(string(dir), ValidDirections) {
		return ErrorInvalidDirection
	}
	// Add edge on graph
	fromCity := m.Graph.GetVertexByStringID(fromCityName)
	toCity := m.Graph.GetVertexByStringID(toCityName)
	_, err := m.Graph.AddEdge(fromCity.Id(), toCity.Id(), dir)
	if err != nil {
		return err
	}
	return nil
}

// DestroyCity marks a city as destroyed and removes the city from the map.
func (m *Map[N, E]) DestroyCity(city *City) error {
	// Remove From Path
	id := m.Graph.GetVertexByStringID(city.Name)
	if id == nil {
		return ErrorCityDoesNotExists
	}
	_ = m.Graph.RemoveVertex(*m.Graph.GetVertexByStringID(city.Name))
	// Remove City from Map
	delete(m.Cities, city.Name)
	// Destroy City
	city.Destroy()
	return nil
}

// GetPaths returns the possible paths from the given city
func (m *Map[N, E]) GetPaths(fromCity *City) (map[graph.VertexID]*graph.Edge[City, Direction], error) {
	vertex := m.Graph.GetVertexByStringID(fromCity.Name)
	if vertex == nil {
		return nil, ErrorCityDoesNotExists
	}

	return vertex.OutgoingEdges, nil

}

// GetCitiesNames returns the available cities on a map as a string slice
func (m *Map[N, E]) GetCitiesNames() []string {
	keys := []string{}
	for key, _ := range m.Cities {
		keys = append(keys, key)
	}
	return keys
}

// validTextRow validates that a string complies with the expected syntax for map parsing
func validTextRow(rowFields []string) error {
	if len(rowFields) == 0 {
		return errors.New(fmt.Sprintf("Invalid row: '%v' Needs to have at least a city name", rowFields))
	}
	if len(rowFields) > 1 {
		for i, field := range rowFields {
			if i == 0 {
				continue
			}
			splitted := strings.Split(field, "=")
			if len(splitted) != 2 {
				return errors.New(fmt.Sprintf("Invalid row: '%v' each path needs to have format 'direction=city'", rowFields))
			}

			if !StringInSlice(splitted[0], ValidDirections) {
				return errors.New(fmt.Sprintf("Invalid row: '%v' direction needs to be one from %v", rowFields, ValidDirections))
			}

		}
	}
	return nil
}

var InvalidToken = errors.New("Invalid token syntax.")

// buildPathFromToken creates an edge from one city into another one in a given direction from token string.
func (m *Map[N, E]) buildPathFromToken(fromCity *City, elm string) error {
	splitted := strings.Split(elm, "=")
	if len(splitted) != 2 {
		return InvalidToken
	}
	dir := splitted[0]
	cityName := splitted[1]
	if cityName == fromCity.Name {
		return ErrorPathToSameCity
	}
	cityObj := m.getOrCreateCity(cityName)
	err := m.AddPath(fromCity.Name, cityObj.Name, Direction(dir))
	if err != nil {
		return err
	}
	return nil
}

// getOrCreateCity finds a city on the map or creates a new one by adding vertex on graph
func (m *Map[N, E]) getOrCreateCity(cityName string) *City {
	if _, ok := m.Cities[cityName]; !ok {

		newCity := NewCityFromName(cityName)
		m.Cities[cityName] = &newCity
		m.Graph.AddVertex(newCity)
		return &newCity
	} else {
		return m.Cities[cityName]
	}
}
