package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"syscall"
)

type Updater struct {
	CurrentVersion string
	Repo           string
	LatestRelease  *Release
}

// SetLatestReleaseInfo retrieves the latest release info from the apps github repo
func (u *Updater) SetLatestReleaseInfo() error {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", u.Repo))
	if err != nil {
		return fmt.Errorf("SetLatestReleaseInfo: failed Get request: %w", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("SetLatestReleaseInfo: failed reading response body: %w", err)
	}

	latestRelease := Release{}
	err = json.Unmarshal(b, &latestRelease)
	if err != nil {
		return fmt.Errorf("SetLatestReleaseInfo: failed Unmarshal: %w", err)
	}
	u.LatestRelease = &latestRelease

	return nil
}

// NeedUpdate makes a simple comparison between the current app version and
// the latest version from the app's github repo and determines if an update is needed.
func (u *Updater) NeedUpdate() bool {
	currentVersion, err := strconv.ParseFloat(u.CurrentVersion, 10)
	if err != nil || currentVersion == 0 {
		fmt.Println("NeedUpdate: failed to parse version %w", err)
		return false
	}

	latestVersion, err := strconv.ParseFloat(u.LatestRelease.TagName, 10)
	if err != nil || latestVersion == 0 {
		fmt.Println("NeedUpdate: failed to parse version %w", err)
		return false
	}

	return currentVersion < latestVersion
}

// getDownloadURL filters the available assets from the latest release
// for a binary that matches the current platform
func getDownloadURL(assets []Asset) (string, error) {
	path := fmt.Sprintf("self-update-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		path = path + ".exe"
	}

	for _, asset := range assets {
		if asset.Name == path {
			return asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("getDownloadURL: Could not find binary for current platform")
}

// Update handles downloading the new binary and setting the correct permissions
// before restarting the server
func (u *Updater) Update() error {
	if u.LatestRelease == nil {
		err := u.SetLatestReleaseInfo()
		if err != nil {
			return err
		}
	}

	executable, err := os.Executable()
	if err != nil {
		return err
	}

	// Move existing binary to backup
	executableBackup := executable + ".old"
	if _, err := os.Stat(executable); err == nil {
		os.Rename(executable, executableBackup)
	}

	// I don't have access to a Mac or Windows machine, I suspect this may cause issues
	// with the permissions on Windows vs *nix.
	f, err := os.OpenFile(executable, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		os.Rename(executableBackup, executable)
		return err
	}
	defer f.Close()

	binaryPath, err := getDownloadURL(u.LatestRelease.Assets)
	if err != nil {
		return err
	}

	fmt.Printf("downloading new version to %s\n", executable)
	resp, err := http.Get(binaryPath)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	f.Close()
	os.Remove(executableBackup)

	fmt.Println("Update complete: restarting server")
	if err := syscall.Exec(executable, os.Args, os.Environ()); err != nil {
		return err
	}

	return nil
}
