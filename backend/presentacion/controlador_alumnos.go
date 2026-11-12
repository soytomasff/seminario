package presentacion

import (
	"net/http"

	"seminario4/backend/aplicacion/alumnos"
	"seminario4/backend/dominio"
)

type ControladorAlumnos struct {
	fachada *alumnos.Fachada
}

type entradaAlumno struct {
	Nombre string `json:"nombre"`
	Edad   int    `json:"edad"`
}

func (c *ControladorAlumnos) consultar(w http.ResponseWriter, r *http.Request) {
	lista, err := c.fachada.Consultar()
	if err != nil {
		responderError(w, err)
		return
	}
	responder(w, http.StatusOK, lista)
}

func (c *ControladorAlumnos) guardar(w http.ResponseWriter, r *http.Request) {
	id, valido := leerID(w, r)
	if !valido {
		return
	}
	var entrada entradaAlumno
	if !leerJSON(w, r, &entrada) {
		return
	}
	var alumno dominio.Alumno
	var err error
	estado := http.StatusOK
	if r.Method == http.MethodPost {
		alumno, err = c.fachada.Registrar(entrada.Nombre, entrada.Edad)
		estado = http.StatusCreated
	} else {
		alumno, err = c.fachada.Actualizar(id, entrada.Nombre, entrada.Edad)
	}
	if err != nil {
		responderError(w, err)
		return
	}
	responder(w, estado, alumno)
}

func (c *ControladorAlumnos) eliminar(w http.ResponseWriter, r *http.Request) {
	id, valido := leerID(w, r)
	if !valido {
		return
	}
	if err := c.fachada.Eliminar(id); err != nil {
		responderError(w, err)
		return
	}
	responder(w, http.StatusNoContent, nil)
}
