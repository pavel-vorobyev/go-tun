package packet

type Callback interface {
	Call(Ptc int, Src string, Dst string, N int)
}

type TrafficCallback struct {
	T int
}

func (c *TrafficCallback) Call(_ int, _ string, _ string, N int) {
	c.T = c.T + N
}
