package presentacion

import (
	"net/http"

	"seminario4/backend/aplicacion"
	"seminario4/backend/dominio"
)

type ControladorSistema struct {
	servicio *aplicacion.ServicioAcademico
}

func (c *ControladorSistema) salud(w http.ResponseWriter, r *http.Request) {
	responder(w, http.StatusOK, map[string]string{"estado": "ok"})
}

func (c *ControladorSistema) consultar(w http.ResponseWriter, r *http.Request) {
	datos, err := c.servicio.Consultar()
	if err != nil {
		responderError(w, err)
		return
	}
	responder(w, http.StatusOK, map[string]any{
		"alumnos": datos.Alumnos, "matriculas": c.servicio.Matriculas.Describir(datos.Matriculas), "cursos": dominio.CursosDisponibles(),
	})
}

func (c *ControladorSistema) reiniciar(w http.ResponseWriter, r *http.Request) {
	if err := c.servicio.Reiniciar(); err != nil {
		responderError(w, err)
		return
	}
	responder(w, http.StatusNoContent, nil)
}
