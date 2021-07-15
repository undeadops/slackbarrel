package agent

import (
	"fmt"
	"os/exec"
)

func (a *Agent) ReloadService(services []Service) error {
	for _, service := range services {
		// First Reload Systemd
		c := exec.Command("systemctl", "daemon-reload")
		_, err := c.CombinedOutput()
		if err != nil {
			return fmt.Errorf("Error reloading systemd: ", err)
		}

		// Reload or Restart the service depending on config
		switch service.Action {
		case "reload":
			c = exec.Command("systemctl", "reload", service.Name)
			_, err = c.CombinedOutput()
			if err != nil {
				return fmt.Errorf("Error reloading service: ", err)
			}
		case "restart":
			c = exec.Command("systemctl", "restart", service.Name)
			_, err = c.CombinedOutput()
			if err != nil {
				return fmt.Errorf("Error restarting service: ", err)
			}
		}
	}
	return nil
}
