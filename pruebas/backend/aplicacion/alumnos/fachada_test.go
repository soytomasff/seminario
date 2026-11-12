package alumnos_test

import (
	"errors"
	"testing"

	"seminario4/backend/aplicacion/alumnos"
	"seminario4/backend/dominio"
)

// Este doble implementa solo el contrato de alumnos, sin métodos de matrículas.
type repositorioAlumnos struct {
	guardados []dominio.Alumno
	eliminado int
	fallo     error
}

func (r *repositorioAlumnos) Consultar() (dominio.DatosAcademicos, error) {
	return dominio.DatosAcademicos{Alumnos: []dominio.Alumno{{ID: 7, Nombre: "Ana López", Edad: 20}}}, r.fallo
}

func (r *repositorioAlumnos) GuardarAlumno(a dominio.Alumno) (dominio.Alumno, error) {
	r.guardados = append(r.guardados, a)
	if a.ID == 0 {
		a.ID = 7
	}
	return a, r.fallo
}

func (r *repositorioAlumnos) EliminarAlumno(id int) error {
	r.eliminado = id
	return r.fallo
}

func TestFachadaCoordinaLosCuatroCasosDeUso(t *testing.T) {
	repo := &repositorioAlumnos{}
	f := alumnos.NuevaFachada(repo)
	registrado, err := f.Registrar("  Ana   López  ", 20)
	if err != nil || registrado.ID != 7 || registrado.Nombre != "Ana López" || repo.guardados[0].ID != 0 {
		t.Fatalf("Registro inesperado: %+v, %v", registrado, err)
	}
	lista, err := f.Consultar()
	if err != nil || len(lista) != 1 || lista[0] != registrado {
		t.Fatalf("Consulta inesperada: %+v, %v", lista, err)
	}
	actualizado, err := f.Actualizar(7, "Ana Torres", 21)
	if err != nil || actualizado.ID != 7 || actualizado.Nombre != "Ana Torres" || repo.guardados[1].ID != 7 {
		t.Fatalf("Actualización inesperada: %+v, %v", actualizado, err)
	}
	if err := f.Eliminar(7); err != nil || repo.eliminado != 7 {
		t.Fatal("No se delegó la eliminación del alumno")
	}
}

func TestFachadaValidaAntesDeGuardarYPropagaErrores(t *testing.T) {
	repo := &repositorioAlumnos{}
	f := alumnos.NuevaFachada(repo)
	if _, err := f.Registrar(" ", 20); err == nil {
		t.Fatal("Se aceptó un nombre vacío")
	}
	if _, err := f.Actualizar(7, "Ana López", 121); err == nil {
		t.Fatal("Se aceptó una edad inválida")
	}
	if len(repo.guardados) != 0 {
		t.Fatal("Se guardaron datos que no pasaron la validación")
	}
	repo.fallo = dominio.ErrAlmacenamiento
	_, registro := f.Registrar("Ana López", 20)
	_, consulta := f.Consultar()
	_, actualizacion := f.Actualizar(7, "Ana López", 20)
	for _, err := range []error{registro, consulta, actualizacion, f.Eliminar(7)} {
		if !errors.Is(err, dominio.ErrAlmacenamiento) {
			t.Fatalf("Se perdió el error del repositorio: %v", err)
		}
	}
}
