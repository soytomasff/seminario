package dominio

// CrearElementoMatricula selecciona el elemento; el visitante define la operación.
func CrearElementoMatricula(datos Matricula) ElementoMatricula {
	if datos.Modalidad == "Virtual" {
		return MatriculaVirtual{Datos: datos}
	}
	return MatriculaPresencial{Datos: datos}
}
