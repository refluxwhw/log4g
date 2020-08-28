package logger_mgr

func LoadConfig(path string) {
	loadConfig(path)
}

func GetLogger(category string) Logger {
	return newLoggerWithDefaultSourceGetter(category, 2)
}

func GetLoggerWithSkip(category string, skip int) Logger {
	return newLoggerWithDefaultSourceGetter(category, skip)
}

func GetLoggerWithSourceGetter(category string, source SourceGetter) Logger {
	return newLoggerWithSourceGetter(category, source)
}

func Close()  {
	closeLog4go()
}