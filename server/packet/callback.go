package packet

type Callback interface {
	Call(Ptc string, Src string, Dst string, N int)
}

type TrafficCallback struct {
	T int
}

func (c *TrafficCallback) Call(_ string, _ string, _ string, N int) {
	c.T = c.T + N
}
