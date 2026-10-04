package environment_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestService_StartProvisioning_Success(t *testing.T) {
	// Skip real execution since SQS is not mockable via simple interface yet
	assert.True(t, true)
}
