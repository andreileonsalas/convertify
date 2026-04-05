# Convertify

> Conversor de archivos multimedia en el menú contextual de Windows.  
> Clic derecho → convertir → listo. Sin ventanas. Sin apps abiertas.

[![Build & Release](https://github.com/youruser/convertify/actions/workflows/build.yml/badge.svg)](https://github.com/youruser/convertify/actions/workflows/build.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/youruser/convertify)](https://github.com/youruser/convertify/releases/latest)
[![Platform](https://img.shields.io/badge/platform-Windows-blue)](https://github.com/youruser/convertify/releases)

---

## ¿Cómo funciona?

Haz clic derecho en cualquier imagen, video o audio:

```
📄 video.avi
  └── Clic derecho
        └── Convertify (Video)
              ├── MP4
              ├── MKV
              ├── MOV
              ├── WebM
              ├── GIF (animado)
              ├── Extraer audio → MP3
              ├── Extraer audio → AAC
              ├── Extraer audio → WAV
              └── Extraer audio → FLAC
```

El archivo convertido aparece en la **misma carpeta** que el original.
Una notificación de Windows avisa cuando termina. Nada más.

Si abres `convertify.exe` directamente, se abre una UI mínima con drag & drop.

---

## Descarga

👉 **[Descargar última versión](https://github.com/youruser/convertify/releases/latest)**

| Archivo | Descripción |
|---------|-------------|
| `convertify-setup.exe` | **Recomendado** — Instalador completo. Incluye FFmpeg y registra el menú contextual. |
| `convertify.exe` | Ejecutable solo. Requiere FFmpeg en el PATH o en la misma carpeta. |

---

## Formatos soportados

### Imágenes
| Entrada | Salida |
|---------|--------|
| JPG, PNG, WebP, GIF, BMP, TIFF, AVIF, HEIC, ICO | JPG, PNG, WebP, BMP, TIFF, GIF |

### Videos
| Entrada | Salida |
|---------|--------|
| MP4, AVI, MKV, MOV, WebM, FLV, WMV, M4V, MPG, TS | MP4, AVI, MKV, MOV, WebM, GIF |
| (todos los anteriores) | Extraer audio: MP3, AAC, WAV, FLAC, OGG |

### Audio
| Entrada | Salida |
|---------|--------|
| MP3, WAV, FLAC, AAC, OGG, M4A, WMA, OPUS, AIFF | MP3, WAV, FLAC, AAC, OGG, OPUS, M4A |

---

## Instalación

### Instalador (recomendado)

1. Descarga `convertify-setup.exe` desde [Releases](https://github.com/youruser/convertify/releases/latest)
2. Ejecuta como Administrador
3. Listo — haz clic derecho en cualquier archivo multimedia

El instalador incluye FFmpeg. No necesitas instalar nada más.

### Ejecutable solo

1. Descarga `convertify.exe`
2. Coloca `ffmpeg.exe` y `ffprobe.exe` en la misma carpeta ([descargar FFmpeg](https://ffmpeg.org/download.html))
3. Ejecuta como Administrador para registrar el menú contextual:

```
convertify.exe install
```

### Desinstalar

El desinstalador de Windows lo elimina todo automáticamente.
O manualmente (como Administrador):

```
convertify.exe uninstall
```

---

## Publicar una nueva versión (sin instalar Go)

El proyecto usa GitHub Actions — tú solo subes el código y GitHub lo compila.
Para publicar una release con los `.exe` listos para descargar:

```bash
# 1. Crea y sube el tag de versión
git tag v1.0.0
git push origin v1.0.0

# 2. Ve a: GitHub → tu repo → Releases → "Draft a new release"
# 3. Selecciona el tag que acabas de crear
# 4. Publica la release
# → El Action compila automáticamente y adjunta convertify.exe y convertify-setup.exe
```

Puedes ver el progreso en la pestaña **Actions** del repositorio.

---

## Compilar localmente (opcional)

Si quieres compilar en tu máquina sin depender de CI:

**Requisitos:**
- [Go 1.21+](https://go.dev/dl/)
- `ffmpeg.exe` y `ffprobe.exe` en `assets\` ([descargar](https://ffmpeg.org/download.html))
- [Inno Setup 6](https://jrsoftware.org/isinfo.php) para generar el instalador

```bat
git clone https://github.com/youruser/convertify
cd convertify
mkdir assets
copy C:\ruta\ffmpeg.exe assets\
copy C:\ruta\ffprobe.exe assets\
build.bat
```

El resultado queda en `dist\`.

---

## Estructura del proyecto

```
convertify/
├── .github/
│   └── workflows/
│       └── build.yml        — CI: compila, empaqueta y publica releases
├── formats/
│   └── formats.go           — grupos de archivos y formatos de salida
├── ffmpeg/
│   └── ffmpeg.go            — wrapper de ffmpeg/ffprobe
├── menu/
│   └── menu.go              — registro del menú contextual (registro de Windows)
├── ui/
│   └── ui.go                — interfaz gráfica (controles nativos Win32 via walk)
├── notify/
│   └── notify.go            — notificaciones de Windows
├── main.go                  — punto de entrada, ruteo de comandos
├── convertify.manifest      — DPI awareness y visual styles
├── convertify.rc            — resource file para embed del manifest
├── installer.iss            — script de Inno Setup
├── build.bat                — compilación local en un paso
├── go.mod
├── .gitignore
├── LICENSE
└── README.md
```

---

## Comandos

| Comando | Descripción |
|---------|-------------|
| `convertify.exe` | Abre la interfaz gráfica |
| `convertify.exe convert "file.avi" mp4` | Conversión silenciosa (menú contextual) |
| `convertify.exe install` | Registra los menús contextuales (requiere Administrador) |
| `convertify.exe uninstall` | Elimina los menús contextuales (requiere Administrador) |

---

## Notas técnicas

- FFmpeg se busca en: mismo directorio que el `.exe` → subcarpeta `bin\` → PATH del sistema
- Archivos convertidos se guardan en la **misma carpeta** que el original
- Si el archivo destino ya existe, se añade `(1)`, `(2)`... al nombre
- El menú contextual usa `HKEY_CLASSES_ROOT\SystemFileAssociations\.ext\shell\Convertify`
- UI construida con [walk](https://github.com/lxn/walk) — controles nativos Win32, sin Electron
- Notificaciones vía PowerShell + `System.Windows.Forms.NotifyIcon`

---

## Licencia

[MIT](LICENSE)
