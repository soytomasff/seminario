package matriculas

func (g *Gestor) Eliminar(id int) error {
	return g.repositorio.EliminarMatricula(id)
}
