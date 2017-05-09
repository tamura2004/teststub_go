package main

import (
	"time"
)

func TimeStamp() string {
	return time.Now().Format("20060102150405")
}
