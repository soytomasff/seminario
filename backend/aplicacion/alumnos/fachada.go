package alumnos

import "seminario4/backend/dominio"

// Fachada ofrece una entrada al CRUD y delega en cuatro casos de uso privados.
// Cada caso de uso tiene su implementación en un archivo independiente.
type Fachada struct {
	registro      registrarAlumno
	consulta      consultarAlumnos
	actualizacion actualizarAlumno
	eliminacion   eliminarAlumno
}

func NuevaFachada(repositorio dominio.RepositorioAlumnos) *Fachada {
	return &Fachada{
		registro:      registrarAlumno{repositorio: repositorio},
		consulta:      consultarAlumnos{repositorio: repositorio},
		actualizacion: actualizarAlumno{repositorio: repositorio},
		eliminacion:   eliminarAlumno{repositorio: repositorio},
	}
}

func (f *Fachada) Registrar(nombre string, edad int) (dominio.Alumno, error) {
	return f.registro.ejecutar(nombre, edad)
}

func (f *Fachada) Consultar() ([]dominio.Alumno, error) {
	return f.consulta.ejecutar()
}

func (f *Fachada) Actualizar(id int, nombre string, edad int) (dominio.Alumno, error) {
	return f.actualizacion.ejecutar(id, nombre, edad)
}

func (f *Fachada) Eliminar(id int) error {
	return f.eliminacion.ejecutar(id)
}
