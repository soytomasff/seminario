package infraestructura

import (
	"fmt"
	"net/http"
	"strconv"

	"seminario4/backend/dominio"
)

func (r *RepositorioSupabase) GuardarMatricula(matricula dominio.Matricula) (dominio.Matricula, error) {
	metodo, ruta := http.MethodPost, "/matriculas"
	if matricula.ID != 0 {
		metodo, ruta = http.MethodPatch, "/matriculas?id=eq."+strconv.Itoa(matricula.ID)
	}
	datos := struct {
		AlumnoID  int    `json:"alumno_id"`
		Curso     string `json:"curso"`
		Modalidad string `json:"modalidad"`
	}{matricula.AlumnoID, matricula.Curso, matricula.Modalidad}
	var guardadas []dominio.Matricula
	err := r.solicitar(metodo, ruta, datos, &guardadas)
	if err != nil {
		return dominio.Matricula{}, err
	}
	if len(guardadas) == 0 && matricula.ID != 0 {
		return dominio.Matricula{}, dominio.ErrNoEncontrado
	}
	if len(guardadas) != 1 || guardadas[0].ID < 1 {
		return dominio.Matricula{}, fmt.Errorf("%w: Supabase no devolvió la matrícula guardada.", dominio.ErrAlmacenamiento)
	}
	return guardadas[0], nil
}

func (r *RepositorioSupabase) EliminarMatricula(id int) error {
	return r.eliminar("matriculas", id)
}
