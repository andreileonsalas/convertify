package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"convertify/ffmpeg"
	"convertify/formats"
	"convertify/menu"
	"convertify/notify"
	"convertify/ui"
)

func main() {
	args := os.Args[1:]

	switch {
	case len(args) == 0:
		// No arguments → launch GUI
		runUI()

	case args[0] == "convert" && len(args) == 3:
		// Called from context menu:  convert <filepath> <ext>
		runSilentConvert(args[1], args[2])

	case args[0] == "install":
		// Register context menus (needs elevation)
		runInstall()

	case args[0] == "uninstall":
		// Remove context menus (needs elevation)
		runUninstall()

	default:
		fmt.Fprintln(os.Stderr, "Uso:")
		fmt.Fprintln(os.Stderr, "  convertify.exe                    — abre la interfaz gráfica")
		fmt.Fprintln(os.Stderr, "  convertify.exe convert <file> <ext> — convierte un archivo (menú contextual)")
		fmt.Fprintln(os.Stderr, "  convertify.exe install            — registra el menú contextual (requiere admin)")
		fmt.Fprintln(os.Stderr, "  convertify.exe uninstall          — elimina el menú contextual (requiere admin)")
		os.Exit(1)
	}
}

// runUI launches the graphical interface
func runUI() {
	ui.Run()
}

// runSilentConvert is called from the Windows context menu.
// It converts the file silently and shows a toast notification when done.
func runSilentConvert(inputPath, targetExt string) {
	targetExt = strings.ToLower(targetExt)

	// Detect input group
	group := formats.DetectGroup(inputPath)
	if group == formats.Unknown {
		notify.Show("Convertify - Error", "Formato de archivo no reconocido")
		os.Exit(1)
	}

	// Find the matching output format
	outFmt, ok := formats.FindFormat(group, targetExt)
	if !ok {
		notify.Show("Convertify - Error", "Formato de salida desconocido: "+targetExt)
		os.Exit(1)
	}

	// Show "working" notification for long conversions
	baseName := filepath.Base(inputPath)
	notify.Show("Convertify", fmt.Sprintf("Convirtiendo %s...", baseName))

	outPath, err := ffmpeg.Convert(inputPath, outFmt)
	if err != nil {
		notify.Show("Convertify - Error", err.Error())
		os.Exit(1)
	}

	outName := filepath.Base(outPath)
	notify.Show("Convertify ✓", fmt.Sprintf("%s listo", outName))
}

// runInstall registers context menus in the Windows registry
func runInstall() {
	fmt.Println("Instalando menús contextuales...")
	if err := menu.Install(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Asegúrate de ejecutar como Administrador.")
		os.Exit(1)
	}
	fmt.Println("✓  Menús contextuales instalados correctamente.")
	fmt.Println("   Haz clic derecho en cualquier imagen, video o audio para ver 'Convertify'.")
}

// runUninstall removes all Convertify registry entries
func runUninstall() {
	fmt.Println("Desinstalando menús contextuales...")
	if err := menu.Uninstall(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓  Menús contextuales eliminados.")
}
