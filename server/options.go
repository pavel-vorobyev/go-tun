package server

import (
	"go-tun/server/config"
	"go-tun/server/packet"
	"go-tun/server/storage/address"
)

type Options struct {
	configProvider  config.Provider
	cAddrKeyFactory address.CAddrKeyFactory
	cAddrStore      address.CAddrStore
	rxModifiers     []packet.Modifier
	txModifiers     []packet.Modifier
	rxCallbacks     []packet.Callback
	txCallbacks     []packet.Callback
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

func (opt *Options) AddRxModifier(m packet.Modifier) *Options {
	opt.rxModifiers = append(opt.rxModifiers, m)
	return opt
}

func (opt *Options) AddTxModifier(m packet.Modifier) *Options {
	opt.txModifiers = append(opt.txModifiers, m)
	return opt
}

func (opt *Options) AddRxCallback(c packet.Callback) *Options {
	opt.rxCallbacks = append(opt.rxCallbacks, c)
	return opt
}

func (opt *Options) AddTxCallback(c packet.Callback) *Options {
	opt.txCallbacks = append(opt.txCallbacks, c)
	return opt
}
