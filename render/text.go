package render

import (
	"fmt"
	"net/http"
)

type String struct {
	Format string
	Data   []any
}

var textContentType = []string{"text/plain; charset=utf-8"}

func (r String) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	_, err := fmt.Fprintf(w, r.Format, r.Data...)
	return err
}

func (r String) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, textContentType)
}
