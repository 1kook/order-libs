package cex_test

import (
	"os"
	"testing"

	"github.com/goonma/order-libs/cex"
	"github.com/goonma/order-libs/defines"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestNewCexClient(t *testing.T) {
	err := godotenv.Load(os.ExpandEnv("/config/.env"))
	if err != nil {
		err := godotenv.Load(os.ExpandEnv(".env"))
		if err != nil {
			panic(err)
		}
	}

	opts := []cex.ClientOption{
		cex.WithApiKeyAndSecretKey(os.Getenv("BIN_API_KEY"), os.Getenv("BIN_API_SECRET")),
	}

	_, err = cex.NewCEXClient(defines.CEXCodeBinance, opts...)
	assert.NoError(t, err)

	_, err = cex.NewCEXClient(defines.CEXCodeBinance)
	assert.Error(t, err)
}
