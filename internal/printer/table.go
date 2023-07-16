package printer

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/koki-develop/awsls/internal/aws"
)

var _ Printer = (*tablePrinter)(nil)

type tablePrinter struct{}

func newTablePrinter() *tablePrinter {
	return &tablePrinter{}
}

func (*tablePrinter) Print(w io.Writer, rs aws.Resources) error {
	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
	defer tw.Flush()

	if _, err := fmt.Fprintln(tw, "Service\tRegion\tResource\tARN"); err != nil {
		return err
	}

	for _, r := range rs {
		if _, err := fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", r.Service, r.Region, r.Resource, r.ARN); err != nil {
			return err
		}
	}

	return nil
}
