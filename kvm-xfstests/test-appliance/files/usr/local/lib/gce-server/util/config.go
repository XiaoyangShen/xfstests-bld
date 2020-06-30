package util

import (
	"regexp"
)

const (
	// GceConfigFile contains GCE related configs
	GceConfigFile = "/usr/local/lib/gce_xfstests.config"
	// KcsConfigFile contains KCS related configs, located on LTM
	KcsConfigFile = "/root/xfstests_bld/kvm-xfstests/.kcs_instance"
)

// Config dictionary retrieved from gce_xfstests.config.
type Config struct {
	kv map[string]string
}

// GetConfig reads from the config file and returns a struct Config.
// It attempts to match each line with two possible config patterns.
func GetConfig(configFile string) Config {
	c := Config{make(map[string]string)}
	var re *regexp.Regexp
	if configFile == GceConfigFile {
		re = regexp.MustCompile(`^declare -- (.*?)="(.*?)"$`)
	} else {
		re = regexp.MustCompile(`^(.*?)=(.*?)$`)
	}

	lines, err := ReadLines(configFile)
	Check(err)

	for _, line := range lines {
		tokens := re.FindStringSubmatch(line)
		if len(tokens) == 3 {
			c.kv[tokens[1]] = tokens[2]
		}
	}

	return c
}

// Get a certain config value according to key.
// Returns empty string if key is not present in config.
func (c *Config) Get(key string) string {
	if val, ok := c.kv[key]; ok {
		return val
	}
	return ""
}
