package settings

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"

	"github.com/kirsle/configdir"
	"github.com/mreliasen/scrolls-cli/internal/flags"
	"github.com/spf13/viper"
)

var (
	settings *Settings
	mu       sync.Mutex
)

type Settings struct {
	changed bool
}

func GetConfigDir() (string, error) {
	configPath := configdir.LocalConfig("scrolls")

	configpPathFlag := viper.GetString("config-path")
	if len(configpPathFlag) > 0 {
		configPath = configpPathFlag
	}

	err := configdir.MakePath(configPath)
	if err != nil {
		return "", err
	}

	return configPath, nil
}

func LoadSettings() (*Settings, error) {
	mu.Lock()
	defer mu.Unlock()

	if settings != nil {
		return settings, nil
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	settings = &Settings{}
	viper.BindEnv("config-path", "SCROLLS_CLI_CONFIG_DIR")
	viper.SetConfigName("settings")
	viper.SetConfigType("json")
	viper.AddConfigPath(configDir)
	configFile := path.Join(configDir, "settings.json")

	if abs, err := filepath.Abs(configFile); err == nil {
		configFile = abs
	}

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// Force config creation
			if err := viper.SafeWriteConfig(); err != nil {
				return nil, err
			}
		case viper.ConfigParseError:
			if flags.ResetConfig() {
				viper.WriteConfig()
				break
			}

			fmt.Printf("Warning: could not parse JSON config from file %s\n", configFile)
			fmt.Printf("Fix the syntax errors on the file, or use the --reset-config flag to replace it with a fresh one.\n")
			fmt.Printf("E.g. scrolls config init --reset-config\n")

			return nil, err
		default:
			return nil, err
		}
	}

	return settings, nil
}

func (s *Settings) PersistChanges() {
	if settings == nil || !settings.changed {
		return
	}

	if err := viper.WriteConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to persist scrolls-cli settings file: %s\n", err.Error())
	}
}

func (s *Settings) GetEditor() string {
	e := viper.GetString("editor")
	if e != "" {
		return e
	}

	// defaults
	_, err := exec.LookPath("vim")
	if err == nil {
		return "vim"
	}

	_, err = exec.LookPath("vi")
	if err == nil {
		return "vi"
	}

	_, err = exec.LookPath("zed")
	if err == nil {
		return "zed"
	}

	_, err = exec.LookPath("notepad")
	if err == nil {
		return "notepad"
	}

	return ""
}

func (s *Settings) SetEditor(editor string) error {
	viper.Set("editor", editor)
	s.changed = true
	return nil
}

func (s *Settings) SetAutoupdate(auto_update string) {
	viper.Set("auto_update", auto_update)
	s.changed = true
}

func (s *Settings) SetLastUpdateCheck(t int64) {
	viper.Set("last_update_check", t)
	s.changed = true
}

func (s *Settings) GetLastUpdateCheck() int64 {
	return viper.GetInt64("last_update_check")
}

func (s *Settings) GetAutoupdate() bool {
	return viper.GetBool("auto_update")
}

func (s *Settings) GetLibrary() string {
	lib := viper.GetString("library")
	if lib != "" {
		return lib
	}

	configDir, err := GetConfigDir()
	if err != nil {
		panic("failed to get configuration path")
	}

	return configDir
}

func (s *Settings) SetLibrary(path string) {
	viper.Set("library", path)
	s.changed = true
}

func (s *Settings) SetMigrationVersion(v string) {
	viper.Set("migration_version", v)
	s.changed = true
}

func (s *Settings) GetMigrationVersion() string {
	return viper.GetString("migration_version")
}
