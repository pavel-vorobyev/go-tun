package network

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	IffNoPi   = 0x10
	IffTun    = 0x01
	TunSetIff = 0x400454CA
)

type Interface struct {
	Mtu  int
	Name string
	fd   *os.File
}

func NewInterface(name string, mtu int) (*Interface, error) {
	fd, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	ifr := make([]byte, 18)
	copy(ifr, name)
	ifr[16] = IffTun
	ifr[17] = IffNoPi

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd.Fd(), uintptr(TunSetIff), uintptr(unsafe.Pointer(&ifr[0])))
	if errno != 0 {
		return nil, fmt.Errorf("ioctl open tun failed")
	}

	err = syscall.SetNonblock(int(fd.Fd()), false)
	if err != nil {
		return nil, err
	}

	return &Interface{
		Mtu:  mtu,
		Name: name,
		fd:   fd,
	}, nil
}

func (t *Interface) Read(data []byte) (int, error) {
	return t.fd.Read(data)
}

func (t *Interface) Write(data []byte) (int, error) {
	return t.fd.Write(data)
}

func (t *Interface) Close() error {
	return t.fd.Close()
}

func (t *Interface) GetMtu() int {
	return t.Mtu
}
