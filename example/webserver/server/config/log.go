package config

type Log struct {
	Folder string `json:"folder"` // folder path of log file, empty use system log
	Level  string `json:"level"`  // output level: error | warning | info | trace | debug
}
