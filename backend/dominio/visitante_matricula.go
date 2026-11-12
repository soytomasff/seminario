package dominio

// VisitanteMatricula define una operación para cada modalidad concreta.
type VisitanteMatricula interface {
	VisitarPresencial(MatriculaPresencial) string
	VisitarVirtual(MatriculaVirtual) string
}
