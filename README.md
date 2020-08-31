我是光年实验室高级招聘经理。
我在github上访问了你的开源项目，你的代码超赞。你最近有没有在看工作机会，我们在招软件开发工程师，拉钩和BOSS等招聘网站也发布了相关岗位，有公司和职位的详细信息。
我们公司在杭州，业务主要做流量增长，是很多大型互联网公司的流量顾问。公司弹性工作制，福利齐全，发展潜力大，良好的办公环境和学习氛围。
公司官网是http://www.gnlab.com,公司地址是杭州市西湖区古墩路紫金广场B座，若你感兴趣，欢迎与我联系，
电话是0571-88839161，手机号：18668131388，微信号：echo 'bGhsaGxoMTEyNAo='|base64 -D ,静待佳音。如有打扰，还请见谅，祝生活愉快工作顺利。

# sower
[![GitHub release](http://img.shields.io/github/release/wweir/sower.svg?style=popout)](https://github.com/wweir/sower/releases)
[![Actions Status](https://github.com/wweir/sower/workflows/Go/badge.svg)](https://github.com/wweir/sower/actions)
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/wweir/sower.svg?style=popout)](https://hub.docker.com/r/wweir/sower)
[![GitHub issue](https://img.shields.io/github/issues/wweir/sower.svg?style=popout)](https://github.com/wweir/sower/issues)
[![GitHub star](https://img.shields.io/github/stars/wweir/sower.svg?style=popout)](https://github.com/wweir/sower/stargazers)
[![GitHub license](https://img.shields.io/github/license/wweir/sower.svg?style=popout)](LICENSE)


中文介绍见 [Wiki](https://github.com/wweir/sower/wiki)

The sower is a cross-platform intelligent transparent proxy tool.

The first time you visit a new website, the sower will detect if the domain is accessible and add it in the dynamic detect list. So, you do not need to care about the rules, sower will handle it intelligently.

Sower provider both http_proxy/https_proxy and DNS-based proxy. All these kinds of proxy support intelligent router. You can also port-forward any TCP request to remote, such as SSH / SMTP / POP3.

You can enjoy it by setting http_proxy or your DNS without any other settings.

If you already have another proxy solution, you can use it's socks5(h) service as a parent proxy to enjoy the sower's intelligent router.


## Installation
To enjoy the sower, you need to deploy sower on both server-side and client-side.

The installation script has been integrated into the sower. You can install sower as system service by running `./sower -install 'xxx'`

## Server
*If you already have another proxy solution with socks5h support, you can skip server-side.*

At the server-side, the sower runs just like a web server proxy.
It redirects HTTP requests to HTTPS and proxy https requests to the upstream HTTP service.
You can use your certificate or use the auto-generated certificate by the sower.

What you must set is the upstream HTTP service. You can set it by parameter `-s`, eg:
``` shell
# sower -s 127.0.0.1:8080
```

## Client
The easiest way to run it is:
``` shell
# sower -c aa.bb.cc # the `aa.bb.cc` can also be `socks5h://127.0.0.1:1080`
```
But a configuration file is recommended to persist dynamic rules in client side.

There are 3 kinds of proxy solutions, they are HTTP(S)_PROXY / DNS-based proxy / port-forward.

### HTTP(S)_PROXY
An HTTP(S)_PROXY listening on `:8080` is set by default if you run sower as client mode.

### DNS-based proxy
You can set the `serve_ip` field in the `dns` section in the configuration file to start the DNS-based proxy. You should also set the value of `serve_ip` as your default DNS in OS.

If you want to enjoy the full experience provided by the sower, you can take sower as your private DNS on a long-running server and set it as your default DNS in your router.

### port-forward
The port-forward can be only setted in configuration file, you can set it in section `client.router.port_mapping`, eg:
``` toml
[client.router.port_mapping]
":2222"="aa.bb.cc:22"
```


## Architecture
```
  relay   <--+       +-> target
http service |       |   service
     +-------+-------+----+
     |    sower server    |
     +----^-------^-------+
          80     443
301 http -+       +----- https
to https          |     service
               protected
                by tls
          socks5  |
 dns <---+   ^    |  +--> direct
relay    |   |    |  |   request
     +---+---+----+--+----+
     |    sower client    |
     +----^--^----^---^---+
          |  |    |   |
    dns --+  +    +   +-- port
            80  http(s)  forward
           443  proxy
```
