package matriculas_test

import (
	"errors"
	"testing"

	"seminario4/backend/aplicacion/matriculas"
	"seminario4/backend/dominio"
)

type repositorioMatriculas struct {
	guardadas []dominio.Matricula
	eliminada int
	fallo     error
}

func (r *repositorioMatriculas) Consultar() (dominio.DatosAcademicos, error) {
	return dominio.DatosAcademicos{Matriculas: []dominio.Matricula{
		{ID: 1, AlumnoID: 7, Curso: "Desarrollo Web"},
		{ID: 2, AlumnoID: 7, Curso: "Base de Datos", Modalidad: "Virtual"},
	}}, r.fallo
}

func (r *repositorioMatriculas) GuardarMatricula(m dominio.Matricula) (dominio.Matricula, error) {
	r.guardadas = append(r.guardadas, m)
	if m.ID == 0 {
		m.ID = 10
	}
	return m, r.fallo
}

func (r *repositorioMatriculas) EliminarMatricula(id int) error {
	r.eliminada = id
	return r.fallo
}

// Una segunda operación Visitor demuestra el polimorfismo sin cambiar entidades.
type visitanteEtiquetas struct{}

func (visitanteEtiquetas) VisitarPresencial(m dominio.MatriculaPresencial) string {
	return "Campus: " + m.Datos.Curso
}

func (visitanteEtiquetas) VisitarVirtual(m dominio.MatriculaVirtual) string {
	return "En línea: " + m.Datos.Curso
}

func TestGestorCRUDConVisitanteInyectado(t *testing.T) {
	repo := &repositorioMatriculas{}
	g := matriculas.NuevoGestor(repo, visitanteEtiquetas{})
	creada, err := g.Registrar(7, "Desarrollo Web", "Presencial")
	if err != nil || creada.ID != 10 || creada.Indicacion != "Campus: Desarrollo Web" || repo.guardadas[0].ID != 0 {
		t.Fatalf("Registro inesperado: %+v, %v", creada, err)
	}
	editada, err := g.Actualizar(10, 7, "Base de Datos", "Virtual")
	if err != nil || editada.ID != 10 || editada.Indicacion != "En línea: Base de Datos" || repo.guardadas[1].ID != 10 {
		t.Fatalf("Actualización inesperada: %+v, %v", editada, err)
	}
	lista, err := g.Consultar()
	if err != nil || len(lista) != 2 || lista[0].Modalidad != "Presencial" || lista[0].Indicacion != "Campus: Desarrollo Web" || lista[1].Indicacion != "En línea: Base de Datos" {
		t.Fatalf("Consulta y compatibilidad de datos anteriores: %+v, %v", lista, err)
	}
	if err := g.Eliminar(10); err != nil || repo.eliminada != 10 {
		t.Fatal("No se delegó la eliminación de la matrícula")
	}
}

func TestGestorNoPersisteDatosInvalidosYPropagaErrores(t *testing.T) {
	repo := &repositorioMatriculas{}
	g := matriculas.NuevoGestor(repo, matriculas.VisitanteIndicaciones{})
	if _, err := g.Registrar(0, "Desarrollo Web", "Virtual"); err == nil {
		t.Fatal("Se aceptó un alumno inválido")
	}
	if _, err := g.Actualizar(10, 7, "Curso desconocido", "Virtual"); err == nil {
		t.Fatal("Se aceptó un curso inválido")
	}
	if len(repo.guardadas) != 0 {
		t.Fatal("Se persistieron datos inválidos")
	}
	repo.fallo = dominio.ErrDuplicado
	_, registro := g.Registrar(7, "Desarrollo Web", "Virtual")
	_, actualizacion := g.Actualizar(10, 7, "Desarrollo Web", "Virtual")
	_, consulta := g.Consultar()
	for _, err := range []error{registro, actualizacion, consulta, g.Eliminar(10)} {
		if !errors.Is(err, dominio.ErrDuplicado) {
			t.Fatalf("Se perdió el error del repositorio: %v", err)
		}
	}
}
