package ogmigo

type circular struct {
	index int
	data  [][]byte
}

func newCircular(cap int) *circular {
	return &circular{
		data: make([][]byte, cap),
	}
}

func (c *circular) add(data []byte) {
	c.data[c.index] = data
	c.index = (c.index + 1) % len(c.data)
}

func (c *circular) list() (data [][]byte) {
	for i := 0; i < len(c.data); i++ {
		offset := (c.index + i) % len(c.data)
		if v := c.data[offset]; len(v) > 0 {
			data = append(data, v)
		}
	}
	return data
}

func (c *circular) prefix(data ...[]byte) [][]byte {
	return append(data, c.list()...)
}
