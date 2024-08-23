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
		rxModifiers:     make([]packet.Modifier, 0),
		txModifiers:     make([]packet.Modifier, 0),
		rxCallbacks:     make([]packet.Callback, 0),
		txCallbacks:     make([]packet.Callback, 0),
	}
}

func (opt *Options) SetCustomConfigProvider(cp config.Provider) {
	opt.configProvider = cp
}

func (opt *Options) SetCustomSrcAddressKeyFactory(kf address.CAddrKeyFactory) {
	opt.cAddrKeyFactory = kf
}

func (opt *Options) SetCustomSrcAddressStore(s address.CAddrStore) {
	opt.cAddrStore = s
}

func (opt *Options) AddRxModifier(m packet.Modifier) {
	opt.rxModifiers = append(opt.rxModifiers, m)
}

func (opt *Options) AddTxModifier(m packet.Modifier) {
	opt.txModifiers = append(opt.txModifiers, m)
}

func (opt *Options) AddRxCallback(c packet.Callback) {
	opt.rxCallbacks = append(opt.rxCallbacks, c)
}

func (opt *Options) AddTxCallback(c packet.Callback) {
	opt.txCallbacks = append(opt.txCallbacks, c)
}
