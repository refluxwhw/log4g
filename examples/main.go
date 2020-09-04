package main

import (
	"time"

	log "log4g"
)

func main() {
	testLog()
}

func testLog()  {
	// _ = log.LoadJsonFile("example.json")
	_ = log.LoadYamlFile("example.yaml")
	t := log.GetLogger("TestA")
	t.Debug("test debug")
	t.Info("test info")
	t.Error("test error")
	t.Critical("test critical")

	go func() {
		tb := log.GetLogger("TestB")
		tb.Debug("test debug")
		tb.Info("test info")
		tb.Error("test error")
		tb.Critical("test critical")
	}()

	go func() {
		tb := log.GetLogger("TestB")
		tb.Debug("test debug")
		tb.Info("test info")
		tb.Error("test error")
		tb.Critical("test critical")
	}()

	log.Info("console info")
	time.Sleep(time.Second * 1)
	log.Close()
	time.Sleep(time.Second * 1)
}