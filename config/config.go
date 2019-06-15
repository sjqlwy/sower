package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/golang/glog"
	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/wweir/sower/util"
)

var (
	version, date string
	cfg           = &Cfg{}
	// onRefresh execute hooks while refesh config
	onRefresh = []func(*Cfg) (string, error){
		func(c *Cfg) (string, error) {
			if c.Client.Suggest.OnSuggest == "" {
				return "", nil
			}
			if err := execute(c.Client.Suggest.OnSuggest); err != nil {
				return "", err
			}
			return "flushed dns cache", nil
		},
	}
)

// Cfg is the config file definition
type Cfg struct {
	ConfigFile string `toml:"-"`
	LogVerbose int    `toml:"log_verbose"`

	Client struct {
		ServeIP string `toml:"serve_ip"`
		DNS     string `toml:"dns"`

		Rule struct {
			BlockList []string `toml:"blocklist"`
			WhiteList []string `toml:"whitelist"`
		} `toml:"rule"`

		Suggest struct {
			SuggestLevel string   `toml:"suggest_level"`
			OnSuggest    string   `toml:"on_suggest"`
			Suggestions  []string `toml:"suggestions"`
		} `toml:"suggest"`
	} `toml:"client"`

	Transport struct {
		SelfURI   string   `toml:"self_uri"`
		OutletURI string   `toml:"outlet_uri"`
		Brokers   []string `toml:"brokers"`
	} `toml:"transport"`

	HTTPProxys []struct {
		ListenAddr string `toml:"listen_addr"`
		OutletURI  string `toml:"outlet_uri"`
	} `toml:"http_proxy"`

	DirectProxys []struct {
		ListenAddr string `toml:"listen_addr"`
		TargetAddr string `toml:"target_addr"`
		OutletURI  string `toml:"outlet_uri"`
	} `toml:"direct_proxy"`

	mu sync.Mutex
}

// GetCfg return the default config
func GetCfg() *Cfg {
	return cfg
}

// Init initialize the config from config file
func (c *Cfg) Init() error {
	f, err := os.OpenFile(c.ConfigFile, os.O_RDONLY, 0644)
	if err != nil {
		return errors.Wrapf(err, "load config (%s)", c.ConfigFile)
	}
	defer f.Close()

	file := c.ConfigFile // keep config file path
	if err = toml.NewDecoder(f).Decode(c); err != nil {
		return errors.Wrapf(err, "decode config (%s) fail: %s", c.ConfigFile, err)
	}
	c.ConfigFile = file

	for i := range onRefresh {
		if action, err := onRefresh[i](c); err != nil {
			return errors.Wrap(err, action)
		}
	}

	return nil
}

//AddHook add hook function at refresh point
func (c *Cfg) AddHook(fn func(*Cfg) (string, error), initRun bool) (string, error) {
	c.mu.Lock()
	onRefresh = append(onRefresh, fn)
	c.mu.Unlock()
	if initRun {
		return fn(c)
	}
	return "", nil
}

// AddSuggestion add new domain into suggest rules
func (c *Cfg) AddSuggestion(domain string) {
	c.mu.Lock()
	c.Client.Suggest.Suggestions = append(c.Client.Suggest.Suggestions, domain)
	c.Client.Suggest.Suggestions = util.NewReverseSecSlice(c.Client.Suggest.Suggestions).Sort().Uniq()
	c.mu.Unlock()

	if err := c.store(); err != nil {
		glog.Errorln(err)
	}

	// reload config
	for i := range onRefresh {
		if action, err := onRefresh[i](c); err != nil {
			glog.Errorln(action+":", err)
		}
	}
}

// store safely persist config from memory to file
func (c *Cfg) store() error {
	f, err := os.OpenFile(c.ConfigFile+"~", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "open %s~", c.ConfigFile)
	}
	defer func() {
		f.Close()
		os.Remove(c.ConfigFile + "~")
	}()

	if err := toml.NewEncoder(f).ArraysWithOneElementPerLine(true).Encode(c); err != nil {
		return errors.Wrap(err, "encode config")
	}

	err = os.Rename(c.ConfigFile+"~", c.ConfigFile)
	return errors.Wrapf(err, "move config %s", c.ConfigFile)
}

// printVersion print config and version info
func (c *Cfg) printVersion() {
	config, _ := json.MarshalIndent(cfg, "", "\t")
	fmt.Printf("Version:\n\t%s %s\nCfgig:\n%s", version, date, config)
	os.Exit(0)
}
