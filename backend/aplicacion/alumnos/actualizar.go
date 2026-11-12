package alumnos

import "seminario4/backend/dominio"

type actualizarAlumno struct {
	repositorio dominio.RepositorioAlumnos
}

func (c actualizarAlumno) ejecutar(id int, nombre string, edad int) (dominio.Alumno, error) {
	alumno, err := dominio.NuevoAlumno(nombre, edad)
	if err != nil {
		return dominio.Alumno{}, err
	}
	alumno.ID = id
	return c.repositorio.GuardarAlumno(alumno)
}
