package logger

import "log"

var Debug bool = true

const (
	green = "\033[1;32m"
	cyan  = "\033[1;36m"
	red   = "\033[1;31m"
	gold  = "\033[1;93m"
	reset = "\033[0m"
)

func Log_cli(str string) {
	if Debug {
		log.Println(str)
	}
}

func Log_fatal(str string) {
	if Debug {
		log.Println(red + str + reset)
	}
}

func Log(str string) {
	if Debug {
		log.Println(cyan + str + reset)
	}
}

func Log_import(str string) {
	if Debug {
		log.Println(gold + str + reset)
	}
}