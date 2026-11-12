package matriculas

import "seminario4/backend/dominio"

// guardar comparte la validación y persistencia de registrar y actualizar.
func (g *Gestor) guardar(id, alumnoID int, curso, modalidad string) (DetalleMatricula, error) {
	matricula, err := dominio.NuevaMatricula(alumnoID, curso, modalidad)
	if err != nil {
		return DetalleMatricula{}, err
	}
	matricula.ID = id
	guardada, err := g.repositorio.GuardarMatricula(matricula)
	if err != nil {
		return DetalleMatricula{}, err
	}
	return g.describirMatricula(guardada), nil
}
