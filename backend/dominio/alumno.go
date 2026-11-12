package dominio

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// Alumno es una entidad identificada por ID.
type Alumno struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Edad   int    `json:"edad"`
}

// NuevoAlumno normaliza y valida los datos antes de guardar la entidad.
func NuevoAlumno(nombre string, edad int) (Alumno, error) {
	nombre = strings.Join(strings.Fields(nombre), " ")
	longitud := utf8.RuneCountInString(nombre)
	if longitud < 3 || longitud > 100 {
		return Alumno{}, errors.New("El nombre debe tener entre 3 y 100 caracteres.")
	}
	if edad < 1 || edad > 120 {
		return Alumno{}, errors.New("La edad debe ser un número entero entre 1 y 120.")
	}
	return Alumno{Nombre: nombre, Edad: edad}, nil
}
