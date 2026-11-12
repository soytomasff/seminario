package dominio

type RepositorioMatriculas interface {
	ConsultaAcademica
	GuardarMatricula(Matricula) (Matricula, error)
	EliminarMatricula(int) error
}
