package xml2json

import (
	"bytes"
	"unicode/utf8"
)

type (
	// an encodePlugin is added to an encoder to allow custom functionality at runtime
	encoderPlugin interface {
		AddTo(*Encoder) *Encoder
	}
	// a sanitizer overides the default string sanitization for encoding json
	encoderSanitizer interface {
		Sanitize(string) string
	}
	// CustomSanitizer santizes JSON using a best guess approach, used for converting all data to appropriate types
	customSanitizer struct {
		parseTypes []JSType
	}
)

// NewCustomSanitizer allows customized parsing behavior by passing in the desired JSTypes
func NewCustomSanitizer(ts ...JSType) *customSanitizer {
	return &customSanitizer{parseTypes: ts}
}

func (cs *customSanitizer) parseAsString(t JSType) bool {
	if t == String {
		return true
	}
	for i := 0; i < len(cs.parseTypes); i++ {
		if cs.parseTypes[i] == t {
			return false
		}
	}
	return true
}

func (cs *customSanitizer) AddTo(e *Encoder) *Encoder {
	e.s = cs
	return e
}

func (cs *customSanitizer) Sanitize(s string) string {
	var buf bytes.Buffer
	// prefix output according to santizer settings
	jsType := Str2JSType(s)
	if cs.parseAsString(jsType) {
		buf.WriteByte('"')
	}

	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' && b != '&' {
				i++
				continue
			}
			if start < i {
				buf.WriteString(s[start:i])
			}
			switch b {
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			case '\t':
				buf.WriteByte('\\')
				buf.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \n and \r,
				// as well as <, > and &. The latter are escaped because they
				// can lead to security holes when user-controlled strings
				// are rendered into JSON and served to some browsers.
				buf.WriteString(`\u00`)
				buf.WriteByte(hex[b>>4])
				buf.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\u202`)
			buf.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		buf.WriteString(s[start:])
	}
	// suffix output according to sanitizer settings
	if cs.parseAsString(jsType) {
		buf.WriteByte('"')
	}
	return buf.String()
}
