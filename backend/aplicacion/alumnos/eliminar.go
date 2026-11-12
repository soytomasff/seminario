package alumnos

import "seminario4/backend/dominio"

type eliminarAlumno struct {
	repositorio dominio.RepositorioAlumnos
}

func (c eliminarAlumno) ejecutar(id int) error {
	return c.repositorio.EliminarAlumno(id)
}
