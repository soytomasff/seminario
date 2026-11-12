package infraestructura

import "seminario4/backend/dominio"

func (r *RepositorioArchivo) GuardarAlumno(alumno dominio.Alumno) (dominio.Alumno, error) {
	err := r.modificar(func(datos *dominio.DatosAcademicos) error {
		if alumno.ID == 0 {
			alumno.ID = datos.SiguienteAlumno
			datos.SiguienteAlumno++
			datos.Alumnos = append(datos.Alumnos, alumno)
			return nil
		}
		for indice, actual := range datos.Alumnos {
			if actual.ID == alumno.ID {
				datos.Alumnos[indice] = alumno
				return nil
			}
		}
		return dominio.ErrNoEncontrado
	})
	return alumno, err
}

func (r *RepositorioArchivo) EliminarAlumno(id int) error {
	return r.modificar(func(datos *dominio.DatosAcademicos) error {
		indiceAlumno := -1
		for indice, alumno := range datos.Alumnos {
			if alumno.ID == id {
				indiceAlumno = indice
				break
			}
		}
		if indiceAlumno == -1 {
			return dominio.ErrNoEncontrado
		}
		datos.Alumnos = append(datos.Alumnos[:indiceAlumno], datos.Alumnos[indiceAlumno+1:]...)
		// Borrado en cascada: al eliminar un alumno quitamos también sus matrículas.
		matriculasRestantes := []dominio.Matricula{}
		for _, matricula := range datos.Matriculas {
			if matricula.AlumnoID != id {
				matriculasRestantes = append(matriculasRestantes, matricula)
			}
		}
		datos.Matriculas = matriculasRestantes
		return nil
	})
}
