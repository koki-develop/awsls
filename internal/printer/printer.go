package printer

import (
	"fmt"
	"io"

	"github.com/koki-develop/awsls/internal/aws"
)

type Printer interface {
	Print(w io.Writer, rs aws.Resources) error
}

func New(format string) (Printer, error) {
	switch format {
	case "json":
		return newJsonPrinter(), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
