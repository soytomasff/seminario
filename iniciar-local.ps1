# Inicia Aula con datos y autenticación locales, sin una cuenta de Supabase.
$ErrorActionPreference = 'Stop'
$raizAula = $PSScriptRoot
$backendAula = Join-Path $raizAula 'backend'
$frontendAula = Join-Path $raizAula 'frontend'
$logsAula = Join-Path $raizAula 'dist/local'
$null = New-Item -ItemType Directory -Path $logsAula -Force
$entornoAula = Join-Path $backendAula '.env'

if (-not (Test-Path -LiteralPath $entornoAula)) {
    $aleatorioAula = New-Object byte[] 12
    $generadorAula = [Security.Cryptography.RandomNumberGenerator]::Create()
    try { $generadorAula.GetBytes($aleatorioAula) } finally { $generadorAula.Dispose() }
    $claveAula = 'Aula-' + [BitConverter]::ToString($aleatorioAula).Replace('-', '')
    $configuracionAula = @(
        'ALMACENAMIENTO=archivo', 'ARCHIVO_DATOS=datos/academia.json',
        'AUTENTICACION=local', 'AUTH_LOCAL_CORREO=admin@aula.local',
        "AUTH_LOCAL_CLAVE=$claveAula", 'DIRECCION_HTTP=127.0.0.1:8080',
        'ORIGEN_FRONTEND=http://localhost:3000'
    ) -join "`n"
    [IO.File]::WriteAllText($entornoAula, $configuracionAula + "`n", [Text.UTF8Encoding]::new($false))
    [IO.File]::WriteAllText((Join-Path $backendAula 'ACCESO-ADMIN.txt'), "Correo=admin@aula.local`nClave=$claveAula`n", [Text.UTF8Encoding]::new($false))
}

function Test-ServicioAula([string]$Url) {
    try { return (Invoke-WebRequest -Uri $Url -UseBasicParsing -TimeoutSec 3).StatusCode -eq 200 }
    catch { return $false }
}

function Esperar-ServicioAula([string]$Url) {
    for ($intentoAula = 0; $intentoAula -lt 40; $intentoAula++) {
        if (Test-ServicioAula $Url) { return }
        Start-Sleep -Milliseconds 500
    }
    throw "No arranco $Url. Revisa los registros en dist/local."
}

if (-not (Test-ServicioAula 'http://127.0.0.1:8080/healthz')) {
    Push-Location $raizAula
    try {
        go build -o backend/aula-local.exe ./backend
        if ($LASTEXITCODE -ne 0) { throw 'No se pudo compilar el backend.' }
    } finally { Pop-Location }
    $null = Start-Process -FilePath (Join-Path $backendAula 'aula-local.exe') -WorkingDirectory $backendAula -WindowStyle Hidden -RedirectStandardOutput (Join-Path $logsAula 'backend.log') -RedirectStandardError (Join-Path $logsAula 'backend-error.log') -PassThru
    Esperar-ServicioAula 'http://127.0.0.1:8080/healthz'
}

if (-not (Test-ServicioAula 'http://localhost:3000')) {
    if (-not (Test-Path -LiteralPath (Join-Path $frontendAula 'node_modules/next'))) {
        Push-Location $frontendAula
        try {
            npm.cmd ci --no-audit --no-fund
            if ($LASTEXITCODE -ne 0) { throw 'No se pudieron instalar las dependencias.' }
        } finally { Pop-Location }
    }
    $nodeAula = (Get-Command node.exe).Source
    $nextAula = '"' + (Join-Path $frontendAula 'node_modules/next/dist/bin/next') + '"'
    $null = Start-Process -FilePath $nodeAula -ArgumentList @($nextAula, 'dev', '--hostname', '127.0.0.1', '--port', '3000') -WorkingDirectory $frontendAula -WindowStyle Hidden -RedirectStandardOutput (Join-Path $logsAula 'frontend.log') -RedirectStandardError (Join-Path $logsAula 'frontend-error.log') -PassThru
    Esperar-ServicioAula 'http://localhost:3000'
}

Start-Process 'http://localhost:3000'
Write-Host 'Aula disponible en http://localhost:3000'
Write-Host 'Acceso local: backend/ACCESO-ADMIN.txt. Los datos se guardan en backend/datos/academia.json.'
