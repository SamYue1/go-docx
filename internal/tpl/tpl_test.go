package tpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeDefaultTemplate(t *testing.T) {
	t.Run("it_opens_default_docx", func(t *testing.T) {
		data, err := OpenDefault()
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
	})

	t.Run("it_is_valid_zip_archive", func(t *testing.T) {
		data, err := OpenDefault()
		assert.NoError(t, err)
		assert.True(t, len(data) > 0)
	})
}
