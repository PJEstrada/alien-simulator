# Alien Invasion Simulator

Cli tool for an alien invasion simulation. Simulates aliens moving randomly, fighting
and destroying cities.

## Installation

1. Clone Repo

```
git clone https://github.com/PJEstrada/alien-simulator
```

2. Build Binaries.

```
cd aliensim && go install
```


3. Run CLI command With sample file with 3 aliens

```
alien-invasion-simulator sampleMapFiles/cities1.txt 3 
```


3. Optional, run verbose mode.

```
alien-invasion-simulator sampleMapFiles/cities1.txt 3 --verbose
```

### Assumptions

* I'm modeling the city as a directed graph with the constraint that
  when an edge is created on one city only in that direction. 
* Paths cannot be duplicate and should be explicitly defined
  on text file. A path from A => B on north does not mean a path from B => A on
  south exists.
* An alien that is trapped, tries to move on each iteration, so it counts as a movement.
* A city can only receive 1 alien at a time. Otherwise, aliens would have perfect arrival timing.
* Duplicate city names are not supported.
* I assume aliens arrival to cities are instantaneous.
