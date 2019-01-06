package storage

import "errors"

type Memory struct {
	last  string
	isSet bool
}

// Save saves the last position into memory variable. This is not restart save.
func (m *Memory) Save(position string) error {
	m.last = position
	m.isSet = true
	return nil
}

// Last returns the last read position or an error if not set yet.
func (m *Memory) Last() (string, error) {
	if !m.isSet {
		return "", errors.New("not set during this program life time")
	}
	return m.last, nil
}
