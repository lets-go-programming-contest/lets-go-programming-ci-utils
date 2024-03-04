package ci

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"path"
	"strconv"
	"testing"
)

const TESTS_COUNT = 8

func TestMyProgram(t *testing.T) {
	pathToCommon := path.Join(os.Getenv("STATIC_DIR"), "cases")
	for i := 1; i <= TESTS_COUNT; i++ {
		t.Run(fmt.Sprintf("test-2-2-case-%d", i), func(t *testing.T) {
			assert := assert.New(t)
			fileInput, err := os.Open(path.Join(pathToCommon, strconv.Itoa(i)))
			if err != nil {
				assert.FailNow(err.Error())
			}
			defer fileInput.Close()
			cmd := exec.Command(path.Join(os.Getenv("BUILD_BIN"), "service"))
			cmd.Stdin = fileInput
			output, err := cmd.CombinedOutput()
			assert.NoError(err, "Failed to execute command: %s", cmd)
			fileOutput, err := os.ReadFile(path.Join(pathToCommon, fmt.Sprintf("%d.a", i)))
			if err != nil {
				assert.FailNow(err.Error())
			}
			assert.Equal(string(fileOutput), string(output))
		})
	}
}
