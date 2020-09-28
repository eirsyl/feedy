package open

import (
	"os/exec"
	"runtime"
)

// OpenURL opens an url in the configured browser
func OpenURL(url string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Start() // nolint: errcheck, gas
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start() // nolint: errcheck, gas
	case "darwin":
		exec.Command("open", url).Start() // nolint: errcheck, gas
	}
}
