package scrolls

import (
	"fmt"
	"net/http"

	"github.com/mreliasen/scrolls-cli/internal/utils"
)

type VersionClient client

type VersionInfo struct {
	Version string `json:"version"`
}

type VersionInfoResponse struct {
	Version VersionInfo `json:"version"`
}

func (u *VersionClient) GetLatestRelease() (VersionInfo, error) {
	res, err := u.client.Get("/releases/latest", nil)
	if err != nil {
		return VersionInfo{}, fmt.Errorf("failed to get release version info: %s", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return VersionInfo{}, fmt.Errorf("failed to get release version info: %s", err)
	}

	data, err := utils.Unmarshal[VersionInfo](res)
	if err != nil {
		return VersionInfo{}, fmt.Errorf("failed to deserialise response: %w", err)
	}

	return data, nil
}

func (u *VersionClient) CheckForUpdates() {
	/* if u.client.settings.GetAutoupdate() == "on" && time.Now().Unix() >= u.client.settings.GetLastUpdateCheck()+int64(24*60*60) {
		latest, err := fetchLatestVersion()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error fetching latest version:", err)
			return
		}

		parsedVersion, err := semver.NewVersion(GetCurrentVersion())
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing current version:", err)
			return
		}

		parsedLatest, err := semver.NewVersion(latest)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error parsing latest version:", err)
			return
		}

		if parsedVersion.LessThan(parsedLatest) {
			fmt.Println("Updating to the latest version")

			err := Update()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Error updating:", err)
			}

			fmt.Printf("\nYou can disable automatic updates with %s\n", "scrolls config set autoupdate off")
		}
		u.client.settings.SetLastUpdateCheck(time.Now().Unix())
		u.client.settings.PersistChanges()
	} */
}
