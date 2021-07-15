package agent

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/user"
	"strconv"
)

func (a *Agent) Files(files []File) error {
	// HTTP Client
	client := SetupClient()

	// This is now broken.
	for _, file := range files {
		if file.State == "present" {
			fileDiff, err := a.fileDifferent(file)
			if err != nil {
				a.Logger.Printf("Error Comparing File: %v", err)
			}
			if fileDiff {
				a.Logger.Println("Fetchin " + a.ServerUrl + "/files/" + file.Source)
				err := a.downloadFile(file, client)
				if err != nil {
					a.Logger.Printf("Error Downloading File: %v", err)
				}
				// Tell Service to reload
				a.Reload = true
			}

			err = a.fixPerms(file)
			if err != nil {
				a.Logger.Printf("Failed setting permissions: %v", err)
			}
		}
		if file.State == "absent" {
			err := a.removeFile(file)
			if err != nil {
				a.Logger.Printf("Failed to remove file: %v", err)
			}
			// Tell Service to reload
			a.Reload = true
		}
	}
	return nil
}

func (a *Agent) removeFile(f File) error {
	if _, err := os.Stat(f.Path); err == nil {
		err := os.Remove(f.Path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Agent) downloadFile(f File, c *http.Client) error {
	// Where to download file from
	fileUrl := a.ServerUrl + "/files/" + f.Source

	resp, err := http.Get(fileUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Open file for create/truncate
	mode, _ := strconv.ParseUint(f.Mode, 8, 32)
	file, err := os.OpenFile(f.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.FileMode(mode))
	if err != nil {
		return err
	}
	defer file.Close()

	// Write out downloaded file to local
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) fileDifferent(f File) (bool, error) {
	// Check if file exists
	_, err := os.Stat(f.Path)
	if err != nil {
		return true, nil
	}

	// Open file
	file, err := os.Open(f.Path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	h := sha256.New()
	if _, err = io.Copy(h, file); err != nil {
		return false, err
	}
	// Return true if files are the same
	if f.ShaSum == fmt.Sprintf("%x", h.Sum(nil)) {
		return false, nil
	}
	return true, nil
}

func (a *Agent) fixPerms(f File) error {
	// Does file exist?
	_, err := os.Stat(f.Path)
	if err != nil {
		return err
	}

	// Get UID int
	uid, err := user.Lookup(f.Owner)
	if err != nil {
		return err
	}
	// Get GID int
	gid, err := user.LookupGroup(f.Group)
	if err != nil {
		return err
	}

	// Set file Mode
	mode, _ := strconv.ParseUint(f.Mode, 8, 32)
	err = os.Chmod(f.Path, fs.FileMode(mode))
	if err != nil {
		a.Logger.Println("Error chmoding ", f.Path)
		return err
	}

	// Set Owner and Group of File
	u, _ := strconv.Atoi(uid.Uid)
	g, _ := strconv.Atoi(gid.Gid)
	err = os.Chown(f.Path, u, g)
	if err != nil {
		return err
	}
	return nil
}
