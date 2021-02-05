package config

import (
	"fmt"
	"log"
)

// ReadErrorDefaultHandler - default error handler for reading config error
func ReadErrorDefaultHandler(err error) {
	log.Fatal(fmt.Sprintf("error occured while reading config:\n%s", err.Error()))
}

// WriteErrorDefaultHandler - default error handler for wrtiting into config error
func WriteErrorDefaultHandler(err error) {
	log.Fatal(fmt.Sprintf("error occured while writing to config:\n%s", err.Error()))
}
