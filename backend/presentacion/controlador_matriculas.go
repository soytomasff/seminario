package presentacion

import (
	"net/http"

	"seminario4/backend/aplicacion/matriculas"
)

type ControladorMatriculas struct {
	gestor *matriculas.Gestor
}

type entradaMatricula struct {
	AlumnoID  int    `json:"alumno_id"`
	Curso     string `json:"curso"`
	Modalidad string `json:"modalidad"`
}

func (c *ControladorMatriculas) consultar(w http.ResponseWriter, r *http.Request) {
	lista, err := c.gestor.Consultar()
	if err != nil {
		responderError(w, err)
		return
	}
	responder(w, http.StatusOK, lista)
}

func (c *ControladorMatriculas) guardar(w http.ResponseWriter, r *http.Request) {
	id, valido := leerID(w, r)
	if !valido {
		return
	}
	var entrada entradaMatricula
	if !leerJSON(w, r, &entrada) {
		return
	}
	var detalle matriculas.DetalleMatricula
	var err error
	estado := http.StatusOK
	if r.Method == http.MethodPost {
		detalle, err = c.gestor.Registrar(entrada.AlumnoID, entrada.Curso, entrada.Modalidad)
		estado = http.StatusCreated
	} else {
		detalle, err = c.gestor.Actualizar(id, entrada.AlumnoID, entrada.Curso, entrada.Modalidad)
	}
	if err != nil {
		responderError(w, err)
		return
	}
	responder(w, estado, detalle)
}

func (c *ControladorMatriculas) eliminar(w http.ResponseWriter, r *http.Request) {
	id, valido := leerID(w, r)
	if !valido {
		return
	}
	if err := c.gestor.Eliminar(id); err != nil {
		responderError(w, err)
		return
	}
	responder(w, http.StatusNoContent, nil)
}
