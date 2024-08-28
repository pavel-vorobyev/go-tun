package server

type PacketModifier interface {
	Process(data []byte) ([]byte, error)
}
