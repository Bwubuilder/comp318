package jsontogo

import (
	"log/slog"

	"github.com/Bwubuilder/owldb/jsonvisitor/jsonvisit"
)

type jsonVisitor struct {
	id identity
}

func New() jsonVisitor {
	return jsonVisitor{id: NewID()}
}

// Process JSON Map by iterating through map and calling Accept
func (v jsonVisitor) Map(m map[string]any) (any, error) {
	returnMap := make(map[string]any)

	for key, val := range m {
		res, err := jsonvisit.Accept(val, v)
		if err != nil {
			return "", err
		}
		slog.Info("Created Key: ", key, "Created Value: ", res)
		returnMap[key] = res
	}

	return returnMap, nil
}

// Process JSON slice by iterating through slice and calling Accept
func (v jsonVisitor) Slice(s []any) (any, error) {
	returnSlice := make([]any, 0, 10)

	for i, val := range s {
		res, err := jsonvisit.Accept(val, v)
		if err != nil {
			return "", err
		}
		slog.Info("Created: ", res, "at index: ", i)
		returnSlice[i] = res
	}

	return returnSlice, nil
}

// Process JSON bool by printing bool
func (v jsonVisitor) Bool(b bool) (any, error) {
	return jsonvisit.Accept(b, v.id)
}

// Process JSON float
func (v jsonVisitor) Float64(f float64) (any, error) {
	return jsonvisit.Accept(f, v.id)
}

// Process JSON string
func (v jsonVisitor) String(s string) (any, error) {
	return jsonvisit.Accept(s, v.id)
}

// Process JSON null value
func (v jsonVisitor) Null() (any, error) {
	return jsonvisit.Accept(nil, v.id)
}
