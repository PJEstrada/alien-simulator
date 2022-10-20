package types

import "fmt"

// Alien represents an alien that can invade cities and track the number of movements and current city.
type Alien struct {
	ID              int
	Name            string
	CurrentCityName string
	Map             *Map[City, Direction]
	IsDead          bool
	NumMovements    int
	CanMove         bool
}

func (a *Alien) ToString() string {
	s := fmt.Sprintf("Alien[%d] - Name: %s - CurrentCity: %s - Dead: %v - NumMovements: %d - Can Move: %v",
		a.ID,
		a.Name,
		a.CurrentCityName,
		a.IsDead,
		a.NumMovements,
		a.CanMove)
	return s
}

// NewAlien create a new alien on with given name, city and ID.
func NewAlien(id int, name string, city string, mapObj *Map[City, Direction]) Alien {
	return Alien{
		ID:              id,
		Name:            name,
		CurrentCityName: city,
		NumMovements:    0,
		Map:             mapObj,
		CanMove:         true,
		IsDead:          false,
	}
}

// isTrapped return true if alien has no paths to go from its CurrentCityName otherwise returns false.
func (a *Alien) isTrapped() (bool, error) {
	city := a.Map.Cities[a.CurrentCityName]
	if city == nil {
		return true, ErrorCityDoesNotExists
	}
	paths, _ := a.Map.GetPaths(city)
	if len(paths) == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// invadeCity changes the current city of an alien, adds the alien to the city obj, and increments the number of movements counter.
func (a *Alien) invadeCity(c *City) *City {
	if !c.hasAlien(a) {
		c.Aliens = append(c.Aliens, a)
		a.CurrentCityName = c.Name
		a.NumMovements += 1
	}
	return c
}
