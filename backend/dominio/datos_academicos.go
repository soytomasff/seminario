package dominio

type DatosAcademicos struct {
	Alumnos            []Alumno    `json:"alumnos"`
	Matriculas         []Matricula `json:"matriculas"`
	SiguienteAlumno    int         `json:"siguiente_alumno"`
	SiguienteMatricula int         `json:"siguiente_matricula"`
}
