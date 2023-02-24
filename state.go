package main

import (
	"math/rand"
	"sort"
)

// State represents a game state.
type State struct {
	Tiles    map[XY]Tile
	Entities map[ID]Entity
	at       map[XY][]ID
}

// NewState returns a new empty State.
func NewState() *State {
	return &State{
		Tiles:    map[XY]Tile{},
		Entities: map[ID]Entity{},
		at:       map[XY][]ID{},
	}
}

// Update updates all entities except from the specified ones.
func (s *State) Update(except ...ID) {
	for id, e := range s.Entities {
		if !in(except, id) {
			e.Update()
		}
	}
	s.at = map[XY][]ID{}
	for id, e := range s.Entities {
		p := e.Pos()
		s.at[p] = append(s.at[p], id)
	}
}

// Add adds the specified entity to the world.
// Returns a new ID of the added entity.
func (s *State) Add(e Entity) ID {
	id := MakeID()
	s.Entities[id] = e
	return id
}

// EntitiesAt returns all entities at the given position sorted by ID
// in increasing order.
func (s *State) EntitiesAt(p XY) []Entity {
	ids := []ID{}
	for _, id := range s.at[p] {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	e := make([]Entity, 0, len(ids))
	for _, id := range ids {
		e = append(e, s.Entities[id])
	}
	return e
}

// RandomPosition returns a random unoccupied position.
func (s *State) RandomPosition() XY {
	empty := []XY{}
	for xy, tile := range s.Tiles {
		if tile == Floor {
			empty = append(empty, xy)
		}
	}
	return empty[rand.Intn(len(empty))]
}
