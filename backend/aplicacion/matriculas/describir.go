package matriculas

import "seminario4/backend/dominio"

func (g *Gestor) Describir(matriculas []dominio.Matricula) []DetalleMatricula {
	detalles := make([]DetalleMatricula, 0, len(matriculas))
	for _, matricula := range matriculas {
		detalles = append(detalles, g.describirMatricula(matricula))
	}
	return detalles
}

func (g *Gestor) describirMatricula(matricula dominio.Matricula) DetalleMatricula {
	if matricula.Modalidad == "" {
		matricula.Modalidad = "Presencial"
	}
	elemento := dominio.CrearElementoMatricula(matricula)
	return DetalleMatricula{Matricula: matricula, Indicacion: elemento.Aceptar(g.visitante)}
}
