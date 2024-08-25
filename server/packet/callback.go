package packet

type Callback interface {
	Call(call *CallbackCall)
}

type CallbackCall struct {
	Ptc string
	Src string
	Dst string
	N   int
}

type TrafficCallback struct {
	T int
}

func (c *TrafficCallback) Call(args *CallbackCall) {
	c.T = c.T + args.N
}
