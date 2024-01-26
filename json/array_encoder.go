package json

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
)

var (
	comma    = []byte(",\n")
	start    = []byte("[\n")
	end      = []byte("\n]\n")
	noindent = []string{}
)

type ArrayEncoder[T any] struct {
	w        io.Writer
	prefixer func() error
	m        sync.Mutex
	indent   []string
}

func NewArrayEncoder[T any](w io.Writer) *ArrayEncoder[T] {
	ae := &ArrayEncoder[T]{
		w:      w,
		indent: noindent,
	}
	ae.Reset()
	return ae
}

func (ae *ArrayEncoder[T]) SetIndent(prefix, indent string) {
	if prefix == "" && indent == "" {
		ae.indent = noindent
		return
	}
	ae.indent = []string{prefix, indent}
}

func (ae *ArrayEncoder[T]) hasIndent() bool {
	if len(ae.indent) < 2 || ((ae.indent[0] == "") && (ae.indent[1] == "")) {
		return false
	}
	return true
}

//TODO: Better indent implementation. when indenting,
//all of the input except for the [] should be indented by one notch

func (ae *ArrayEncoder[T]) Encode(v T) error {
	ae.m.Lock()
	defer ae.m.Unlock()
	err := ae.prefixer()
	if err != nil {
		return err
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if !ae.hasIndent() || len(b) <= 2 {
		_, err = ae.w.Write(b)
		return err
	}

	var buf bytes.Buffer
	err = json.Indent(&buf, b, ae.indent[0], ae.indent[1])
	if err != nil {
		return err
	}
	_, err = buf.WriteTo(ae.w)
	return err
}

func (ae *ArrayEncoder[T]) Finish() error {
	var buf bytes.Buffer
	buf.Write(end)
	_, err := ae.w.Write(end)
	return err
}

func (ae *ArrayEncoder[T]) Reset() {
	ae.prefixer = func() error {
		ae.prefixer = func() error {
			_, err := ae.w.Write(comma)
			return err
		}
		_, err := ae.w.Write(start)
		return err
	}
}
