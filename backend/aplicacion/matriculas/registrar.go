package matriculas

func (g *Gestor) Registrar(alumnoID int, curso, modalidad string) (DetalleMatricula, error) {
	return g.guardar(0, alumnoID, curso, modalidad)
}
