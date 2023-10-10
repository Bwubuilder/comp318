package authorization

import (
	"fmt"
	"log/slog"

	"github.com/comp318/tutorials/json/visitor/jsonvisit"
)

type authVisitor struct {
}

func New() authVisitor {
	return authVisitor{}
}

// Process JSON Map by iterating through map and calling Accept
func (v authVisitor) Map(m map[string]any) (string, error) {
	var user string
	for key, val := range m {
		slog.Info("This is", key)
		res, err := jsonvisit.Accept(val, v)
		if err != nil {
			return "", err
		}
		user += res
	}

	return user, nil
}

// Process JSON slice by iterating through slice and calling Accept
func (v authVisitor) Slice(s []any) (string, error) {
	str := fmt.Sprint("[")
	first := true // see if comma needs to be added before item
	for _, val := range s {
		if !first {
			str += ", "
		} else {
			first = false
		}
		res, err := jsonvisit.Accept(val, v)
		if err != nil {
			return "", err
		}
		str += res
	}
	str += "]"

	return str, nil
}

// Process JSON bool by printing bool
func (v authVisitor) Bool(b bool) (string, error) {
	if b {
		return "true", nil
	} else {
		return "false", nil
	}
}

// Process JSON float
func (v authVisitor) Float64(f float64) (string, error) {
	return fmt.Sprintf("%f", f), nil
}

// Process JSON string
func (v authVisitor) String(s string) (string, error) {
	return fmt.Sprintf("\"%s\"", s), nil
}

// Process JSON null value
func (v authVisitor) Null() (string, error) {
	return "null", nil
}
