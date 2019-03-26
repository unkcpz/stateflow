package gflow

type Port struct {
	channel chan interface{}
	cache   interface{}
}

// Feed send data to port channel, if nil send cache
func (p *Port) Feed(data interface{}) {
	if data == nil {
		data = p.cache
	}
	p.channel <- data
}

// Extract write port channel value to its cache and return cache
func (p *Port) Extract() interface{} {
	// return <-port.channel
	p.cache = <-p.channel
	return p.cache
}
