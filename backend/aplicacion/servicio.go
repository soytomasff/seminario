// Aplicación organiza los casos de uso: registrar, consultar, editar y eliminar.
package aplicacion

import (
	"seminario4/backend/aplicacion/alumnos"
	"seminario4/backend/aplicacion/matriculas"
	"seminario4/backend/dominio"
)

// ServicioAcademico compone ambos módulos y coordina consultas del sistema.
type ServicioAcademico struct {
	repositorio dominio.RepositorioAcademico
	Alumnos     *alumnos.Fachada
	Matriculas  *matriculas.Gestor
}

// Esto es inyección de dependencias manual: recibimos el repositorio desde main.go.
func NuevoServicio(repositorio dominio.RepositorioAcademico) *ServicioAcademico {
	return &ServicioAcademico{
		repositorio: repositorio,
		Alumnos:     alumnos.NuevaFachada(repositorio),
		Matriculas:  matriculas.NuevoGestor(repositorio, matriculas.VisitanteIndicaciones{}),
	}
}

func (s *ServicioAcademico) Consultar() (dominio.DatosAcademicos, error) {
	return s.repositorio.Consultar()
}

func (s *ServicioAcademico) Reiniciar() error {
	return s.repositorio.Reiniciar()
}
