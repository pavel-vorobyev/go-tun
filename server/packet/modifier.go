package packet

type Modifier interface {
	Process(ptc string, src string, dst string, data []byte) ([]byte, error)
}
