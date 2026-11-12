package alumnos

import "seminario4/backend/dominio"

type registrarAlumno struct {
	repositorio dominio.RepositorioAlumnos
}

func (c registrarAlumno) ejecutar(nombre string, edad int) (dominio.Alumno, error) {
	alumno, err := dominio.NuevoAlumno(nombre, edad)
	if err != nil {
		return dominio.Alumno{}, err
	}
	return c.repositorio.GuardarAlumno(alumno)
}
