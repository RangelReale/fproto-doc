package fproto_doc

import (
	"io"

	"github.com/RangelReale/fproto/fdep"
)

type Generator interface {
	Generate(fdep *fdep.Dep, w io.Writer) error
}
