package tui

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenBrowser opens the given URL in the user's default browser. It is used by
// the setup wizard to launch the bank authorization / app-registration pages.
func OpenBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}
