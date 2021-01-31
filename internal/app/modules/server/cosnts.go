package server

import "github.com/denpolischuk/fock-cli/internal/app/consts"

const (
	logOutputFilePath = "/tmp/fock_server_output"
)

var notRunningMessage = consts.Emojis["fail"] + " Fock server is not running"
