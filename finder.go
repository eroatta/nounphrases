package nounphrases

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/jdkato/prose.v2"
)

var (
	phrasePairs = map[string]string{
		"NNP+NNP": "NNP",
		"NN+NN":   "NNI",
		"NNI+NN":  "NNI",
		"JJ+JJ":   "JJ",
		"JJ+NN":   "NNI",
	}

	word = regexp.MustCompile("[a-zA-z0-9-_]+")

	nounTags = map[string]bool{
		"NN":   true, // noun, singular or mass
		"NNP":  true, // noun, proper singular
		"NNPS": true, // noun, proper plural
		"NNS":  true, // noun, plural
	}
)

// Find looks for phrases on the text body.
func Find(text string) ([]string, error) {
	parsedText, err := prose.NewDocument(text)
	if err != nil {
		return []string{}, err
	}
	sentences := parsedText.Sentences()

	// parallelize sentence processing
	var wg sync.WaitGroup
	wg.Add(len(sentences))
	phrc := make(chan []string)
	errc := make(chan error)

	for _, sentence := range sentences {
		go func(text string) {
			defer wg.Done()
			phrases, err := extract(text)
			if err != nil {
				errc <- err
			}
			phrc <- phrases
		}(sentence.Text)
	}

	// closer function
	go func() {
		wg.Wait()
		close(phrc)
		close(errc)
	}()

	// merge results
	phrases := []string{}
	for i := 0; i < len(sentences); i++ {
		select {
		case phr := <-phrc:
			phrases = append(phrases, phr...)
		case err := <-errc:
			return []string{}, err
		}
	}

	return phrases, nil
}

// extract returns a list of noun phrases (strings) from the sentence.
func extract(sentence string) ([]string, error) {
	parsed, err := prose.NewDocument(sentence)
	if err != nil {
		return []string{}, err
	}

	tokens := make([]prose.Token, 0)
	for _, token := range parsed.Tokens() {
		nt := normalize(token)
		if unwantedNoun(nt) {
			continue
		}

		tokens = append(tokens, nt)
	}

	if len(tokens) < 2 {
		return []string{}, nil
	}

	nounPhrases := []string{}

	tag := tokens[0].Tag
	text := tokens[0].Text
	for i := 1; i < len(tokens); i++ {
		t2 := tokens[i]
		value := phrasePairs[fmt.Sprintf("%s+%s", tag, t2.Tag)]
		if value == "" {
			if tag == "NNI" {
				nounPhrases = append(nounPhrases, text)
			}

			tag = t2.Tag
			text = t2.Text
		} else {
			tag = value
			text = fmt.Sprintf("%s %s", text, t2.Text)

			// check for the last element
			if i == len(tokens)-1 && tag == "NNI" {
				nounPhrases = append(nounPhrases, text)
			}
		}
	}

	return nounPhrases, nil
}

// unwantedNoun checks if the token is a noun and doesn't match the required criteria to
// be considered as a valid noun for the phrase.
func unwantedNoun(tok prose.Token) bool {
	return nounTags[tok.Tag] && !word.MatchString(tok.Text)
}

// normalize normalizes the corpus tags: ("NN", "NN-PL", "NNS") -> "NN"
func normalize(tok prose.Token) prose.Token {
	if tok.Tag == "NP-TL" || tok.Tag == "NP" {
		tok.Tag = "NNP"
		return tok
	}

	if strings.HasSuffix(tok.Tag, "-TL") {
		tok.Tag = strings.TrimSuffix(tok.Tag, "-TL")
		return tok
	}

	if strings.HasSuffix(tok.Tag, "S") {
		tok.Tag = strings.TrimSuffix(tok.Tag, "S")
		return tok
	}

	return tok
}
