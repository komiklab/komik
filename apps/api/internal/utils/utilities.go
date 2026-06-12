package utils

import (
	"runtime"

	"github.com/rs/zerolog/log"
)

func IsErrNotNil(err error) bool {
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			log.Error().Err(err).Str("file", file).Int("line", line).Msg("Error occured")
		} else {
			log.Error().Err(err).Msg("Error occurred but caller information is not available")
		}
		return true
	}
	return false
}
