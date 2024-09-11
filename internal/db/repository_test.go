package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository_IsExists(t *testing.T) {
	//t.Skip()

	repository, err := DefaultRepository()
	assert.NoError(t, err)
	assert.NotNil(t, repository)

	//exists, err := repository.IsExists(ctx, 1)
	//assert.NoError(t, err)
	//assert.False(t, exists)
}
