package dominio

// RepositorioAlumnos expone las operaciones que necesita el módulo de alumnos.
type RepositorioAlumnos interface {
	ConsultaAcademica
	GuardarAlumno(Alumno) (Alumno, error)
	EliminarAlumno(int) error
}
