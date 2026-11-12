package dominio

type MatriculaVirtual struct {
	Datos Matricula
}

func (m MatriculaVirtual) Aceptar(visitante VisitanteMatricula) string {
	return visitante.VisitarVirtual(m)
}

var _ ElementoMatricula = MatriculaVirtual{}
