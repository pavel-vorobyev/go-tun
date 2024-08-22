package server

import (
	"go-tun/server/config"
	"go-tun/server/storage/address"
)

type Options struct {
	configProvider  config.Provider
	cAddrKeyFactory address.CAddrKeyFactory
	cAddrStore      address.CAddrStore
}

func CreateOptions() *Options {
	return &Options{
		configProvider:  &config.DefaultConfigProvider{},
		cAddrKeyFactory: &address.DefaultCAddrKeyFactory{},
		cAddrStore:      address.NewDefaultCAddrStore(),
	}
}

func (opt *Options) SetCustomConfigProvider(cp config.Provider) *Options {
	opt.configProvider = cp
	return opt
}

func (opt *Options) SetCustomSrcAddressKeyFactory(kf address.CAddrKeyFactory) *Options {
	opt.cAddrKeyFactory = kf
	return opt
}

func (opt *Options) SetCustomSrcAddressStore(s address.CAddrStore) *Options {
	opt.cAddrStore = s
	return opt
}
