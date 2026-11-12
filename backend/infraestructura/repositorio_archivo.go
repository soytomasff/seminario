// Infraestructura implementa el contrato del dominio usando un archivo JSON.
package infraestructura

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"seminario4/backend/dominio"
)

type RepositorioArchivo struct {
	ruta string
	// Mutex permite atender varias peticiones sin que escriban a la vez.
	// Este ejemplo se ejecuta en UN proceso de Go, no en varios servidores.
	mutex sync.Mutex
}

func NuevoRepositorioArchivo(ruta string) (*RepositorioArchivo, error) {
	repositorio := &RepositorioArchivo{ruta: ruta}
	if err := os.MkdirAll(filepath.Dir(ruta), 0755); err != nil {
		return nil, err
	}
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		if err := repositorio.escribir(datosDeEjemplo()); err != nil {
			return nil, err
		}
	}
	// Si el archivo existe pero está dañado, avisamos; nunca lo reemplazamos en silencio.
	if _, err := repositorio.leer(); err != nil {
		return nil, err
	}
	return repositorio, nil
}

func datosDeEjemplo() dominio.DatosAcademicos {
	return dominio.DatosAcademicos{
		Alumnos: []dominio.Alumno{
			{ID: 1, Nombre: "Juan Pérez", Edad: 20},
			{ID: 2, Nombre: "Ana López", Edad: 19},
		},
		Matriculas: []dominio.Matricula{
			{ID: 1, AlumnoID: 1, Curso: "Programación I", Modalidad: "Presencial"},
			{ID: 2, AlumnoID: 2, Curso: "Desarrollo Web", Modalidad: "Virtual"},
		},
		SiguienteAlumno: 3, SiguienteMatricula: 3,
	}
}

func (r *RepositorioArchivo) leer() (dominio.DatosAcademicos, error) {
	var datos dominio.DatosAcademicos
	contenido, err := os.ReadFile(r.ruta)
	if err == nil {
		err = json.Unmarshal(contenido, &datos)
	}
	if err != nil {
		return datos, fmt.Errorf("%w: %v", dominio.ErrAlmacenamiento, err)
	}
	return datos, nil
}

func (r *RepositorioArchivo) escribir(datos dominio.DatosAcademicos) error {
	// Primero escribimos en un archivo temporal. Si falla, conservamos el anterior.
	temporal, err := os.CreateTemp(filepath.Dir(r.ruta), "academia-*.tmp")
	if err != nil {
		return fmt.Errorf("%w: %v", dominio.ErrAlmacenamiento, err)
	}
	defer os.Remove(temporal.Name())
	codificador := json.NewEncoder(temporal)
	codificador.SetIndent("", "  ")
	err = codificador.Encode(datos)
	errCierre := temporal.Close()
	if err == nil {
		err = errCierre
	}
	if err == nil {
		err = os.Rename(temporal.Name(), r.ruta)
	}
	if err != nil {
		return fmt.Errorf("%w: %v", dominio.ErrAlmacenamiento, err)
	}
	return nil
}

// Este ayudante mantiene juntos leer + cambiar + guardar bajo el mismo bloqueo.
// cambio es una función que recibimos para indicar qué registro se modificará.
func (r *RepositorioArchivo) modificar(cambio func(*dominio.DatosAcademicos) error) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	datos, err := r.leer()
	if err != nil {
		return err
	}
	if err := cambio(&datos); err != nil {
		return err
	}
	return r.escribir(datos)
}

func (r *RepositorioArchivo) Consultar() (dominio.DatosAcademicos, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.leer()
}

func (r *RepositorioArchivo) Reiniciar() error {
	return r.modificar(func(datos *dominio.DatosAcademicos) error {
		*datos = datosDeEjemplo()
		return nil
	})
}

var _ dominio.RepositorioAcademico = (*RepositorioArchivo)(nil)
