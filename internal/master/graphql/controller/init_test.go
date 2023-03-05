package controller_test

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"kroseida.org/slixx/internal/master/application"
	"kroseida.org/slixx/internal/master/datasource"
	"kroseida.org/slixx/pkg/model"
	"os"
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
	datasource.Connect()

	return func() {
		datasource.Close()
		os.Remove("slixx.db")
	}
}

func withPermissions(permissions []string) context.Context {
	user := &model.User{
		Id:          uuid.New(),
		Name:        "Test",
		FirstName:   "Test",
		LastName:    "Test",
		Active:      true,
		Email:       "",
		Description: "",
		Permissions: "[]",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	user.AddPermission(permissions)

	return context.WithValue(context.Background(), "user", user)
}
