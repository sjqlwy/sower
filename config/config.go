package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/google/uuid"
	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/wweir/sower/util"
)

var (
	version, date string
	cfg           = &Cfg{}
	// onRefresh execute hooks while refesh config
	onRefresh = []func(*Cfg) (string, error){
		func(c *Cfg) (action string, err error) {
			action = "clear dns cache"

			if c.Client.Suggest.OnSuggest != "" {
				err = execute(c.Client.Suggest.OnSuggest)
			}
			return
		},
	}
)

// Peer is a p2p peer or server node
type Peer struct {
	AddrUUID  string `toml:"addr_uuid"`
	Transport string `toml:"transport"`
	Cipher    string `toml:"cipher"`
	Password  string `toml:"password"`
}

// IsP2P check if the peer is p2p peer
func (p *Peer) IsP2P() (ok bool, network, addr string) {
	p.AddrUUID = strings.TrimSpace(p.AddrUUID)

	if _, err := uuid.Parse(p.AddrUUID); err == nil {
		return true, "", ""
	}

	secs := strings.Split(p.AddrUUID, "://")
	if len(secs) != 2 {
		glog.Exitf("invalid addr_uuid setting: (%s)", p.AddrUUID)
	}
	if !strings.Contains(secs[1], ":") {
		glog.Exitf("invalid addr_uuid setting: (%s)", p.AddrUUID)
	}

	return false, secs[0], secs[1]
}

// Cfg is the config file definition
type Cfg struct {
	ConfigFile string `toml:"-"`
	LogVerbose int    `toml:"log_verbose"`

	Client struct {
		ClientIP       string `toml:"client_ip"`
		DNSIP          string `toml:"dns_ip"`
		UpstreamDNS    string `toml:"upstream_dns"`
		Socks5Addr string `toml:"upstream_socks5"`

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

	P2P struct {
		Self    Peer   `toml:"self"`
		Peer    Peer   `toml:"peer"`
		Servers []Peer `toml:"servers"`
	} `toml:"p2p"`

	HTTPProxys []struct {
		ListenAddr     string `toml:"listen_addr"`
		Socks5Addr string `toml:"upstream_socks5"`
		Peer           Peer   `toml:"peer"`
	} `toml:"http_proxy"`

	DirectProxys []struct {
		ListenAddr string `toml:"listen_addr"`
				Socks5Addr string `toml:"upstream_socks5"`
		TargetAddr string `toml:"target_addr"`
		Peer       Peer   `toml:"peer"`
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

// store persist config from memory to file
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
