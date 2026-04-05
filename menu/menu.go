package menu

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"

	"convertify/formats"
)

// Install registers Convertify in the Windows right-click context menu
// for all supported file types. Requires elevation (run as Administrator).
func Install() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("no se pudo obtener la ruta del ejecutable: %w", err)
	}
	exePath, _ = filepath.Abs(exePath)

	// Register for each photo extension
	for _, ext := range formats.PhotoExtensions() {
		if err := registerGroup(exePath, ext, formats.Photo, formats.PhotoOutputs); err != nil {
			return fmt.Errorf("error registrando .%s: %w", ext, err)
		}
	}

	// Register for each video extension (video outputs + extract outputs)
	videoAll := append(formats.VideoOutputs, formats.VideoExtractOutputs...)
	for _, ext := range formats.VideoExtensions() {
		if err := registerGroup(exePath, ext, formats.Video, videoAll); err != nil {
			return fmt.Errorf("error registrando .%s: %w", ext, err)
		}
	}

	// Register for each audio extension
	for _, ext := range formats.AudioExtensions() {
		if err := registerGroup(exePath, ext, formats.Audio, formats.AudioOutputs); err != nil {
			return fmt.Errorf("error registrando .%s: %w", ext, err)
		}
	}

	return nil
}

// Uninstall removes all Convertify context menu entries
func Uninstall() error {
	allExts := append(formats.PhotoExtensions(), formats.VideoExtensions()...)
	allExts = append(allExts, formats.AudioExtensions()...)

	for _, ext := range allExts {
		keyPath := fmt.Sprintf(`SystemFileAssociations\.%s\shell\Convertify`, ext)
		if err := deleteKeyRecursive(registry.CLASSES_ROOT, keyPath); err != nil {
			// Ignore "key not found" errors
			fmt.Printf("Advertencia al eliminar .%s: %v\n", ext, err)
		}
	}
	return nil
}

// registerGroup creates cascading context menu entries for a file extension
func registerGroup(exePath, ext string, group formats.Group, outputs []formats.OutputFormat) error {
	// Root key for this extension's Convertify menu
	rootPath := fmt.Sprintf(`SystemFileAssociations\.%s\shell\Convertify`, ext)

	// Create/open the parent key
	parentKey, _, err := registry.CreateKey(
		registry.CLASSES_ROOT,
		rootPath,
		registry.SET_VALUE|registry.CREATE_SUB_KEY,
	)
	if err != nil {
		return fmt.Errorf("no se pudo crear la clave raíz: %w", err)
	}
	defer parentKey.Close()

	// Set display name and mark as cascade menu
	label := fmt.Sprintf("Convertify (%s)", group.String())
	if err := parentKey.SetStringValue("MUIVerb", label); err != nil {
		return err
	}
	// Empty SubCommands signals Windows to use the shell subkey
	if err := parentKey.SetStringValue("SubCommands", ""); err != nil {
		return err
	}

	// Create the shell subkey that holds all child items
	shellKey, _, err := registry.CreateKey(
		registry.CLASSES_ROOT,
		rootPath+`\shell`,
		registry.SET_VALUE|registry.CREATE_SUB_KEY,
	)
	if err != nil {
		return err
	}
	defer shellKey.Close()

	// Add each output format as a child entry
	for i, output := range outputs {
		// Pad index for ordering: 01_, 02_, ...
		itemName := fmt.Sprintf("%02d_%s", i+1, output.Ext)
		itemPath := rootPath + `\shell\` + itemName

		itemKey, _, err := registry.CreateKey(
			registry.CLASSES_ROOT,
			itemPath,
			registry.SET_VALUE|registry.CREATE_SUB_KEY,
		)
		if err != nil {
			return err
		}
		itemKey.SetStringValue("", output.Label)
		itemKey.Close()

		// Create the command subkey
		cmdKey, _, err := registry.CreateKey(
			registry.CLASSES_ROOT,
			itemPath+`\command`,
			registry.SET_VALUE,
		)
		if err != nil {
			return err
		}
		// Command: "convertify.exe" convert "%1" <ext>
		command := fmt.Sprintf(`"%s" convert "%%1" %s`, exePath, output.Ext)
		cmdKey.SetStringValue("", command)
		cmdKey.Close()
	}

	return nil
}

// deleteKeyRecursive deletes a registry key and all its subkeys
func deleteKeyRecursive(root registry.Key, path string) error {
	// Split path into parent and last component
	parts := strings.Split(path, `\`)
	if len(parts) == 0 {
		return nil
	}

	parentPath := strings.Join(parts[:len(parts)-1], `\`)
	childName := parts[len(parts)-1]

	parentKey, err := registry.OpenKey(root, parentPath, registry.SET_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return nil // already gone
	}
	defer parentKey.Close()

	return registry.DeleteKey(parentKey, childName)
}
