package render

import (
	"encoding/json"
	"net/http"
)

type JSON struct {
	Data any
}

var jsonContentType = []string{"application/json; charset=utf-8"}

func (r JSON) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return json.NewEncoder(w).Encode(r.Data)
}

func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}
