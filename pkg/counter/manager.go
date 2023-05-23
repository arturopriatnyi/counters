package counter

import "errors"

type Manager struct {
	s Storage
}

func NewManager(s Storage) *Manager {
	return &Manager{s: s}
}

var ErrExists = errors.New("counter exists")

func (m *Manager) Add(id string) error {
	_, err := m.s.Get(id)
	if err == nil {
		return ErrExists
	}
	if err != ErrNotFound {
		return err
	}

	return m.s.Set(&Counter{ID: id})
}

func (m *Manager) Get(id string) (*Counter, error) {
	return m.s.Get(id)
}

func (m *Manager) Inc(id string) error {
	counter, err := m.s.Get(id)
	if err != nil {
		return err
	}

	counter.Inc()

	return m.s.Set(counter)
}

func (m *Manager) Delete(id string) error {
	return m.s.Delete(id)
}
