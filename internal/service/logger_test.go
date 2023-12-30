package service

import (
	"bytes"
	"testing"
)

func TestLogger(t *testing.T) {
	output := &bytes.Buffer{}
	l := Logger{Writer: output}
	in := "Hello world!\n"
	l.Log(in)
	if out := output.String(); out != in {
		t.Fatalf("logged %q instead of %q", out, in)
	}
}
