// Requiere iniciar-local.ps1 en ejecución. Usa solo registros temporales propios.
import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';

const texto = await readFile(new URL('../backend/ACCESO-ADMIN.txt', import.meta.url), 'utf8');
const credenciales = Object.fromEntries(texto.trim().split(/\r?\n/).map(linea => {
  const indice = linea.indexOf('=');
  return [linea.slice(0, indice), linea.slice(indice + 1)];
}));
let cookie = '';
let alumno;
let matricula;
async function pedir(metodo, ruta, datos, estado = 200) {
  const r = await fetch(`http://localhost:8080/api${ruta}`, {
    method: metodo,
    headers: { 'Content-Type': 'application/json', Origin: 'http://localhost:3000', ...(cookie ? { Cookie: cookie } : {}) },
    body: datos === undefined ? undefined : JSON.stringify(datos),
    signal: AbortSignal.timeout(10000),
  });
  assert.equal(r.status, estado, `${metodo} ${ruta}`);
  return r;
}

await pedir('GET', '/resumen', undefined, 401);
await pedir('POST', '/auth/login', { correo: credenciales.Correo, clave: 'incorrecta' }, 401);
const acceso = await pedir('POST', '/auth/login', { correo: credenciales.Correo, clave: credenciales.Clave });
assert.match(acceso.headers.get('set-cookie'), /HttpOnly/i);
assert.equal(acceso.headers.get('access-control-allow-credentials'), 'true');
cookie = acceso.headers.get('set-cookie').split(';')[0];
const inicial = await (await pedir('GET', '/resumen')).json();
console.log('OK: login local, cookie protegida y consulta autenticada.');
try {
  alumno = await (await pedir('POST', '/alumnos', { nombre: `Prueba local ${Date.now()}`, edad: 20 }, 201)).json();
  await pedir('PUT', `/alumnos/${alumno.id}`, { nombre: alumno.nombre, edad: 21 });
  matricula = await (await pedir('POST', '/matriculas', { alumno_id: alumno.id, curso: inicial.cursos[0], modalidad: 'Presencial' }, 201)).json();
  const editada = await (await pedir('PUT', `/matriculas/${matricula.id}`, { alumno_id: alumno.id, curso: inicial.cursos[0], modalidad: 'Virtual' })).json();
  assert.equal(editada.modalidad, 'Virtual');
  assert.match(editada.indicacion, /instituto|curso/);
  await pedir('POST', '/matriculas', { alumno_id: alumno.id, curso: inicial.cursos[0], modalidad: 'Virtual' }, 409);
  console.log('OK: crear y editar alumnos y matrículas; rechazo de duplicados.');
} finally {
  try {
    if (matricula) await pedir('DELETE', `/matriculas/${matricula.id}`, undefined, 204);
    if (alumno) await pedir('DELETE', `/alumnos/${alumno.id}`, undefined, 204);
    const final = await (await pedir('GET', '/resumen')).json();
    assert.deepEqual(final.alumnos, inicial.alumnos);
    assert.deepEqual(final.matriculas, inicial.matriculas);
    console.log('OK: eliminación y conservación de los registros iniciales.');
  } finally {
    await pedir('POST', '/auth/logout', undefined, 204);
  }
}
await pedir('GET', '/resumen', undefined, 401);
console.log('OK: logout invalida la sesión.');
