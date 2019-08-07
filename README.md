# nounphrases

## Introduction

_nounphrases_ allows us a to find noun phrases on a given text. It splits the input text in several sentences and then processes every sentence, extracting the list of noun phrases encountered on them.

The algorithm is port of the noun phrases extractor from the [TextBlob](https://github.com/sloria/TextBlob/blob/master/textblob/en/np_extractors.py#L135) Python's library, which uses POS tagging.

## Usage

The `Find(string)` function will return an array of phrases found on the given text.

```go
package main

import (
    "fmt"
    "github.com/eroatta/nounphrases"
    "log"
)

func main() {
    phrases, err := nounphrases.Find("We have red cars and yellow trucks")
    if err != nil {
        log.Fatal(err)
    }

    for i, phr := range phrases {
        fmt.Println(fmt.Sprintf("Phrase #%d: %s", i, phr))
    }
}

```

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).
