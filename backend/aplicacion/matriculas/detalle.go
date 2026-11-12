package matriculas

import "seminario4/backend/dominio"

// DetalleMatricula es un DTO: la indicación se calcula, no se persiste.
type DetalleMatricula struct {
	dominio.Matricula
	Indicacion string `json:"indicacion"`
}
