package infraestructura

import "seminario4/backend/dominio"

func (r *RepositorioArchivo) GuardarMatricula(matricula dominio.Matricula) (dominio.Matricula, error) {
	err := r.modificar(func(datos *dominio.DatosAcademicos) error {
		indiceMatricula := -1
		for indice, actual := range datos.Matriculas {
			if actual.ID == matricula.ID {
				indiceMatricula = indice
			}
		}
		if matricula.ID != 0 && indiceMatricula == -1 {
			return dominio.ErrNoEncontrado
		}
		alumnoExiste := false
		for _, alumno := range datos.Alumnos {
			if alumno.ID == matricula.AlumnoID {
				alumnoExiste = true
			}
		}
		if !alumnoExiste {
			return dominio.ErrAlumnoNoExiste
		}
		for _, actual := range datos.Matriculas {
			if actual.ID != matricula.ID && actual.AlumnoID == matricula.AlumnoID && actual.Curso == matricula.Curso {
				return dominio.ErrDuplicado
			}
		}
		if matricula.ID == 0 {
			matricula.ID = datos.SiguienteMatricula
			datos.SiguienteMatricula++
			datos.Matriculas = append(datos.Matriculas, matricula)
		} else {
			datos.Matriculas[indiceMatricula] = matricula
		}
		return nil
	})
	return matricula, err
}

func (r *RepositorioArchivo) EliminarMatricula(id int) error {
	return r.modificar(func(datos *dominio.DatosAcademicos) error {
		for indice, matricula := range datos.Matriculas {
			if matricula.ID == id {
				datos.Matriculas = append(datos.Matriculas[:indice], datos.Matriculas[indice+1:]...)
				return nil
			}
		}
		return dominio.ErrNoEncontrado
	})
}
