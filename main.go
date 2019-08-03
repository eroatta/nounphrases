package main

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/jdkato/prose.v2"
)

func main() {
	/*
		//Create a chain of order 2
		chain := gomarkov.NewChain(2)

		//Feed in training data
		chain.Add(strings.Split("I want a cheese burger", " "))
		chain.Add(strings.Split("I want a chilled sprite", " "))
		chain.Add(strings.Split("I want to go to the movies", " "))

		//Get transition probability of a sequence
		prob, _ := chain.TransitionProbability("a", []string{"I", "want"})
		fmt.Println(prob)
		//Output: 0.6666666666666666

		//You can even generate new text based on an initial seed
		chain.Add(strings.Split("Mother should I build the wall?", " "))
		chain.Add(strings.Split("Mother should I run for President?", " "))
		chain.Add(strings.Split("Mother should I trust the government?", " "))
		next, _ := chain.Generate([]string{"Mother", "should"})
		fmt.Println(next)

		//The chain is JSON serializable
		jsonObj, _ := json.Marshal(chain)
		err := ioutil.WriteFile("model.json", jsonObj, 0644)
		if err != nil {
			fmt.Println(err)
		}*/

	tryPhraseExtractionFromPOSTagging()
}

func tryPhraseExtractionFromPOSTagging() {
	funcComments := `
	Marshal returns the JSON encoding of v.
	Marshal traverses the value v recursively. If an encountered value implements the Marshaler interface and is not a nil pointer, Marshal calls its MarshalJSON method to produce JSON. If no MarshalJSON method is present but the value implements encoding.TextMarshaler instead, Marshal calls its MarshalText method and encodes the result as a JSON string. The nil pointer exception is not strictly necessary but mimics a similar, necessary exception in the behavior of UnmarshalJSON.
	Otherwise, Marshal uses the following type-dependent default encodings:
	Boolean values encode as JSON booleans.
	Floating point, integer, and Number values encode as JSON numbers.
	String values encode as JSON strings coerced to valid UTF-8, replacing invalid bytes with the Unicode replacement rune. The angle brackets "<" and ">" are escaped to "\u003c" and "\u003e" to keep some browsers from misinterpreting JSON output as HTML. Ampersand "&" is also escaped to "\u0026" for the same reason. This escaping can be disabled using an Encoder that had SetEscapeHTML(false) called on it.
	Array and slice values encode as JSON arrays, except that []byte encodes as a base64-encoded string, and a nil slice encodes as the null JSON value.
	Struct values encode as JSON objects. Each exported struct field becomes a member of the object, using the field name as the object key, unless the field is omitted for one of the reasons given below.
	The encoding of each struct field can be customized by the format string stored under the "json" key in the struct field's tag. The format string gives the name of the field, possibly followed by a comma-separated list of options. The name may be empty in order to specify options without overriding the default field name.
	The \"omitempty\" option specifies that the field should be omitted from the encoding if the field has an empty value, defined as false, 0, a nil pointer, a nil interface value, and any empty array, slice, map, or string.
	As a special case, if the field tag is \"-\", the field is always omitted. Note that a field with name "-" can still be generated using the tag "-,". 
	`
	doc, err := prose.NewDocument(funcComments)
	if err != nil {
		log.Fatal(err)
	}

	nouns := map[string]bool{
		"NN":   true, // noun, singular or mass
		"NNP":  true, // noun, proper singular
		"NNPS": true, // noun, proper plural
		"NNS":  true, // noun, plural
	}

	for i, sent := range doc.Sentences() {
		fmt.Println(fmt.Sprintf("#%d. Line: %s", i, sent.Text))
		if i == 0 || i == 2 {
			line, err := prose.NewDocument(sent.Text)
			if err != nil {
				log.Fatal(err)
			}

			for _, tok := range line.Tokens() {
				if nouns[tok.Tag] {
					//	fmt.Println(tok.Text, tok.Tag)
				}
			}
		}
	}

	fmt.Println("Let's try to extract phrases!")
	for i, sent := range doc.Sentences() {
		fmt.Println(fmt.Sprintf("Analyzing sentence #%d: %s", i+1, sent))
		for j, phrase := range extract(sent.Text) {
			fmt.Println(fmt.Sprintf("Phrase #%d: %s", j+1, phrase))
		}
	}

	/*
		for i, ent := range doc.Entities() {
			fmt.Println(fmt.Sprintf("#%d. Entity: %s. Label: %s", i, ent.Text, ent.Label))
		}*/
}

