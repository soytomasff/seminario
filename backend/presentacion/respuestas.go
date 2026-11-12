package presentacion

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"seminario4/backend/dominio"
)

func responder(w http.ResponseWriter, estado int, contenido any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(estado)
	if contenido != nil {
		if err := json.NewEncoder(w).Encode(contenido); err != nil {
			log.Print(err)
		}
	}
}

func responderError(w http.ResponseWriter, err error) {
	estado := http.StatusBadRequest
	if errors.Is(err, dominio.ErrNoEncontrado) {
		estado = http.StatusNotFound
	}
	if errors.Is(err, dominio.ErrDuplicado) {
		estado = http.StatusConflict
	}
	if errors.Is(err, dominio.ErrAlmacenamiento) {
		log.Print(err)
		responder(w, http.StatusInternalServerError, map[string]string{"error": dominio.ErrAlmacenamiento.Error()})
		return
	}
	responder(w, estado, map[string]string{"error": err.Error()})
}
