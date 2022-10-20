package types

// City structure that contains the city name and all the possible paths to other cities. Might also have aliens!
type City struct {
	Name        string
	Aliens      []*Alien
	isDestroyed bool
}

// hasAlien determines if the city has the given alien.
func (c *City) hasAlien(a2 *Alien) bool {
	for _, a := range c.Aliens {
		if a.ID == a2.ID {
			return true
		}
	}
	return false
}

// Destroy kills all aliens on the city and marks the city as destroyed
func (c *City) Destroy() {
	for _, alien := range c.Aliens {
		alien.IsDead = true
	}
	c.isDestroyed = true
	c.Aliens = []*Alien{}
}

// NewCityFromName creates a new city struct from a string name
func NewCityFromName(name string) City {
	return City{
		Name:        name,
		isDestroyed: false,
	}
}

func (c City) ID() string {
	return c.Name
}
