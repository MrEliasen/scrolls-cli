package scrolls

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	semver "github.com/hashicorp/go-version"
	"github.com/mreliasen/scrolls-cli/internal/utils"
)

type VersionClient client

type VersionInfo struct {
	Version string `json:"latest"`
}

type VersionInfoResponse struct {
	Version VersionInfo `json:"latest"`
}

func (u *VersionClient) getLatestRelease() (VersionInfo, error) {
	res, err := u.client.Get("/releases/latest.json", nil)
	if err != nil {
		return VersionInfo{}, fmt.Errorf("failed to get release version info: %s", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return VersionInfo{}, fmt.Errorf("failed to get release version info: %s", err)
	}

	data, err := utils.UnmarshalResp[VersionInfo](res)
	if err != nil {
		return VersionInfo{}, fmt.Errorf("failed to deserialise response: %w", err)
	}

	return data, nil
}

func (u *VersionClient) Update() error {
	command := exec.Command("sh", "-c", "curl -sSfL \"https://cdn.scrolls.sh/releases/install.sh\" | sh")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		return fmt.Errorf("failed to execute update command: %w", err)
	}

	return nil
}

func (u *VersionClient) CheckForUpdates(autoUpdate bool) (currentVersion, latestVersion string, updateAvailable bool, updateError error) {
	latest, err := u.getLatestRelease()
	if err != nil {
		return "", "", false, fmt.Errorf("Error fetching latest version: %w", err)
	}

	parsedVersion, err := semver.NewVersion(utils.Version)
	if err != nil {
		return "", "", false, fmt.Errorf("Error parsing current version: %w", err)
	}

	parsedLatest, err := semver.NewVersion(latest.Version)
	if err != nil {
		return "", "", false, fmt.Errorf("Error parsing latest version: %w", err)
	}

	currentVersion = parsedVersion.String()
	latestVersion = parsedLatest.String()
	updateAvailable = parsedVersion.LessThan(parsedLatest)

	u.client.Settings.SetLastUpdateCheck(time.Now().Unix())
	u.client.Settings.PersistChanges()

	if autoUpdate {
		if updateAvailable {
			fmt.Println("Updating to the latest version")

			err := u.Update()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Error updating:", err)
			}
		} else {
			fmt.Printf("version %s is already latest\n", utils.Version)
		}
	}

	return currentVersion, latestVersion, updateAvailable, nil
}
