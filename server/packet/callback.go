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
	T int
}

func (c *TrafficCallback) Call(args *CallbackCall) {
	l := len(args.Data)
	c.T = c.T + l
}
