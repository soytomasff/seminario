package infraestructura

import (
	"fmt"
	"net/http"
	"strconv"

	"seminario4/backend/dominio"
)

func (r *RepositorioSupabase) GuardarAlumno(alumno dominio.Alumno) (dominio.Alumno, error) {
	metodo, ruta := http.MethodPost, "/alumnos"
	if alumno.ID != 0 {
		metodo, ruta = http.MethodPatch, "/alumnos?id=eq."+strconv.Itoa(alumno.ID)
	}
	datos := struct {
		Nombre string `json:"nombre"`
		Edad   int    `json:"edad"`
	}{alumno.Nombre, alumno.Edad}
	var guardados []dominio.Alumno
	err := r.solicitar(metodo, ruta, datos, &guardados)
	if err != nil {
		return dominio.Alumno{}, err
	}
	if len(guardados) == 0 && alumno.ID != 0 {
		return dominio.Alumno{}, dominio.ErrNoEncontrado
	}
	if len(guardados) != 1 || guardados[0].ID < 1 {
		return dominio.Alumno{}, fmt.Errorf("%w: Supabase no devolvió el alumno guardado.", dominio.ErrAlmacenamiento)
	}
	return guardados[0], nil
}

func (r *RepositorioSupabase) EliminarAlumno(id int) error {
	// PostgreSQL elimina sus matrículas mediante la clave foránea ON DELETE CASCADE.
	return r.eliminar("alumnos", id)
}
