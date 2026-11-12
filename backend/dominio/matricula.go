package dominio

import "errors"

// Matricula vincula un alumno con un curso.
type Matricula struct {
	ID        int    `json:"id"`
	AlumnoID  int    `json:"alumno_id"`
	Curso     string `json:"curso"`
	Modalidad string `json:"modalidad"`
}

func CursosDisponibles() []string {
	return []string{"Programación I", "Desarrollo Web", "Base de Datos", "Diseño de Interfaces"}
}

func NuevaMatricula(alumnoID int, curso, modalidad string) (Matricula, error) {
	if alumnoID < 1 {
		return Matricula{}, errors.New("Debes elegir un alumno registrado.")
	}
	// Los datos anteriores a la modalidad se interpretan como presenciales.
	if modalidad == "" {
		modalidad = "Presencial"
	}
	if modalidad != "Presencial" && modalidad != "Virtual" {
		return Matricula{}, errors.New("La modalidad debe ser Presencial o Virtual.")
	}
	for _, disponible := range CursosDisponibles() {
		if curso == disponible {
			return Matricula{AlumnoID: alumnoID, Curso: curso, Modalidad: modalidad}, nil
		}
	}
	return Matricula{}, errors.New("Selecciona uno de los cursos disponibles.")
}
