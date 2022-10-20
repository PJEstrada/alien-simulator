package types

import (
	"alien-invasion-simulator/pkg/graph"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// AlienSimulator is the struct storing the state of cities and alien movements.
type AlienSimulator struct {
	Map                      *Map[City, Direction]
	Aliens                   []*Alien
	AliensInCity             map[*City][]*Alien
	NumDeadAliens            int
	MaxIterations            int
	NumAliensReachedMaxMoves int
	CurrentIteration         int
	NumAliensCannotMove      int
	Verbose                  bool
	Seed                     int64
}

// NewAlienSimulator creates a new alien invasion simulator.
func NewAlienSimulator(mapData *Map[City, Direction], aliens []*Alien, maxIterations int, verbose bool) AlienSimulator {
	return AlienSimulator{
		Map:                      mapData,
		Aliens:                   aliens,
		MaxIterations:            maxIterations,
		Verbose:                  verbose,
		NumAliensReachedMaxMoves: 0,
		NumDeadAliens:            0,
		NumAliensCannotMove:      0,
		CurrentIteration:         0,
	}
}

// printStats logs string information of the current simulation object.
func (sim *AlienSimulator) getStats() string {
	res := fmt.Sprintf("Iteration # %d - Total Aliens: %d - Dead Aliens: %d - Cities Left: %d - Num Aliens Reached Max Moves: %d",
		sim.CurrentIteration,
		len(sim.Aliens),
		sim.NumDeadAliens,
		len(sim.Map.Cities),
		sim.NumAliensReachedMaxMoves,
	)
	return res
}

// SimulateInvasion iterates and moves aliens until all cities are destroyed or until MaxIterations number is reached.
func (sim *AlienSimulator) SimulateInvasion() error {
	sim.CurrentIteration = 0
	if len(sim.Map.Cities) == 0 {
		log.Print("No cities on map. Stopping")
		return nil
	}
	for true {
		if sim.Verbose {
			stats := sim.getStats()
			log.Printf("%s", stats)

		}

		for _, alien := range sim.Aliens {
			if sim.Verbose {
				log.Printf("%s", alien.ToString())
			}
			if alien.IsDead {
				continue
			}
			if alien.NumMovements == sim.MaxIterations && alien.CanMove {
				sim.NumAliensReachedMaxMoves += 1
				sim.NumAliensCannotMove += 1
				alien.CanMove = false
				continue
			}
			// each alien randomly decides to invade a city.

			willMove := sim.alienWillMove()
			if willMove {
				_, err := sim.alienMove(alien)
				if err != nil {
					return err
				}
			}

		}
		if len(sim.Map.Cities) == 0 {
			log.Print("No more cities left. Stopping simulation.")
			break
		}
		if sim.NumDeadAliens == len(sim.Aliens) {
			log.Print("All aliens are dead. Stopping simulation.")
			break
		}
		if sim.NumAliensReachedMaxMoves >= len(sim.Aliens) {
			log.Printf("All aliens have moved %d times. Stopping simulation.", sim.MaxIterations)
			break
		}
		if sim.NumAliensCannotMove == len(sim.Aliens) {
			log.Print("All aliens done. Stopping simulation.")
			break
		}
		sim.CurrentIteration += 1
	}
	log.Printf("Finished Simulation. Map is: \n------- \n\n%s \n", sim.Map.ToString())
	stats := sim.getStats()
	log.Printf("%s", stats)
	return nil
}

// alienWillMove decides randomly if an alien will move.
func (sim *AlienSimulator) alienWillMove() bool {
	seed := time.Now().UTC().UnixNano()
	rand.Seed(seed)
	return rand.Intn(2) == 1
}

// alienMove simulates an alien movement, destroying a city and aliens if more than 2 aliens collide.
func (sim *AlienSimulator) alienMove(alien *Alien) (*City, error) {
	city := sim.Map.Cities[alien.CurrentCityName]

	trapped, err := alien.isTrapped()
	if err != nil {
		return nil, err
	}
	if trapped {
		if sim.Verbose {
			log.Printf("Alien %s is trapped on %v! [Movement #%d]", alien.Name, alien.CurrentCityName, alien.NumMovements)
		}
		alien.NumMovements += 1
		return city, nil
	}
	paths, _ := sim.Map.GetPaths(city)
	var pathKeys []graph.VertexID
	for k, _ := range paths {
		pathKeys = append(pathKeys, k)
	}
	chosenPathKey := pathKeys[rand.Intn(len(pathKeys))]
	chosenPath := paths[chosenPathKey]
	prevCity := alien.CurrentCityName
	invadedCity := alien.invadeCity(&chosenPath.To.Data)
	if len(invadedCity.Aliens) == 2 {
		log.Printf("paths are %v. chose path was %v. current city of alien is %s", pathKeys, chosenPath, prevCity)
		aliensBefore := []*Alien{}
		aliensBefore = append(aliensBefore, invadedCity.Aliens[0])
		aliensBefore = append(aliensBefore, invadedCity.Aliens[1])
		sim.Map.DestroyCity(invadedCity)
		log.Printf("[DESTROYED] Aliens %s and %s are fighting! City %s is destroyed.",
			aliensBefore[0].Name,
			aliensBefore[1].Name,
			invadedCity.Name)
		alien.IsDead = true
		sim.NumDeadAliens += 2
		sim.NumAliensCannotMove += 2

	}
	return invadedCity, nil
}
