package notify

import (
	"fmt"
	"os/exec"
	"syscall"
)

// Show displays a Windows system tray balloon notification.
// Non-blocking - fires and forgets.
func Show(title, message string) {
	// Sanitize for PowerShell string embedding
	title = sanitize(title)
	message = sanitize(message)

	script := fmt.Sprintf(`
Add-Type -AssemblyName System.Windows.Forms
$n = New-Object System.Windows.Forms.NotifyIcon
$n.Icon = [System.Drawing.SystemIcons]::Application
$n.BalloonTipIcon = [System.Windows.Forms.ToolTipIcon]::Info
$n.BalloonTipTitle = '%s'
$n.BalloonTipText = '%s'
$n.Visible = $true
$n.ShowBalloonTip(4000)
Start-Sleep -Milliseconds 5000
$n.Dispose()
`, title, message)

	cmd := exec.Command("powershell",
		"-WindowStyle", "Hidden",
		"-NonInteractive",
		"-Command", script,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// Start without waiting
	_ = cmd.Start()
}

// sanitize removes characters that could break PowerShell string literals
func sanitize(s string) string {
	result := ""
	for _, r := range s {
		if r == '\'' || r == '"' || r == '`' {
			continue
		}
		result += string(r)
	}
	return result
}
