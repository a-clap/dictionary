package dict

import "github.com/a-clap/dictionary/internal/logger"

var log logger.Logger = logger.NewDummy()

func InitLogger(logger logger.Logger) {
	log = logger
}
