package flowmat

type Port struct {
  channel chan interface{}
  cache interface{}
}

func (p *Port) Feed(data interface{}) {
  if data == nil {
    data = p.cache
  }
  p.channel <- data
}

func (p *Port) Extract() interface{} {
  // return <-port.channel
  <-p.channel
  return p.cache
}
