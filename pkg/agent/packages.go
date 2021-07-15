package agent

import (
	"fmt"
	"os"
	"os/exec"
)

func (a *Agent) Packages(pkgs []Package) error {
	for _, p := range pkgs {
		switch p.State {
		case "present":
			err := a.InstallPkg(p.Name)
			if err != nil {
				fmt.Errorf("Error installing package: ", err)
			}
		case "absent":
			err := a.RemovePkg(p.Name)
			if err != nil {
				fmt.Errorf("Error removing package: ", err)
			}
		}
	}
	return nil
}

func (a *Agent) RemovePkg(name string) error {
	c := exec.Command("apt-get", "remove", "-qq", name)
	c.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
	_, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error in installation: ", err)
	}
	c = exec.Command("apt-get", "autoremove", "-y", "-qq")
	c.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
	_, err = c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error removing unneeded packages", err)
	}
	return nil
}

func (a *Agent) InstallPkg(name string) error {
	c := exec.Command("apt-get", "update", "-qq")
	c.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
	_, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error Updating Apt package cache", err)
	}
	c = exec.Command("apt-get", "install", "-y", "--no-install-recommends", "-qq", name)
	c.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
	_, err = c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error Fetching combined output of package install", err)
	}
	return nil
}
