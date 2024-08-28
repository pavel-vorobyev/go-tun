package transport

type Data struct {
	data  []byte
	cAddr string
}

func NewData(data []byte, cAddr string) *Data {
	return &Data{
		data:  data,
		cAddr: cAddr,
	}
}

func (d *Data) GetData() []byte {
	return d.data
}

func (d *Data) GetCAddr() string {
	return d.cAddr
}
