package infraestructura_test

import (
	"context"
	"errors"
	"testing"

	"seminario4/backend/aplicacion"
	"seminario4/backend/infraestructura"
)

func TestAutenticacionLocalCredencialesRotacionYCierre(t *testing.T) {
	auth, err := infraestructura.NuevaAutenticacionLocal("admin@aula.local", "clave-local-de-prueba", "127.0.0.1:8080", "http://localhost:3000")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	for _, caso := range []struct{ correo, clave string }{
		{"admin@aula.local", "incorrecta"},
		{"otro@aula.local", "clave-local-de-prueba"},
	} {
		if _, err := auth.IniciarSesion(ctx, caso.correo, caso.clave); !errors.Is(err, aplicacion.ErrCredenciales) {
			t.Fatal("Se aceptaron credenciales incorrectas")
		}
	}
	primera, err := auth.IniciarSesion(ctx, "admin@aula.local", "clave-local-de-prueba")
	if err != nil {
		t.Fatal(err)
	}
	segunda, err := auth.IniciarSesion(ctx, "admin@aula.local", "clave-local-de-prueba")
	if err != nil || primera.Acceso == segunda.Acceso {
		t.Fatal("Las sesiones deben ser independientes")
	}
	usuario, err := auth.VerificarSesion(ctx, primera.Acceso)
	if err != nil || usuario.ID != infraestructura.IDAdministradorLocal {
		t.Fatal("No se verificó al administrador local")
	}
	for _, token := range []string{"inventado", primera.Refresco} {
		if _, err := auth.VerificarSesion(ctx, token); !errors.Is(err, aplicacion.ErrCredenciales) {
			t.Fatal("Se aceptó un token de acceso inválido")
		}
	}
	if _, err := auth.RenovarSesion(ctx, primera.Acceso); !errors.Is(err, aplicacion.ErrCredenciales) {
		t.Fatal("Se usó el token de acceso como refresco")
	}
	nueva, err := auth.RenovarSesion(ctx, primera.Refresco)
	if err != nil || nueva.Acceso == primera.Acceso || nueva.Refresco == primera.Refresco {
		t.Fatal("No se rotaron los tokens")
	}
	if _, err := auth.VerificarSesion(ctx, primera.Acceso); !errors.Is(err, aplicacion.ErrCredenciales) {
		t.Fatal("Se reutilizó el acceso anterior")
	}
	if _, err := auth.RenovarSesion(ctx, primera.Refresco); !errors.Is(err, aplicacion.ErrCredenciales) {
		t.Fatal("Se reutilizó el refresco anterior")
	}
	if err := auth.CerrarSesion(ctx, nueva.Acceso); err != nil {
		t.Fatal(err)
	}
	if _, err := auth.VerificarSesion(ctx, nueva.Acceso); !errors.Is(err, aplicacion.ErrCredenciales) {
		t.Fatal("Se reutilizó una sesión cerrada")
	}
	if _, err := auth.RenovarSesion(ctx, nueva.Refresco); !errors.Is(err, aplicacion.ErrCredenciales) {
		t.Fatal("El cierre no revocó el refresco")
	}
	if _, err := auth.VerificarSesion(ctx, segunda.Acceso); err != nil {
		t.Fatal("Cerrar una sesión afectó a otra")
	}
}

func TestAutenticacionLocalExigeLoopbackYClave(t *testing.T) {
	for _, caso := range []struct{ direccion, origen string }{
		{"0.0.0.0:8080", "http://localhost:3000"},
		{":8080", "http://localhost:3000"},
		{"127.0.0.1:8080", "https://ejemplo.com"},
		{"127.0.0.1:8080", "http://localhost:3000@ejemplo.com"},
	} {
		if _, err := infraestructura.NuevaAutenticacionLocal("admin@aula.local", "clave-local-de-prueba", caso.direccion, caso.origen); err == nil {
			t.Errorf("Se permitió acceso externo: %+v", caso)
		}
	}
	if _, err := infraestructura.NuevaAutenticacionLocal("admin@aula.local", "corta", "127.0.0.1:8080", "http://localhost:3000"); err == nil {
		t.Fatal("Se aceptó una contraseña corta")
	}
}
