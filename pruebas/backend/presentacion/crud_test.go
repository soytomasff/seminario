package presentacion_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"seminario4/backend/aplicacion"
	"seminario4/backend/aplicacion/matriculas"
	"seminario4/backend/dominio"
	"seminario4/backend/infraestructura"
	"seminario4/backend/presentacion"
)

func preparar(t *testing.T) (http.Handler, *infraestructura.RepositorioArchivo, string) {
	t.Helper()
	ruta := filepath.Join(t.TempDir(), "academia.json")
	repo, err := infraestructura.NuevoRepositorioArchivo(ruta)
	if err != nil {
		t.Fatal(err)
	}
	return presentacion.NuevoManejador(aplicacion.NuevoServicio(repo), "http://localhost:3000"), repo, ruta
}

func solicitar(t *testing.T, h http.Handler, metodo, ruta string, datos any, estado int, destino any) {
	t.Helper()
	var cuerpo []byte
	if datos != nil {
		var err error
		cuerpo, err = json.Marshal(datos)
		if err != nil {
			t.Fatal(err)
		}
	}
	r := httptest.NewRequest(metodo, ruta, bytes.NewReader(cuerpo))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != estado {
		t.Fatalf("%s %s: esperado %d, recibido %d: %s", metodo, ruta, estado, w.Code, w.Body.String())
	}
	if estado == http.StatusNoContent && w.Body.Len() != 0 {
		t.Fatal("HTTP 204 debe tener el cuerpo vacío")
	}
	if destino != nil {
		if err := json.Unmarshal(w.Body.Bytes(), destino); err != nil {
			t.Fatal(err)
		}
	}
}

func TestControladoresCRUDPersistenciaYCascada(t *testing.T) {
	h, repo, ruta := preparar(t)
	inicial, err := repo.Consultar()
	if err != nil {
		t.Fatal(err)
	}
	var alumno dominio.Alumno
	solicitar(t, h, "POST", "/api/alumnos", map[string]any{"nombre": "  María   García  ", "edad": 20}, 201, &alumno)
	if alumno.ID < 1 || alumno.Nombre != "María García" {
		t.Fatalf("Respuesta de creación inesperada: %+v", alumno)
	}
	rutaAlumno := fmt.Sprintf("/api/alumnos/%d", alumno.ID)
	solicitar(t, h, "PUT", rutaAlumno, map[string]any{"nombre": "María Torres", "edad": 21}, 200, &alumno)
	if alumno.Nombre != "María Torres" || alumno.Edad != 21 {
		t.Fatal("No se actualizaron los datos del alumno")
	}
	var alumnos []dominio.Alumno
	solicitar(t, h, "GET", "/api/alumnos", nil, 200, &alumnos)
	if len(alumnos) != len(inicial.Alumnos)+1 || alumnos[len(alumnos)-1] != alumno {
		t.Fatal("El alumno no aparece en la consulta")
	}
	entrada := map[string]any{"alumno_id": alumno.ID, "curso": "Desarrollo Web", "modalidad": "Presencial"}
	var primera matriculas.DetalleMatricula
	solicitar(t, h, "POST", "/api/matriculas", entrada, 201, &primera)
	if primera.Indicacion != "Asiste al instituto para el curso Desarrollo Web." {
		t.Fatalf("Visitor presencial inesperado: %+v", primera)
	}
	entrada["modalidad"] = "Virtual"
	rutaMatricula := fmt.Sprintf("/api/matriculas/%d", primera.ID)
	solicitar(t, h, "PUT", rutaMatricula, entrada, 200, &primera)
	if primera.Indicacion != "Ingresa al aula virtual para el curso Desarrollo Web." || primera.Modalidad != "Virtual" {
		t.Fatal("Visitor no recalculó la indicación al actualizar")
	}
	solicitar(t, h, "POST", "/api/matriculas", entrada, 409, nil)
	var lista []matriculas.DetalleMatricula
	solicitar(t, h, "GET", "/api/matriculas", nil, 200, &lista)
	if len(lista) != len(inicial.Matriculas)+1 || lista[len(lista)-1] != primera {
		t.Fatal("La consulta de matrículas cambió su contrato JSON")
	}
	entrada["curso"] = "Base de Datos"
	var segunda matriculas.DetalleMatricula
	solicitar(t, h, "POST", "/api/matriculas", entrada, 201, &segunda)
	solicitar(t, h, "DELETE", rutaMatricula, nil, 204, nil)
	intermedio, err := repo.Consultar()
	if err != nil || len(intermedio.Alumnos) != len(inicial.Alumnos)+1 || len(intermedio.Matriculas) != len(inicial.Matriculas)+1 {
		t.Fatal("Eliminar una matrícula debe conservar al alumno y sus otras matrículas")
	}
	solicitar(t, h, "DELETE", rutaAlumno, nil, 204, nil)
	solicitar(t, h, "DELETE", rutaAlumno, nil, 404, nil)
	reabierto, err := infraestructura.NuevoRepositorioArchivo(ruta)
	if err != nil {
		t.Fatal(err)
	}
	final, err := reabierto.Consultar()
	if err != nil || !reflect.DeepEqual(final.Alumnos, inicial.Alumnos) || !reflect.DeepEqual(final.Matriculas, inicial.Matriculas) {
		t.Fatal("El borrado en cascada o la persistencia afectaron otros registros")
	}
	var resumen struct {
		Alumnos    []dominio.Alumno              `json:"alumnos"`
		Matriculas []matriculas.DetalleMatricula `json:"matriculas"`
		Cursos     []string                      `json:"cursos"`
	}
	solicitar(t, h, "GET", "/api/resumen", nil, 200, &resumen)
	if len(resumen.Alumnos) != len(inicial.Alumnos) || len(resumen.Matriculas) != len(inicial.Matriculas) || len(resumen.Cursos) != len(dominio.CursosDisponibles()) {
		t.Fatal("El resumen perdió información de un módulo")
	}
}

