package output

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// OpenBrowser opens a file in the default browser
func OpenBrowser(filename string) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("falha ao obter caminho absoluto: %w", err)
	}

	url := "file://" + absPath

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		return fmt.Errorf("plataforma não suportada: %s", runtime.GOOS)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("erro ao abrir navegador: %w", err)
	}

	return nil
}

// SaveHTML saves HTML content to a file
// filePath should be the full path including directory and filename
// The directory will be created if it doesn't exist
func SaveHTML(filePath string, content string) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	return os.WriteFile(filePath, []byte(content), 0644)
}
