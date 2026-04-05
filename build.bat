@echo off
REM ============================================================
REM  Convertify - Build Script
REM  Requires: Go 1.21+, windres (MinGW), Inno Setup (optional)
REM ============================================================

echo [1/4] Descargando dependencias...
go mod tidy
if %ERRORLEVEL% neq 0 goto :error

echo [2/4] Compilando recurso (manifest)...
REM windres es parte de MinGW. Si no lo tienes, salta este paso.
REM El .syso embeds el manifest en el exe para visual styles y DPI.
where windres >nul 2>&1
if %ERRORLEVEL% equ 0 (
    windres convertify.rc -O coff -o convertify.syso
) else (
    echo    windres no encontrado, omitiendo embed del manifest.
)

echo [3/4] Compilando convertify.exe...
REM -ldflags "-H windowsgui" oculta la consola cuando se lanza desde el escritorio
REM -ldflags "-s -w" reduce el tamanio del binario
go build -ldflags "-H windowsgui -s -w" -o dist\convertify.exe .
if %ERRORLEVEL% neq 0 goto :error

echo [4/4] Exe generado en dist\convertify.exe
echo.

REM Opcional: generar instalador con Inno Setup
where iscc >nul 2>&1
if %ERRORLEVEL% equ 0 (
    echo Generando instalador con Inno Setup...
    iscc installer.iss
    echo Instalador generado en dist\convertify-setup.exe
) else (
    echo Inno Setup no encontrado. Para generar el instalador:
    echo   1. Descarga Inno Setup de https://jrsoftware.org/isinfo.php
    echo   2. Ejecuta: iscc installer.iss
)

echo.
echo === Build completado ===
goto :end

:error
echo ERROR: La compilacion fallo.
exit /b 1

:end
