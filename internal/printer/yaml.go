package printer

import (
	"io"

	"github.com/koki-develop/awsls/internal/aws"
	"gopkg.in/yaml.v3"
)

var _ Printer = (*yamlPrinter)(nil)

type yamlPrinter struct{}

func newYamlPrinter() *yamlPrinter {
	return &yamlPrinter{}
}

func (y *yamlPrinter) Print(w io.Writer, rs aws.Resources) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	return enc.Encode(rs)
}
