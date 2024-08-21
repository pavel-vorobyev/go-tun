package config

type Config struct {
	RemoteIp   string `json:"RemoteIp"`
	RemotePort int    `json:"RemotePort"`
	TunName    string `json:"TunName"`
	TunIp      string `json:"TunIp"`
	TunCidr    int    `json:"TunCidr"`
	Mtu        int    `json:"Mtu"`
}
