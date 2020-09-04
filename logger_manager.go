package log4g

var (
	gLoggerMgr     = map[string]*category{}
	gDefaultLogger = newDefaultCategory("console")

	Critical       = gDefaultLogger.Critical
	CriticalF      = gDefaultLogger.CriticalF
	Error          = gDefaultLogger.Error
	ErrorF         = gDefaultLogger.ErrorF
	Warn           = gDefaultLogger.Warn
	WarnF          = gDefaultLogger.WarnF
	Info           = gDefaultLogger.Info
	InfoF          = gDefaultLogger.InfoF
	Debug          = gDefaultLogger.Debug
	DebugF         = gDefaultLogger.DebugF
	Trace          = gDefaultLogger.Trace
	TraceF         = gDefaultLogger.TraceF
	Log            = gDefaultLogger.Log
	LogF           = gDefaultLogger.LogF
)

func SetDefaultLogger(name string) {
	gDefaultLogger = GetLogger(name)

	Critical = gDefaultLogger.Critical
	CriticalF = gDefaultLogger.CriticalF
	Error = gDefaultLogger.Error
	ErrorF = gDefaultLogger.ErrorF
	Warn = gDefaultLogger.Warn
	WarnF = gDefaultLogger.WarnF
	Info = gDefaultLogger.Info
	InfoF = gDefaultLogger.InfoF
	Debug = gDefaultLogger.Debug
	DebugF = gDefaultLogger.DebugF
	Trace = gDefaultLogger.Trace
	TraceF = gDefaultLogger.TraceF
	Log = gDefaultLogger.Log
	LogF = gDefaultLogger.LogF
}

func LoadYamlFile(path string) error {
	return loadYamlFile(path)
}

func LoadJsonFile(path string) error {
	return loadJsonFile(path)
}

func LoadJsonString(js string) error {
	return loadJsonString(js)
}

// 关闭所有的文件
func Close() {
	for _, c := range gLoggerMgr {
		for _, f := range c.filters {
			for _, w := range f.writers {
				w.Close()
			}
		}
	}
	gLoggerMgr = map[string]*category{}
}

func GetLogger(name string) Logger {
	lg, ok := gLoggerMgr[name]
	if ok {
		return lg
	}
	return newDefaultCategory(name)
}
