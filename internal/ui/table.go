package ui

import (
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
)

type TableOptions struct {
	Headers     []string
	ColumnAlign []int
	NoHeader    bool
	Border      bool
}

func NewTable(opts TableOptions) *tablewriter.Table {
	return NewTableWriter(os.Stdout, opts)
}

func NewTableWriter(w io.Writer, opts TableOptions) *tablewriter.Table {
	table := tablewriter.NewWriter(w)

	if len(opts.Headers) > 0 && !opts.NoHeader {
		// Convert []string to []any for the variadic Header method
		headers := make([]any, len(opts.Headers))
		for i, h := range opts.Headers {
			headers[i] = h
		}
		table.Header(headers...)
	}

	return table
}

func RenderTable(headers []string, rows [][]string) {
	table := NewTable(TableOptions{
		Headers: headers,
	})
	for _, row := range rows {
		table.Append(row)
	}
	table.Render()
}

func RenderKeyValue(pairs [][]string) {
	table := NewTable(TableOptions{
		NoHeader: true,
	})
	for _, pair := range pairs {
		table.Append(pair)
	}
	table.Render()
}
