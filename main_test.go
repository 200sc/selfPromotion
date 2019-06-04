package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadTemplates(t *testing.T) {
	err := loadTemplates("testdata")
	require.Nil(t, err)
	err = loadTemplates("notadirectory")
	require.NotNil(t, err)
}
