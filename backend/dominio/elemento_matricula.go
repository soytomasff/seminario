package dominio

// ElementoMatricula es el contrato de los elementos que aceptan un visitante.
type ElementoMatricula interface {
	Aceptar(VisitanteMatricula) string
}
