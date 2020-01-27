package crawler

import (
	"fmt"
)

type DocsWriter interface {
	WriteDocs([]News) (n int, err error)
}

type printerWriter struct{}

func (*printerWriter) WriteDocs(docs []News) (n int, err error) {
	fmt.Println(docs)
	return len(docs), nil
}

func NewPrinterWriter() *printerWriter {
	return &printerWriter{}
}
