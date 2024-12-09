package scrolls

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

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

func (u *VersionClient) CheckForUpdates() {
	latest, err := u.getLatestRelease()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error fetching latest version:", err)
		return
	}

	parsedVersion, err := semver.NewVersion(utils.Version)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing current version:", err)
		return
	}

	parsedLatest, err := semver.NewVersion(latest.Version)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error parsing latest version:", err)
		return
	}

	if parsedVersion.LessThan(parsedLatest) {
		fmt.Println("Updating to the latest version")

		err := u.Update()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error updating:", err)
		}

		return
	}

	fmt.Printf("version %s is already latest\n", utils.Version)
}
