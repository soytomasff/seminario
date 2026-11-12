package matriculas

func (g *Gestor) Consultar() ([]DetalleMatricula, error) {
	datos, err := g.repositorio.Consultar()
	if err != nil {
		return nil, err
	}
	return g.Describir(datos.Matriculas), nil
}
