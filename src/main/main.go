package main

import (
	log "logger_mgr"
)

func main() {
	log.LoadConfig("test/example.json")
	t := log.GetLogger("Test")
	t.Debug("hello 123")
	log.Close()
}