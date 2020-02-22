package reproksie

import (
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	os.Chdir("..")
	t.Run("Run.Config.JSON", func(t *testing.T) {
		args := []string{
			"",
			"-c=example/config.json",
			"-b=true",
		}

		app := NewApp(AppConfig{})
		err := app.Run(args)

		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Run.Config.YAML", func(t *testing.T) {
		args := []string{
			"",
			"-c=example/config.yml",
			"-b=true",
		}

		app := NewApp(AppConfig{})
		err := app.Run(args)

		if err != nil {
			t.Error(err)
		}
	})
	t.Run("Run.Clean", func(t *testing.T) {
		app := NewApp(AppConfig{})
		err := app.Run([]string{""})

		if err == nil {
			t.Error("Expected error but got nil")
		}
	})
	t.Run("Run.NonExist", func(t *testing.T) {
		app := NewApp(AppConfig{})
		err := app.Run([]string{"", "-c=somenotexistingdir/noneexistingconfig.json"})

		if err == nil {
			t.Error("Expected error but got nil")
		}
	})
}
