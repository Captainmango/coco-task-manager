package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ItGetsCrontabFile(t *testing.T) {
	t.Parallel()
	os.Setenv("CRONTAB_FILE", "test.value")

	BootstrapConfig()

	assert.Equal(t, "test.value", Config.CrontabFile)
}

func Test_ItSuppliesDefaultCronTab(t *testing.T) {
	t.Parallel()
	BootstrapConfig()

	assert.NotEqual(t, "", Config.CrontabFile)
}
