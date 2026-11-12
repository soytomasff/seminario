# Aula académica

## Estructura

| Carpeta | Responsabilidad |
| --- | --- |
| `backend/dominio` | Entidades, reglas de validación, errores y contrato de persistencia. No depende de HTTP, Supabase ni archivos. |
| `backend/aplicacion` | Casos de uso de alumnos y matrículas; coordina el dominio con el repositorio y define el contrato de autenticación. |
| `backend/infraestructura` | Implementaciones de persistencia, autenticación con Supabase y lectura de configuración. |
| `backend/presentacion` | Adaptadores HTTP, respuestas y control de acceso. |
| `backend/main.go` | Punto de entrada que conecta las capas. |
| `frontend` | Interfaz React. |
| `pruebas/backend` | Pruebas Go organizadas por la capa que verifican. |
| `pruebas/*.mjs` | Pruebas de integración, navegador y despliegue. |
| `deploy/aws` | Configuración y scripts de despliegue. |

Las dependencias van hacia aplicación y dominio: presentación e infraestructura usan sus contratos; `main.go` elige las implementaciones. El módulo Go está en la raíz para que `go test ./...` incluya las pruebas separadas del código de producción.

## Módulos y patrones

- `backend/aplicacion/alumnos`: `fachada.go` ofrece el CRUD y delega en `registrar.go`, `consultar.go`, `actualizar.go` y `eliminar.go`.
- `backend/aplicacion/matriculas`: `gestor.go` recibe el repositorio y el visitante; las cuatro operaciones CRUD tienen archivos propios. `detalle.go` define la respuesta y `visitante_indicaciones.go` implementa Visitor.
- `backend/dominio`: las interfaces `ElementoMatricula` y `VisitanteMatricula`, los elementos presencial y virtual y su selección están en archivos separados. Cada CRUD dispone de su contrato de repositorio.
- `backend/presentacion`: los controladores de alumnos, matrículas y sistema están separados. `http.go` conecta las rutas; `solicitudes.go`, `respuestas.go` y `cors.go` contienen las tareas HTTP comunes.
- `backend/infraestructura`: cada almacenamiento separa sus operaciones en archivos terminados en `_alumnos.go` y `_matriculas.go`; comparten la misma instancia de repositorio.

La estructura completa y el recorrido de los patrones están en [docs/ARQUITECTURA.md](docs/ARQUITECTURA.md).

## Desarrollo

Para abrir el sistema completo sin Supabase, ejecuta desde la raíz:

```powershell
powershell -ExecutionPolicy Bypass -File .\iniciar-local.ps1
```

El script genera una contraseña aleatoria en `backend/ACCESO-ADMIN.txt` y configura `backend/.env` cuando todavía no existe. Inicia la API y el frontend y abre `http://localhost:3000`. Los datos se guardan en `backend/datos/academia.json`. La autenticación local acepta únicamente direcciones de loopback para la API y el frontend. Para volver a abrirlo, ejecuta el mismo script; conserva los datos y las credenciales existentes.

Desde `backend`, configura `.env` a partir de `.env.example` y ejecuta `go run .`. Desde `frontend`, ejecuta `npm ci` y `npm run dev`. Consulta [la guía de Supabase](supabase/LEEME.md) y [la guía de AWS](deploy/aws/LEEME.md) para configurar esos servicios.

## Pruebas

Desde la raíz del repositorio:

```powershell
go test ./...
```

Las pruebas Go usan servidores HTTP simulados y directorios temporales; no requieren credenciales reales. Los scripts `pruebas/*.mjs` requieren servicios en ejecución según cada caso.
