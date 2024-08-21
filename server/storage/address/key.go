package address

import "fmt"

type CAddrKeyFactory interface {
	Get(ptc string, src string, dst string) string
}

type DefaultCAddrKeyFactory struct{}

func (kp *DefaultCAddrKeyFactory) Get(ptc string, src string, dst string) string {
	return fmt.Sprintf("%s/%s:%s", ptc, src, dst)
}
