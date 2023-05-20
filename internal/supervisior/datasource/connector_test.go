package datasource_test

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kroseida.org/slixx/internal/supervisior/application"
	"kroseida.org/slixx/internal/supervisior/datasource"
	"os"
	"testing"
	"time"
)

func setupSuite() func() {
	os.Remove("slixx.db")

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.Encoding = "console"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	application.Logger = logger.Sugar()
	return func() {
		os.Remove("slixx.db")
	}
}

func Test_Connect(t *testing.T) {
	teardownSuite := setupSuite()
	err := datasource.Connect()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	_, err = os.ReadFile(application.CurrentSettings.Database.Configuration["file"])
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	assert.NotNil(t, datasource.StorageProvider)
	assert.NotNil(t, datasource.UserProvider)
	err = datasource.Close()
	if err != nil {
		t.Error(err)
		teardownSuite()
		return
	}
	teardownSuite()
}
