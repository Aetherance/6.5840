package logger

import "log"

const (
	green = "\033[1;32m"
	cyan  = "\033[1;36m"
	red   = "\033[1;31m"
	gold  = "\033[1;93m"
	reset = "\033[0m"
)

func Log_cli(str string) {
	log.Println("Client: " + green + str + reset)
}

func Log_fatal(str string) {
	log.Println(red + str + reset)
}

func Log(str string) {
	log.Println(cyan + str + reset)
}

func Log_import(str string) {
	log.Println()
}