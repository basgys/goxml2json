package xml2json

import (
	"strings"
)

type (
	// an encodePlugin is added to an encoder to allow custom functionality at runtime
	encoderPlugin interface {
		AddTo(*Encoder) *Encoder
	}
	// a type converter overides the default string sanitization for encoding json
	encoderTypeConverter interface {
		Convert(string) string
	}
	// customTypeConverter converts strings to JSON types using a best guess approach, only parses the JSON types given
	// when initialized via WithTypeConverter
	customTypeConverter struct {
		parseTypes []JSType
	}

	attrPrefixer    string
	contentPrefixer string
)

// WithTypeConverter allows customized js type conversion behavior by passing in the desired JSTypes
func WithTypeConverter(ts ...JSType) *customTypeConverter {
	return &customTypeConverter{parseTypes: ts}
}

func (tc *customTypeConverter) parseAsString(t JSType) bool {
	if t == String {
		return true
	}
	for i := 0; i < len(tc.parseTypes); i++ {
		if tc.parseTypes[i] == t {
			return false
		}
	}
	return true
}

// Adds the type converter to the encoder
func (tc *customTypeConverter) AddTo(e *Encoder) *Encoder {
	e.tc = tc
	return e
}

func (tc *customTypeConverter) Convert(s string) string {
	// remove quotes if they exists
	if strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
		s = s[1 : len(s)-1]
	}
	jsType := Str2JSType(s)
	if tc.parseAsString(jsType) {
		// add the quotes removed at the start of this func
		s = `"` + s + `"`
	}
	return s
}

// WithAttrPrefix appends the given prefix to the json output of xml attribute fields to preserve namespaces
func WithAttrPrefix(prefix string) *attrPrefixer {
	ap := attrPrefixer(prefix)
	return &ap
}

func (a *attrPrefixer) AddTo(e *Encoder) {
	e.attributePrefix = string((*a))
}

// WithContentPrefix appends the given prefix to the json output of xml content fields to preserve namespaces
func WithContentPrefix(prefix string) *contentPrefixer {
	c := contentPrefixer(prefix)
	return &c
}

func (c *contentPrefixer) AddTo(e *Encoder) {
	e.contentPrefix = string((*c))
}
