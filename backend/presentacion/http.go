// Package presentacion adapta los casos de uso a HTTP.
package presentacion

import (
	"net/http"

	"seminario4/backend/aplicacion"
)

// NuevoManejador conecta las rutas con el controlador de cada módulo.
func NuevoManejador(servicio *aplicacion.ServicioAcademico, origen string) http.Handler {
	alumnos := &ControladorAlumnos{fachada: servicio.Alumnos}
	matriculas := &ControladorMatriculas{gestor: servicio.Matriculas}
	sistema := &ControladorSistema{servicio: servicio}
	rutas := http.NewServeMux()
	rutas.HandleFunc("GET /healthz", sistema.salud)
	rutas.HandleFunc("GET /api/resumen", sistema.consultar)
	rutas.HandleFunc("GET /api/alumnos", alumnos.consultar)
	rutas.HandleFunc("POST /api/alumnos", alumnos.guardar)
	rutas.HandleFunc("PUT /api/alumnos/{id}", alumnos.guardar)
	rutas.HandleFunc("DELETE /api/alumnos/{id}", alumnos.eliminar)
	rutas.HandleFunc("GET /api/matriculas", matriculas.consultar)
	rutas.HandleFunc("POST /api/matriculas", matriculas.guardar)
	rutas.HandleFunc("PUT /api/matriculas/{id}", matriculas.guardar)
	rutas.HandleFunc("DELETE /api/matriculas/{id}", matriculas.eliminar)
	rutas.HandleFunc("POST /api/ejemplos", sistema.reiniciar)
	return conCORS(rutas, origen)
}
