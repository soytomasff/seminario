package matriculas

func (g *Gestor) Actualizar(id, alumnoID int, curso, modalidad string) (DetalleMatricula, error) {
	return g.guardar(id, alumnoID, curso, modalidad)
}
