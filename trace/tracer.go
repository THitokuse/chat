package trace

import (
  "io"
)
//Tracerはコード内での出来事を記録できるオブジェクトを表すインターフェースです。
type Tracer interface {
  Trace(...interface{})
}
type tracer struct {
  out io.Writer
}
func (t *tracer) Trace(a...interface{}) {}

func New(w io.Writer) Tracer {
  return &tracer{out: w}
}
