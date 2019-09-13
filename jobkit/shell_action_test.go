package jobkit

import (
	"fmt"
	"os"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func TestShellAction(t *testing.T) {
	assert := assert.New(t)

	envVar := uuid.V4().String()
	envValue := uuid.V4().String()
	os.Setenv(envVar, envValue)
	defer os.Unsetenv(envVar)

	action := ShellAction([]string{"sh", fmt.Sprintf("$%s/foo", envVar)})
	assert.NotNil(action)
}
