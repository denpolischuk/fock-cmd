package utils_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/denpolischuk/fock-cli/internal/app/utils"
)

func TestReplaceInFile(t *testing.T) {
	testStr := fmt.Sprintf("some testing\nfile content")
	t.Run("success with regexp mode OFF", func(t *testing.T) {
		var buf bytes.Buffer
		buf.WriteString(testStr)
		res, err := utils.ReplaceInFile(&buf, "testing", "test", false)
		if err != nil {
			t.Fail()
		}
		if res != fmt.Sprintf("some test\nfile content\n") {
			t.Fail()
		}
	})
	t.Run("success with regexp mode ON", func(t *testing.T) {
		var buf bytes.Buffer
		buf.WriteString(testStr)
		res, err := utils.ReplaceInFile(&buf, `tes\w+`, "test", true)
		if err != nil {
			t.Fail()
		}
		if res != "some test\nfile content\n" {
			t.Fail()
		}
	})

}
