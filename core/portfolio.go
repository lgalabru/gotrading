package core

import (
	"sync"

	uuid "github.com/satori/go.uuid"
)

var instance *PortfolioStateManager
var once sync.Once

func SharedPortfolioManager() *PortfolioStateManager {
	once.Do(func() {
		instance = &PortfolioStateManager{}
		instance.States = make(map[string]PortfolioState)
	})
	return instance
}

type PortfolioStateManager struct {
	States      map[string]PortfolioState
	LastStateID string
}

// Portfolio wraps all your positions
type PortfolioState struct {
	StateID   string                          `json:"stateID"`
	Positions map[string]map[Currency]float64 `json:"positions"`
}

// Portfolio wraps all your positions
type PortfolioStateSlice struct {
	Exch      string
	Positions map[Currency]float64
}

// Portfolio wraps all your positions
type Portfolio struct {
	Positions map[string]map[Currency]float64
}

// Update position
func (m *PortfolioStateManager) PushState(state PortfolioState) {
	m.LastStateID = state.StateID
	m.States[state.StateID] = state
}

func NewPortfolioStateFromPositions(positions map[string]map[Currency]float64) PortfolioState {
	state := NewPortfolioState()
	state.Positions = positions
	return state
}

func NewPortfolioState() PortfolioState {
	state := PortfolioState{}
	uuid := (uuid.NewV4()).String()
	state.StateID = uuid
	state.Positions = make(map[string]map[Currency]float64)
	return state
}

func (m *PortfolioStateManager) LastPositions() map[string]map[Currency]float64 {
	return m.States[m.LastStateID].Positions
}

func (m *PortfolioStateManager) Position(stateID, exch string, curr Currency) float64 {
	return m.States[stateID].Positions[exch][curr]
}

func (m *PortfolioStateManager) CurrentPosition(exch string, curr Currency) float64 {
	return m.States[m.LastStateID].Positions[exch][curr]
}

// Update position
func (m *PortfolioStateManager) UpdateWithNewState(state PortfolioState, override bool) {
	if override || len(m.States) == 0 {
		m.PushState(state)
	} else {
		last := m.States[m.LastStateID]
		new := NewPortfolioState()
		new.StateID = (uuid.NewV4()).String()

		for exch := range last.Positions {
			new.Positions[exch] = make(map[Currency]float64)
			for currency := range state.Positions[exch] {
				new.Positions[exch][currency] = last.Positions[exch][currency]
			}
		}
		for exch := range state.Positions {
			for currency := range state.Positions[exch] {
				new.Positions[exch][currency] = state.Positions[exch][currency]
			}
		}
		m.PushState(new)
	}
}

// Update position
func (m *PortfolioStateManager) UpdateWithNewPosition(exch string, c Currency, amount float64) {
	current := m.States[m.LastStateID]
	next := current.Copy()
	next.UpdatePosition(exch, c, amount)
	uuid := (uuid.NewV4()).String()
	next.StateID = uuid
	m.PushState(next)
}

// Fork current state
func (m *PortfolioStateManager) ForkCurrentState() PortfolioState {
	current := m.States[m.LastStateID]
	fork := current.Copy()
	uuid := (uuid.NewV4()).String()
	fork.StateID = uuid
	return fork
}

// Update position
func (s *PortfolioState) UpdatePosition(exch string, c Currency, amount float64) {
	if s.Positions == nil {
		s.Positions = make(map[string]map[Currency]float64)
	}
	if s.Positions[exch] == nil {
		s.Positions[exch] = make(map[Currency]float64)
	}
	s.Positions[exch][c] = Trunc8(amount)
}

// Copy state
func (s *PortfolioState) Copy() PortfolioState {
	copy := PortfolioState{}
	copy.Positions = make(map[string]map[Currency]float64)

	for exch := range s.Positions {
		copy.Positions[exch] = make(map[Currency]float64)
		for currency := range s.Positions[exch] {
			copy.Positions[exch][currency] = s.Positions[exch][currency]
		}
	}
	return copy
}

// Update position
func (s *Portfolio) UpdatePosition(exch string, c Currency, amount float64) {
	if s.Positions == nil {
		s.Positions = make(map[string]map[Currency]float64)
	}
	if s.Positions[exch] == nil {
		s.Positions[exch] = make(map[Currency]float64)
	}
	s.Positions[exch][c] = Trunc8(amount)
}
