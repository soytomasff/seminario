package dominio

import "errors"

var ErrNoEncontrado = errors.New("El registro no existe.")
var ErrDuplicado = errors.New("El alumno ya está matriculado en ese curso.")
var ErrAlumnoNoExiste = errors.New("El alumno elegido no está registrado.")
var ErrAlmacenamiento = errors.New("No se pudieron leer o guardar los datos.")
