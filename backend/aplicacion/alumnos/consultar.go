package alumnos

import "seminario4/backend/dominio"

type consultarAlumnos struct {
	repositorio dominio.RepositorioAlumnos
}

func (c consultarAlumnos) ejecutar() ([]dominio.Alumno, error) {
	datos, err := c.repositorio.Consultar()
	return datos.Alumnos, err
}
