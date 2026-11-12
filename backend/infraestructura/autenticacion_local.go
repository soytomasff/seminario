package infraestructura

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"seminario4/backend/aplicacion"
)

const IDAdministradorLocal = "administrador-local"

type tokenLocal struct {
	acceso, refresco [32]byte
	vence            time.Time
}

// AutenticacionLocal permite desarrollar sin servicios externos. El constructor
// exige que tanto la API como el frontend estén limitados al equipo local.
type AutenticacionLocal struct {
	usuario aplicacion.UsuarioAutenticado
	clave   [32]byte
	mu      sync.Mutex
	tokens  map[[32]byte]*tokenLocal
}

var _ aplicacion.Autenticador = (*AutenticacionLocal)(nil)

func NuevaAutenticacionLocal(correo, clave, direccion, origen string) (*AutenticacionLocal, error) {
	host, _, err := net.SplitHostPort(direccion)
	u, errURL := url.Parse(origen)
	if err != nil || !esHostLocal(host) || errURL != nil || !esHostLocal(u.Hostname()) || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
		return nil, errors.New("La autenticación local requiere DIRECCION_HTTP y ORIGEN_FRONTEND de localhost o loopback.")
	}
	correo = strings.ToLower(strings.TrimSpace(correo))
	if !strings.Contains(correo, "@") || len(correo) > 254 || len(clave) < 12 || len(clave) > 1024 {
		return nil, errors.New("Configura AUTH_LOCAL_CORREO y AUTH_LOCAL_CLAVE con una contraseña de al menos 12 caracteres.")
	}
	return &AutenticacionLocal{
		usuario: aplicacion.UsuarioAutenticado{ID: IDAdministradorLocal, Correo: correo},
		clave:   sha256.Sum256([]byte(clave)),
		tokens:  make(map[[32]byte]*tokenLocal),
	}, nil
}

func esHostLocal(host string) bool {
	return host == "localhost" || net.ParseIP(host).IsLoopback()
}

func (a *AutenticacionLocal) IniciarSesion(ctx context.Context, correo, clave string) (aplicacion.SesionAutenticada, error) {
	if err := ctx.Err(); err != nil {
		return aplicacion.SesionAutenticada{}, err
	}
	huella := sha256.Sum256([]byte(clave))
	if subtle.ConstantTimeCompare(huella[:], a.clave[:]) != 1 || strings.ToLower(strings.TrimSpace(correo)) != a.usuario.Correo {
		return aplicacion.SesionAutenticada{}, aplicacion.ErrCredenciales
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.crearSesion()
}

// crearSesion se llama bajo el mutex; cada sesión tiene tokens independientes.
func (a *AutenticacionLocal) crearSesion() (aplicacion.SesionAutenticada, error) {
	ahora := time.Now()
	for clave, token := range a.tokens {
		if !ahora.Before(token.vence) {
			delete(a.tokens, clave)
		}
	}
	if len(a.tokens) >= 2000 {
		return aplicacion.SesionAutenticada{}, aplicacion.ErrServicioAutenticacion
	}
	var aleatorio [64]byte
	if _, err := rand.Read(aleatorio[:]); err != nil {
		return aplicacion.SesionAutenticada{}, aplicacion.ErrServicioAutenticacion
	}
	acceso := base64.RawURLEncoding.EncodeToString(aleatorio[:32])
	refresco := base64.RawURLEncoding.EncodeToString(aleatorio[32:])
	token := &tokenLocal{acceso: sha256.Sum256([]byte(acceso)), refresco: sha256.Sum256([]byte(refresco)), vence: ahora.Add(time.Hour)}
	a.tokens[token.acceso], a.tokens[token.refresco] = token, token
	return aplicacion.SesionAutenticada{Acceso: acceso, Refresco: refresco, VenceEn: 3600, Usuario: a.usuario}, nil
}

func (a *AutenticacionLocal) VerificarSesion(ctx context.Context, acceso string) (aplicacion.UsuarioAutenticado, error) {
	if err := ctx.Err(); err != nil {
		return aplicacion.UsuarioAutenticado{}, err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	huella := sha256.Sum256([]byte(acceso))
	token := a.tokens[huella]
	if token == nil || token.acceso != huella || !time.Now().Before(token.vence) {
		return aplicacion.UsuarioAutenticado{}, aplicacion.ErrCredenciales
	}
	return a.usuario, nil
}

func (a *AutenticacionLocal) RenovarSesion(ctx context.Context, refresco string) (aplicacion.SesionAutenticada, error) {
	if err := ctx.Err(); err != nil {
		return aplicacion.SesionAutenticada{}, err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	huella := sha256.Sum256([]byte(refresco))
	token := a.tokens[huella]
	if token == nil || token.refresco != huella || !time.Now().Before(token.vence) {
		return aplicacion.SesionAutenticada{}, aplicacion.ErrCredenciales
	}
	nueva, err := a.crearSesion()
	if err != nil {
		return aplicacion.SesionAutenticada{}, err
	}
	delete(a.tokens, token.acceso)
	delete(a.tokens, token.refresco)
	return nueva, nil
}

func (a *AutenticacionLocal) CerrarSesion(ctx context.Context, acceso string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	huella := sha256.Sum256([]byte(acceso))
	if token := a.tokens[huella]; token != nil && token.acceso == huella {
		delete(a.tokens, token.acceso)
		delete(a.tokens, token.refresco)
	}
	return nil
}
