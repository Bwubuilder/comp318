package jsontogo

import (
	"fmt"
	"log/slog"

	"github.com/Bwubuilder/owldb/jsonvisitor/jsonvisit"
)

type jsonVisitor struct {
}

func New() jsonVisitor {
	return jsonVisitor{}
}

// Process JSON Map by iterating through map and calling Accept
func (v jsonVisitor) Map(m map[string]any) (any, error) {
	var returnMap map[string]any

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
	var returnSlice []any
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
	if b {
		return true, nil
	} else {
		return false, nil
	}
}

// Process JSON float
func (v jsonVisitor) Float64(f float64) (any, error) {
	return f, nil
}

// Process JSON string
func (v jsonVisitor) String(s string) (any, error) {
	return fmt.Sprintf("\"%s\"", s), nil
}

// Process JSON null value
func (v jsonVisitor) Null() (any, error) {
	slog.Info("Call to null")
	return nil, nil
}
