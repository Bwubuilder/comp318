package database

import (
	"encoding/json"
	"log/slog"

	"github.com/Bwubuilder/owldb/skiplist"
)

// The Collection struct represents a collection in a database.
type Collection struct {
	Name      string
	Documents skiplist.SkipList[string, Document]
	URI       string `json:"uri"`
}

// NewCollection creates and returns a new Collection struct with the given name.
func NewCollection(name string) *Collection {
	var col Collection
	col.Name = name
	col.Documents = skiplist.NewSkipList[string, Document]()
	col.URI = name
	return &col
}

// GetChildByName implements the function from the PathItem interface.
// If it exists, it returns the document and true, otherwise nil and false.
func (c *Collection) GetChildByName(name string) (PathItem, bool) {
	child, exists := c.Documents.Find(name)
	if exists {
		return &child, true
	}
	return nil, false
}

// Marshal implements the function from the PathItem interface.
// Calling Marshal() marshals and returns the collection as well as an error.
func (c Collection) Marshal() ([]byte, error) {
	colURI := map[string]string{"uri": c.URI}

	response, err := json.MarshalIndent(colURI, "", "  ")
	if err != nil {
		slog.Info("Collection marshaling failed")
		return nil, err
	}
	return response, nil
}
