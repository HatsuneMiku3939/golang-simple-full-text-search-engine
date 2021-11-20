package dump

import (
	"encoding/xml"
	"os"

	"fts/pkg/wikipedia"
)

// LoadDocument loads a Wikipedia article from a dump file.
func LoadDocument(filename string) ([]wikipedia.Document, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := xml.NewDecoder(f)
	dump := struct {
		Documents []wikipedia.Document `xml:"doc"`
	}{}
	if err := dec.Decode(&dump); err != nil {
		return nil, err
	}

	docs := dump.Documents
	for i := range docs {
		docs[i].ID = i
	}
	return docs, nil
}
