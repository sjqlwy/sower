package main

var version, date string

func main() {
	// cfg := &conf.Conf
	// if cfg.VersionOnly {
	// 	config, _ := json.MarshalIndent(cfg, "", "\t")
	// 	fmt.Printf("Version:\n\t%s %s\nConfig:\n%s", version, date, config)
	// 	return
	// }
	// glog.Infof("Starting sower(%s %s): %v", version, date, cfg)

	// tran, err := transport.GetTransport(cfg.NetType)
	// if err != nil {
	// 	glog.Exitln(err)
	// }

	// if cfg.ServerAddr == "" {
	// 	proxy.StartServer(tran, cfg.ServerPort, cfg.Cipher, cfg.Password)

	// } else {
	// 	conf.AddRefreshFn(true, func() (string, error) {
	// 		dns.LoadRules(cfg.BlockList, cfg.Suggestions, cfg.WhiteList, cfg.ServerAddr)
	// 		return "load rules", nil
	// 	})

	// 	isSocks5 := (cfg.NetType == "SOCKS5")
	// 	serverAddr := net.JoinHostPort(cfg.ServerAddr, cfg.ServerPort)

	// 	if cfg.HTTPProxy != "" {
	// 		go proxy.StartHttpProxy(tran, isSocks5, serverAddr, cfg.Cipher, cfg.Password, cfg.HTTPProxy)
	// 	}

	// 	go dns.StartDNS(cfg.DNSServer, cfg.ClientIP, conf.SuggestCh, cfg.SuggestLevel)
	// 	proxy.StartClient(tran, isSocks5, serverAddr, cfg.Cipher, cfg.Password, cfg.ClientIP)
	// }
}
