package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"convertify/formats"
)

// findBinary looks for a binary in:
// 1. Same directory as the running exe
// 2. <exedir>/bin/
// 3. System PATH
func findBinary(name string) (string, error) {
	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		candidates := []string{
			filepath.Join(exeDir, name+".exe"),
			filepath.Join(exeDir, "bin", name+".exe"),
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				return c, nil
			}
		}
	}
	// Fallback to PATH
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}
	return "", fmt.Errorf("%s not found. Download it from https://ffmpeg.org and place it next to convertify.exe", name)
}

// hiddenCmd creates an exec.Cmd with a hidden console window
func hiddenCmd(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

// Convert converts inputPath to the given OutputFormat.
// Output file is placed in the same directory as the input.
// Returns the output file path.
func Convert(inputPath string, fmt formats.OutputFormat) (string, error) {
	ffmpegBin, err := findBinary("ffmpeg")
	if err != nil {
		return "", err
	}

	// Build output path: same dir, same name, new extension
	dir := filepath.Dir(inputPath)
	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	outPath := buildOutputPath(dir, baseName, fmt.Ext)

	// Build args: -y (overwrite) -i <input> [format args] <output>
	args := []string{"-y", "-i", inputPath}
	args = append(args, fmt.Args...)
	args = append(args, outPath)

	cmd := hiddenCmd(ffmpegBin, args...)

	// Capture stderr for error messages
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg falló: %s", lastLines(string(out), 5))
	}

	return outPath, nil
}

// buildOutputPath returns a unique output path, adding (1), (2)... if file exists
func buildOutputPath(dir, baseName, ext string) string {
	candidate := filepath.Join(dir, baseName+"."+ext)
	if _, err := os.Stat(candidate); os.IsNotExist(err) {
		return candidate
	}
	for i := 1; i < 100; i++ {
		candidate = filepath.Join(dir, fmt.Sprintf("%s (%d).%s", baseName, i, ext))
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
	return candidate
}

// lastLines returns the last n lines of a string (for error messages)
func lastLines(s string, n int) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	if len(lines) <= n {
		return s
	}
	return strings.Join(lines[len(lines)-n:], "\n")
}

// Validate checks that ffmpeg and ffprobe are accessible
func Validate() error {
	if _, err := findBinary("ffmpeg"); err != nil {
		return err
	}
	return nil
}
