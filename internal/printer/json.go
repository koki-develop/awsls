package printer

import (
	"encoding/json"
	"io"

	"github.com/koki-develop/awsls/internal/aws"
)

var _ Printer = (*jsonPrinter)(nil)

type jsonPrinter struct{}

func newJsonPrinter() *jsonPrinter {
	return &jsonPrinter{}
}

func (j *jsonPrinter) Print(w io.Writer, rs aws.Resources) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rs)
}
