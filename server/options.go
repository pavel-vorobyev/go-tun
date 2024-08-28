package server

import (
	"go-tun/server/config"
)

type Options struct {
	configProvider config.Provider
	rxModifiers    []PacketModifier
	txModifiers    []PacketModifier
	rxCallbacks    []PacketCallback
	txCallbacks    []PacketCallback
}

func NewOptions() *Options {
	return &Options{
		configProvider: &config.DefaultConfigProvider{},
		rxModifiers:    make([]PacketModifier, 0),
		txModifiers:    make([]PacketModifier, 0),
		rxCallbacks:    make([]PacketCallback, 0),
		txCallbacks:    make([]PacketCallback, 0),
	}
}

func (opt *Options) SetCustomConfigProvider(cp config.Provider) {
	opt.configProvider = cp
}

func (opt *Options) AddRxModifier(m PacketModifier) {
	opt.rxModifiers = append(opt.rxModifiers, m)
}

func (opt *Options) AddTxModifier(m PacketModifier) {
	opt.txModifiers = append(opt.txModifiers, m)
}

func (opt *Options) AddRxCallback(c PacketCallback) {
	opt.rxCallbacks = append(opt.rxCallbacks, c)
}

func (opt *Options) AddTxCallback(c PacketCallback) {
	opt.txCallbacks = append(opt.txCallbacks, c)
}
