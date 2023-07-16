package printer

import (
	"encoding/csv"
	"io"

	"github.com/koki-develop/awsls/internal/aws"
)

var _ Printer = (*csvPrinter)(nil)

type csvPrinter struct{}

func newCsvPrinter() *csvPrinter {
	return &csvPrinter{}
}

func (*csvPrinter) Print(w io.Writer, rs aws.Resources) error {
	cw := csv.NewWriter(w)

	if err := cw.Write([]string{"service", "region", "resource", "arn"}); err != nil {
		return err
	}

	for _, r := range rs {
		if err := cw.Write([]string{r.Service, r.Region, r.Resource, r.ARN}); err != nil {
			return err
		}
	}

	return nil
}
