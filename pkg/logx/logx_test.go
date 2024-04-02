package logx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefault(t *testing.T) {
	logger, err := New()
	assert.Nil(t, err)
	assert.NotNil(t, logger)

	logger.Slog().Info("default logger")
}

func TestFileLogger(t *testing.T) {
	logger, err := New(WithOutput("testdata/test.log"))
	defer logger.Close()

	assert.Nil(t, err)
	assert.NotNil(t, logger)

	logger.Slog().Info("file logger")
}

func TestJSONLogger(t *testing.T) {
	logger, err := New(WithOutput("testdata/test.log"), WithFormat(JSONFormat))
	defer logger.Close()

	assert.Nil(t, err)
	assert.NotNil(t, logger)

	logger.Slog().Info("file logger")
}
