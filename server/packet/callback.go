package packet

type Callback interface {
	Call(call *CallbackCall)
}

type CallbackCall struct {
	Ptc string
	Src string
	Dst string
	Len int
}

type TrafficCallback struct {
	T int
}

func (c *TrafficCallback) Call(args *CallbackCall) {
	c.T = c.T + args.Len
}
