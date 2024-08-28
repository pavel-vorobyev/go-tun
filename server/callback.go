package server

type PacketCallback interface {
	Call(Ptc int, Src string, Dst string, N int)
}

type TrafficPacketCallback struct {
	T int
}

func (c *TrafficPacketCallback) Call(_ int, _ string, _ string, N int) {
	c.T = c.T + N
}
