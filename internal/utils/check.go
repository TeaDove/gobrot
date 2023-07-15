package utils

import "github.com/rs/zerolog/log"

func Check(err error) {
	if err != nil {
		log.Panic().Stack().Err(err).Str("status", "check.failed").Send()
	}
}
