package main

var version, date string

func main() {
	// Conf := &conf.Conf
	// if Conf.VersionOnly {
	// 	config, _ := json.MarshalIndent(Conf, "", "\t")
	// 	fmt.Printf("Version:\n\t%s %s\nConfig:\n%s", version, date, config)
	// 	return
	// }
	// glog.Infof("Starting sower(%s %s): %v", version, date, Conf)

	// tran, err := transport.GetTransport(Conf.NetType)
	// if err != nil {
	// 	glog.Exitln(err)
	// }

	// if Conf.ServerAddr == "" {
	// 	proxy.StartServer(tran, Conf.ServerPort, Conf.Cipher, Conf.Password)

	// } else {
	// 	conf.AddRefreshFn(true, func() (string, error) {
	// 		dns.LoadRules(Conf.BlockList, Conf.Suggestions, Conf.WhiteList, Conf.ServerAddr)
	// 		return "load rules", nil
	// 	})

	// 	isSocks5 := (Conf.NetType == "SOCKS5")
	// 	serverAddr := net.JoinHostPort(Conf.ServerAddr, Conf.ServerPort)

	// 	if Conf.HTTPProxy != "" {
	// 		go proxy.StartHttpProxy(tran, isSocks5, serverAddr, Conf.Cipher, Conf.Password, Conf.HTTPProxy)
	// 	}

	// 	go dns.StartDNS(Conf.DNSServer, Conf.ClientIP, conf.SuggestCh, Conf.SuggestLevel)
	// 	proxy.StartClient(tran, isSocks5, serverAddr, Conf.Cipher, Conf.Password, Conf.ClientIP)
	// }
}
