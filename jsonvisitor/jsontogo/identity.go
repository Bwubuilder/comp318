package jsontogo

import (
	"log/slog"
)

type identity struct{}

func NewID() identity {
	return identity{}
}

// Create map
func (v identity) Map(m map[string]any) (any, error) {
	returnMap := m
	return returnMap, nil
}

func (v identity) Slice(s []any) (any, error) {
	returnSlice := s
	return returnSlice, nil
}

func (v identity) Bool(b bool) (any, error) {
	returnBool := b
	return returnBool, nil
}

func (v identity) Float64(f float64) (any, error) {
	returnFloat := f
	return returnFloat, nil
}

func (v identity) String(s string) (any, error) {
	returnString := s
	return returnString, nil
}

func (v identity) Null() (any, error) {
	slog.Info("Call to null")
	return nil, nil
}
