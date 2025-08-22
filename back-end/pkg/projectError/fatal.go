package projectError

import (
	"log"
	"os"
)

func FatalError(err *Error) {
	log.Fatalf("FATAL ERROR: %s - %s - %s - %v", err.Code, err.Message, err.Path, err.PrevError)
	os.Exit(1)
}
