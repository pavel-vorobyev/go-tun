package packet

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

func (c *TrafficCallback) Call(args CallbackCall) {
	c.t = c.t + len(args.Data)
}
