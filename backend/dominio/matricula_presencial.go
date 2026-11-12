package dominio

// MatriculaPresencial compone los datos de la matrícula y participa en Visitor.
type MatriculaPresencial struct {
	Datos Matricula
}

func (m MatriculaPresencial) Aceptar(visitante VisitanteMatricula) string {
	return visitante.VisitarPresencial(m)
}

var _ ElementoMatricula = MatriculaPresencial{}
