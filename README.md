# golang simple full text search engine

> original article: https://artem.krylysov.com/blog/2020/07/28/lets-build-a-full-text-search-engine/

## Corpus
We are going to search a part of the abstract of English Wikipedia.
The latest dump is available at dumps.wikimedia.org. As of today, the file size after decompression is 913 MB.
The XML file contains over 600K documents.

Document example:

```
<title>Wikipedia: Kit-Cat Klock</title>
<url>https://en.wikipedia.org/wiki/Kit-Cat_Klock</url>
<abstract>The Kit-Cat Klock is an art deco novelty wall clock shaped like a grinning cat with cartoon eyes that swivel in time with its pendulum tail.</abstract>
```

### Loading documents

`pkg/wikipedia/dump/dump.go`

```go
type Document struct {
	Title string `xml:"title"`
	URL   string `xml:"url"`
	Text  string `xml:"abstract"`
	ID    int
}

func LoadDocument(filename string) ([]Document, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := xml.NewDecoder(f)
	dump := struct {
		Documents []Document `xml:"doc"`
	}{}
	if err := dec.Decode(&dump); err != nil {
		return nil, err
	}

	docs := dump.Documents
	for i, _ := range docs {
		docs[i].ID = i
	}
	return docs, nil
}
```

## Analyzing documents


`pkg/search/analyze.go`

```go
func Analyze(text string) []string {
    tokens := tokenize(text)
    tokens = lowercaseFilter(tokens)
    tokens = stopwordFilter(tokens)
    tokens = stemmerFilter(tokens)
    return tokens
}
```

### Tokenize

The tokenizer is the first step of text analysis. Its job is to convert text into a list of tokens.
Our implementation splits the text on a word boundary and removes punctuation marks:

`pkg/search/tokenize.go`

```go
func tokenize(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		// Split on any character that is not a letter or a number.
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}
```

### lowercase

In order to make the search case-insensitive, the lowercase filter converts tokens to lower case.
cAt, Cat and caT are normalized to cat. Later, when we query the index,
we'll lower case the search terms as well. This will make the search term cAt match the text Cat.

`pkg/search/filters.go`

```go
func lowercaseFilter(tokens []string) []string {
	r := make([]string, len(tokens))
	for i, token := range tokens {
		r[i] = strings.ToLower(token)
	}
	return r
}
```

### stopword

Almost any English text contains commonly used words like a, I, the or be. Such words are called stop words. We are going to remove them since almost any document would match the stop words.


`pkg/search/filters.go`

```go
var stopwords = map[string]struct{}{
	"a": {}, "and": {}, "be": {}, "have": {}, "i": {},
	"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
}

func stopwordFilter(tokens []string) []string {
	r := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := stopwords[token]; !ok {
			r = append(r, token)
		}
	}
	return r
}
```

### stemming

Because of the grammar rules, documents may include different forms of the same word.
Stemming reduces words into their base form.
For example, fishing, fished and fisher may be reduced to the base form (stem) fish.

Implementing a stemmer is a non-trivial task, We'll take one of the existing modules:


`pkg/search/filters.go`

```go
import snowballeng "github.com/kljensen/snowball/english"

func stemmerFilter(tokens []string) []string {
    r := make([]string, len(tokens))
    for i, token := range tokens {
        r[i] = snowballeng.Stem(token, false)
    }
    return r
}
```
