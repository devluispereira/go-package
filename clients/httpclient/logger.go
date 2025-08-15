package httpclient

import (
	"os"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func init() {
	logger = zerolog.New(os.Stdout).
		With().Str("layer", "http-client").Logger()
}
