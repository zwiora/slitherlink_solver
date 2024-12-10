package debug

import (
	"fmt"
	"log"
	"time"
)

var IsDebugMode bool

func Println(message any) {
	if IsDebugMode {
		if message != "" {
			log.Println(message)
		} else {
			fmt.Println()
		}
	}
}

func Sleep(t int) {
	if IsDebugMode {
		time.Sleep(time.Duration(t) * time.Millisecond)
	}
}
