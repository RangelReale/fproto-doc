package fproto_doc

import (
	"io"

	"github.com/RangelReale/fdep"
)

type Generator interface {
	Generate(fdep *fdep.Dep, w io.Writer) error
}
