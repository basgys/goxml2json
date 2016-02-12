package xml2json

import (
	"io"
	"strings"
)

// An Encoder writes JSON objects to an output stream.
type Encoder struct {
	w   io.Writer
	err error
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode writes the JSON encoding of v to the stream
func (enc *Encoder) Encode(root *Node) error {
	if enc.err != nil {
		return enc.err
	}
	if root == nil {
		return nil
	}

	enc.err = enc.format(root, 0)

	return enc.err
}

func (enc *Encoder) format(n *Node, lvl int) error {
	if n.IsComplex() {
		enc.write("{")

		i := 0
		tot := len(n.Children)
		for label, children := range n.Children {
			enc.write("\"")
			enc.write(label)
			enc.write("\": ")

			if len(children) > 1 {
				// Array
				enc.write("[")
				for j, c := range children {
					enc.format(c, lvl+1)

					if j < len(children)-1 {
						enc.write(", ")
					}
				}
				enc.write("]")
			} else {
				// Map
				enc.format(children[0], lvl+1)
			}

			if i < tot-1 {
				enc.write(", ")
			}
			i++
		}

		enc.write("}")
	} else {
		// TODO : Extract data type
		enc.write("\"")
		enc.write(escapeJSONString(n.Data))
		enc.write("\"")
	}

	return nil
}

func (enc *Encoder) write(s string) {
	enc.w.Write([]byte(s))
}

func escapeJSONString(s string) string {
	return strings.Replace(s, "\"", "\\\"", -1)
}