/*
	 def extract(self, text):
        '''Return a list of noun phrases (strings) for body of text.'''
        sentences = nltk.tokenize.sent_tokenize(text)
        noun_phrases = []
        for sentence in sentences:
            parsed = self._parse_sentence(sentence)
            # Get the string representation of each subtree that is a
            # noun phrase tree
            phrases = [_normalize_tags(filter_insignificant(each,
                       self.INSIGNIFICANT_SUFFIXES)) for each in parsed
                       if isinstance(each, nltk.tree.Tree) and each.label()
                       == 'NP' and len(filter_insignificant(each)) >= 1
                       and _is_match(each, cfg=self.CFG)]
            nps = [tree2str(phrase) for phrase in phrases]
            noun_phrases.extend(nps)
        return noun_phrases
*/

var defineName = map[string]string{
	"NNP+NNP": "NNP",
	"NN+NN":   "NNI",
	"NNI+NN":  "NNI",
	"JJ+JJ":   "JJ",
	"JJ+NN":   "NNI",
}

// extract returns a list of noun phrases (strings) from the sentence.
func extract(sentence string) []string {
	parsed, err := prose.NewDocument(sentence)
	if err != nil {
		return []string{} //, err
	}

	normalizedTokens := normalizeTokens(parsed.Tokens())
	fmt.Println(normalizedTokens)

	merge := true
	for merge == true {
		merge = false
		for i := 0; i < len(normalizedTokens)-1; i++ {
			t1, t2 := normalizedTokens[i], normalizedTokens[i+1]
			key := fmt.Sprintf("%s+%s", t1.Tag, t2.Tag)
			value := defineName[key]
			if value != "" {
				merge = true
				mergedTok := prose.Token{value, fmt.Sprintf("%s %s", t1.Text, t2.Text), ""}
				//fmt.Println(fmt.Sprintf("MergedTok: %s", mergedTok.Text))
				normalizedTokens[i] = mergedTok
				normalizedTokens = append(normalizedTokens[:i+1], normalizedTokens[i+2:]...)
				break
			}
		}
	}

	nounPhrases := []string{}
	for _, tok := range normalizedTokens {
		if tok.Tag == "NNI" {
			nounPhrases = append(nounPhrases, tok.Text)
		}
	}

	return nounPhrases //, nil
}

var insignificantSuffixes = []string{"DT", "CC", "PRP$", "PRP"}

// filterInsignificantSuffixes filters out insignificant tokens from a chunk of tokens.
func filterInsignificantSuffixes(chunk []prose.Token) []prose.Token {
	filteredTokens := []prose.Token{}
	for _, tok := range chunk {
		significant := true
		for _, tag := range insignificantSuffixes {
			if strings.HasSuffix(tok.Tag, tag) {
				significant = false
				break
			}
		}

		if significant {
			filteredTokens = append(filteredTokens, tok)
		}
	}

	return filteredTokens
}

// normalizeTokens normalizes the corpus tags: ("NN", "NN-PL", "NNS") -> "NN"
func normalizeTokens(chunk []prose.Token) []prose.Token {
	normalizedTokens := []prose.Token{}
	for _, tok := range chunk {
		if tok.Tag == "NP-TL" || tok.Tag == "NP" {
			tok.Tag = "NNP"
			normalizedTokens = append(normalizedTokens, tok)
			continue
		}

		if strings.HasSuffix(tok.Tag, "-TL") {
			tok.Tag = strings.TrimSuffix(tok.Tag, "-TL")
			normalizedTokens = append(normalizedTokens, tok)
			continue
		}

		if strings.HasSuffix(tok.Tag, "S") {
			tok.Tag = strings.TrimSuffix(tok.Tag, "S")
			normalizedTokens = append(normalizedTokens, tok)
			continue
		}

		normalizedTokens = append(normalizedTokens, tok)
	}

	return normalizedTokens
}
