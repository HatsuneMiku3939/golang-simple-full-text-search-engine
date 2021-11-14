package main

import (
	"fmt"
	"fts/pkg/wikipedia/dump"
)

const (
	// The path to the Wikipedia dump file.
	dumpFile = "test/testdata/enwiki-latest-abstract1.xml"
)

func main() {
	docs, err := dump.LoadDocument(dumpFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Length: %d\n", len(docs))
	fmt.Printf("Doc[0]: %v\n", docs[1])
}
