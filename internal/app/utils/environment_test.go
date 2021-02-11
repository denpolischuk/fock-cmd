package utils_test

import (
	"os"
	"testing"

	"github.com/denpolischuk/fock-cli/internal/app/utils"
)

func TestGetUserShell(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		os.Setenv("SHELL", "/usr/bin/test")
		sh, _ := utils.GetUserShell()

		if sh != "test" {
			t.Fail()
		}
	})
	t.Run("fail", func(t *testing.T) {
		os.Setenv("SHELL", "")
		_, err := utils.GetUserShell()

		if err == nil {
			t.Fail()
		}
	})
}

func TestCheckIfAppInstalled(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := utils.CheckIfAppInstalled("bash")

		if !r {
			t.Fail()
		}
	})
	t.Run("fail", func(t *testing.T) {
		r := utils.CheckIfAppInstalled("somenonexistingbinarytotest")

		if r {
			t.Fail()
		}
	})
}
