package dominio

// ConsultaAcademica devuelve una fotografía consistente de los dos módulos.
type ConsultaAcademica interface {
	Consultar() (DatosAcademicos, error)
}

// RepositorioAcademico compone los contratos utilizados por cada CRUD.
// Las implementaciones se construyen una sola vez y comparten la persistencia.
type RepositorioAcademico interface {
	RepositorioAlumnos
	RepositorioMatriculas
	Reiniciar() error
}
