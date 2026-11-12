package matriculas

import "seminario4/backend/dominio"

// VisitanteIndicaciones implementa una operación independiente de las entidades.
type VisitanteIndicaciones struct{}

func (v VisitanteIndicaciones) VisitarPresencial(m dominio.MatriculaPresencial) string {
	return "Asiste al instituto para el curso " + m.Datos.Curso + "."
}

func (v VisitanteIndicaciones) VisitarVirtual(m dominio.MatriculaVirtual) string {
	return "Ingresa al aula virtual para el curso " + m.Datos.Curso + "."
}

var _ dominio.VisitanteMatricula = VisitanteIndicaciones{}
