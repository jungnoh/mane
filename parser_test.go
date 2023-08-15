package mane

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimple(t *testing.T) {
	dir, err := os.Getwd()
	require.NoError(t, err)
	iface, err := NewParser().ParseFile(dir + "/fixtures/simple.go")
	fmt.Println(iface, err)
}
