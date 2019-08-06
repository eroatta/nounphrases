package finder

import (
	"testing"

	//"github.com/eroatta/nounphrases/finder"

	"github.com/stretchr/testify/assert"
)

func TestFind_ShouldReturnNounPhrasesFromText(t *testing.T) {
	text := `
		Marshal returns the JSON encoding of v.
		Marshal traverses the value v recursively. If an encountered value implements the Marshaler interface and is not a nil pointer, Marshal calls its MarshalJSON method to produce JSON. If no MarshalJSON method is present but the value implements encoding.TextMarshaler instead, Marshal calls its MarshalText method and encodes the result as a JSON string. The nil pointer exception is not strictly necessary but mimics a similar, necessary exception in the behavior of UnmarshalJSON.
		Otherwise, Marshal uses the following type-dependent default encodings:
			* Boolean values encode as JSON booleans.
			* Floating point, integer, and Number values encode as JSON numbers.
			* String values encode as JSON strings coerced to valid UTF-8, replacing invalid bytes with the Unicode replacement rune. The angle brackets "<" and ">" are escaped to "\u003c" and "\u003e" to keep some browsers from misinterpreting JSON output as HTML. Ampersand "&" is also escaped to "\u0026" for the same reason. This escaping can be disabled using an Encoder that had SetEscapeHTML(false) called on it.
			* Array and slice values encode as JSON arrays, except that []byte encodes as a base64-encoded string, and a nil slice encodes as the null JSON value.
			* Struct values encode as JSON objects. Each exported struct field becomes a member of the object, using the field name as the object key, unless the field is omitted for one of the reasons given below.
		The encoding of each struct field can be customized by the format string stored under the "json" key in the struct field's tag. The format string gives the name of the field, possibly followed by a comma-separated list of options. The name may be empty in order to specify options without overriding the default field name.
		The "omitempty" option specifies that the field should be omitted from the encoding if the field has an empty value, defined as false, 0, a nil pointer, a nil interface value, and any empty array, slice, map, or string.
		As a special case, if the field tag is "-", the field is always omitted. Note that a field with name "-" can still be generated using the tag "-,". 
	`
	got, err := Find(text)

	assert.NoError(t, err, "no error should be raised")
	assert.Equal(t, 29, len(got))

	expected := []string{
		"value v",
		"encountered value",
		"nil pointer",
		"nil pointer exception",
		"necessary exception",
		"following type-dependent default encodings",
		"invalid bytes",
		"replacement rune",
		"angle brackets",
		"same reason",
		"slice values",
		"base64-encoded string",
		"nil slice encodes",
		"struct field",
		"field name",
		"object key",
		"struct field",
		"format string",
		"json key",
		"struct field",
		"comma-separated list",
		"default field name",
		"option specifies",
		"empty value",
		"nil pointer",
		"nil interface value",
		"empty array",
		"special case",
		"field tag",
	}
	assert.ElementsMatch(t, expected, got)
}
