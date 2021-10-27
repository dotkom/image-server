package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_rootCmd(t *testing.T) {
	if err := Execute(); err != nil {
		assert.FailNowf(t, "Failed to execute root command without any flags.", "Error: %v", err)
	}
}
