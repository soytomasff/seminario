package presentacion

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func leerID(w http.ResponseWriter, r *http.Request) (int, bool) {
	if r.Method == http.MethodPost {
		return 0, true
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		responder(w, http.StatusBadRequest, map[string]string{"error": "El ID debe ser un número entero positivo."})
		return 0, false
	}
	return id, true
}

func leerJSON(w http.ResponseWriter, r *http.Request, destino any) bool {
	lector := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
	lector.DisallowUnknownFields()
	if err := lector.Decode(destino); err != nil {
		responder(w, http.StatusBadRequest, map[string]string{"error": "Envía un JSON válido con los campos y tipos esperados."})
		return false
	}
	if err := lector.Decode(&struct{}{}); err != io.EOF {
		responder(w, http.StatusBadRequest, map[string]string{"error": "Envía un solo objeto JSON."})
		return false
	}
	return true
}
