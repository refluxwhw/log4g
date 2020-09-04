package log4g

import (
	"encoding/json"
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/toolkits/file"
)

// 全局段配置
type globalSecCfg struct {
	ConsoleEnable bool `json:"console_enable" yaml:"console_enable"`
}

// 文件段配置
type fileSecCfg struct {
	Filename string `json:"filename" yaml:"filename"`
	Rotate   bool   `json:"rotate" yaml:"rotate"`   // default is true
	Maxsize  string `json:"maxsize" yaml:"maxsize"` // \d+[KMG]? Suffixes are in terms of 2**10, default is "10M"
	Maxline  string `json:"maxline" yaml:"maxline"` // \d+[KMG]? Suffixes are in terms of thousands, default is "100K"
	Daily    bool   `json:"daily" yaml:"daily"`     // Automatically rotates by day, default is true
}

// 日志分类段配置
type categorySecCfg struct {
	Enable  bool                `json:"enable" yaml:"enable"` // default is false
	Filters []categoryFilterCfg `json:"filters" yaml:"filters"`
}

// 日志分类下的过滤器配置
type categoryFilterCfg struct {
	Level  string   `json:"level" yaml:"level"`   // default is DEBUG
	Layout string   `json:"layout" yaml:"layout"` // layout name
	Output []string `json:"output" yaml:"output"` // output to console and files
}

// 完整配置
type fullConfig struct {
	Global globalSecCfg          `json:"global" yaml:"global"`
	Files  map[string]fileSecCfg `json:"files" yaml:"files"`

	/// TODO: 每个分段后面可以通过 {} 增加格式参数

	// Known format codes:
	// %T - DateTime with format string, default format is (2006-01-03 15:04:05.000)
	//  	eg. %T{2006-01-02 15:04:05}
	// %L - Level (DEBG, TRAC, WARN, EROR, CRIT)
	// %l - Level (D, T, W, E, C)
	// %C - category
	// %S - Source
	// %M - Message
	// Ignores unknown formats
	// Recommended: "[%T] %level %C (%S) %M"
	// 		if config is invalid, will be replaced with default value
	Layouts    map[string]string         `json:"layouts" yaml:"layouts"` // default is "[%T] %level %C (%S) %M"
	Categories map[string]categorySecCfg `json:"categories" yaml:"categories"`
}


func loadFullCfg(cfg *fullConfig) error {
	layouts := map[string]*layoutInfo{}
	writers := map[string]logWriter{}

	writers["console"] = gSingleConsoleWriter

	// 创建分类，以及每个分类下的过滤器
	for name, cateCfg := range cfg.Categories {
		if !cateCfg.Enable || len(cateCfg.Filters) == 0 {
			continue
		}

		c, ok := gLoggerMgr[name]
		if !ok {
			c = &category{
				category: name,
				extSkip:  0,
				filters:  make([]*categoryFilter, 0),
			}
			gLoggerMgr[name] = c
		}

		for _, filterCfg := range cateCfg.Filters {
			layout, ok := layouts[filterCfg.Layout]
			if !ok {
				if layoutCfg, ok := cfg.Layouts[filterCfg.Layout]; !ok {
					return fmt.Errorf("layout not found in layouts config")
				} else {
					layout = newLayoutConf(layoutCfg)
					layouts[filterCfg.Layout] = layout
				}
			}

			filter := newFilter(c.category, strToLevel(filterCfg.Level), layout)
			for _, output := range filterCfg.Output {
				writer, ok := writers[output]
				if !ok {
					if fileCfg, ok := cfg.Files[output]; !ok {
						return fmt.Errorf("output not found in files config")
					} else {
						writer = newFileLogWriter(fileCfg)
						writers[output] = writer
					}
				}

				filter.writers = append(filter.writers, writer)
			}

			c.filters = append(c.filters, filter)
		}
	}

	return nil
}

//// -----------------------------------------------------------------------------------

func readFile(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("[%s] path empty", path)
	}

	if !file.IsExist(path) {
		return nil, fmt.Errorf("config file %s is nonexistent", path)
	}

	content, err := file.ToBytes(path)
	if err != nil {
		return nil, fmt.Errorf("read file %s fail %s", path, err)
	}

	return content, nil
}

func parseJson(content []byte) (cfg *fullConfig, err error) {
	cfg = &fullConfig{}
	err = json.Unmarshal(content, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func parseYaml(content []byte) (cfg *fullConfig, err error) {
	cfg = &fullConfig{}
	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

//// -----------------------------------------------------------------------------------

func loadYamlFile(path string) error {
	content, err := readFile(path)
	if err != nil {
		return err
	}

	categories, err := parseYaml(content)
	if err != nil {
		return err
	}

	return loadFullCfg(categories)
}

func loadJsonFile(path string) error {
	content, err := readFile(path)
	if err != nil {
		return err
	}

	cfg, err := parseJson(content)
	if err != nil {
		return err
	}

	return loadFullCfg(cfg)
}

func loadJsonString(content string) error {
	cfg, err := parseJson([]byte(content))
	if err != nil {
		return err
	}

	return loadFullCfg(cfg)
}
