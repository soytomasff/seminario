package matriculas

import "seminario4/backend/dominio"

// Gestor coordina el CRUD y recibe la operación Visitor mediante composición.
type Gestor struct {
	repositorio dominio.RepositorioMatriculas
	visitante   dominio.VisitanteMatricula
}

func NuevoGestor(repositorio dominio.RepositorioMatriculas, visitante dominio.VisitanteMatricula) *Gestor {
	return &Gestor{repositorio: repositorio, visitante: visitante}
}
