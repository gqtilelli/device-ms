package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnums(t *testing.T) {
	t.Run("success brand type enum", func(t *testing.T) {
		brand := Brand("brand1")
		require.True(t, brand.IsValid())
		brand = Brand("brand3")
		require.True(t, brand.IsValid())
		brand = Brand("blabla")
		require.False(t, brand.IsValid())
	})
}
