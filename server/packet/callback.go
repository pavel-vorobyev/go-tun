package packet

import "log"

type Callback interface {
	Call(call *CallbackCall)
}

type CallbackCall struct {
	Ptc  string
	Src  string
	Dst  string
	Data []byte
}

type TrafficCallback struct {
	t int
}

func (c *TrafficCallback) Call(args *CallbackCall) {
	l := len(args.Data)
	c.t = c.t + l
	log.Println(c.t)
}
