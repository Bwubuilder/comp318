package jsonstringer

import (
	"fmt"

	"github.com/Bwubuilder/owldb/jsonvisitor/jsonvisit"
)

type jsonStringVisitor struct {
}

func New() jsonStringVisitor {
	return jsonStringVisitor{}
}

// Process JSON Map by iterating through map and calling Accept
func (v jsonStringVisitor) Map(m map[string]any) (string, error) {
	str := fmt.Sprint("{")
	first := true // see if comma needs to be added before item
	for key, val := range m {
		if !first {
			str += ", "
		} else {
			first = false
		}
		str += fmt.Sprintf("\"%s\": ", key)
		res, err := jsonvisit.Accept(val, v)
		if err != nil {
			return "", err
		}
		str += res
	}
	str += "}"

	return str, nil
}

// Process JSON slice by iterating through slice and calling Accept
func (v jsonStringVisitor) Slice(s []any) (string, error) {
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
func (v jsonStringVisitor) Bool(b bool) (string, error) {
	if b {
		return "true", nil
	} else {
		return "false", nil
	}
}

// Process JSON float
func (v jsonStringVisitor) Float64(f float64) (string, error) {
	return fmt.Sprintf("%f", f), nil
}

// Process JSON string
func (v jsonStringVisitor) String(s string) (string, error) {
	return fmt.Sprintf("\"%s\"", s), nil
}

// Process JSON null value
func (v jsonStringVisitor) Null() (string, error) {
	return "null", nil
}