func TestValidacionesHTTPYCORSNoModificanDatos(t *testing.T) {
	h, repo, _ := preparar(t)
	antes, err := repo.Consultar()
	if err != nil {
		t.Fatal(err)
	}
	for _, caso := range []struct {
		metodo, ruta, cuerpo string
		estado               int
	}{
		{"POST", "/api/alumnos", `{"nombre":`, 400},
		{"POST", "/api/alumnos", `{"nombre":"Ana López","edad":20,"extra":1}`, 400},
		{"POST", "/api/alumnos", `{"nombre":"Ana López","edad":121}`, 400},
		{"PUT", "/api/alumnos/abc", `{}`, 400},
		{"PUT", "/api/alumnos/999999", `{"nombre":"Ana López","edad":20}`, 404},
		{"POST", "/api/matriculas", `{"alumno_id":999999,"curso":"Desarrollo Web","modalidad":"Virtual"}`, 400},
		{"PUT", "/api/matriculas/999999", `{"alumno_id":1,"curso":"Desarrollo Web","modalidad":"Virtual"}`, 404},
	} {
		r := httptest.NewRequest(caso.metodo, caso.ruta, strings.NewReader(caso.cuerpo))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != caso.estado {
			t.Fatalf("%s %s: esperado %d, recibido %d", caso.metodo, caso.ruta, caso.estado, w.Code)
		}
	}
	r := httptest.NewRequest("POST", "/api/alumnos", strings.NewReader(`{"nombre":"Ana López","edad":20}`))
	r.Header.Set("Origin", "https://ajeno.test")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatal("Se permitió una escritura desde otro origen")
	}
	despues, err := repo.Consultar()
	if err != nil || !reflect.DeepEqual(antes, despues) {
		t.Fatal("Una petición rechazada modificó los datos")
	}
}
