package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTmpl(t *testing.T) {
	err := run("main", []string{"ChatID", "MessageID"}, os.Stdout)
	require.NoError(t, err)
}
