package aliemsim

import (
	"alien-invasion-simulator/pkg/aliemsim/types"
	"github.com/goombaio/namegenerator"
	"log"
	"time"
)

func spawnAliens(numAliens int, mapObj *types.Map[types.City, types.Direction]) []*types.Alien {
	var result []*types.Alien
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	cityKeys := mapObj.GetCitiesNames()
	for i := 0; i < numAliens; i++ {
		name := nameGenerator.Generate()
		cityName := getRandomItem(cityKeys)
		alien := types.NewAlien(i, name, cityName, mapObj)
		result = append(result, &alien)
	}
	return result
}

var newMapFromReader = types.NewMapFromReader

// StartSimulation starts the alien invasion simulation. By parsing text file and building the types.AlienSimulator object
func StartSimulation(filePath string, numAliens int, fs FileSystem, maxIterations int, verbose bool) error {
	log.Printf("Starting Invasion with: %d aliens...", numAliens)
	log.Printf("Building Map from: %s...", filePath)
	file, err := fs.Open(filePath)

	if err != nil {
		return err
	}

	mapObj, err := newMapFromReader(file)
	if err != nil {
		return err
	}
	aliens := spawnAliens(numAliens, mapObj)
	sim := types.NewAlienSimulator(mapObj, aliens, maxIterations, verbose)

	sim.SimulateInvasion()
	return nil
}
