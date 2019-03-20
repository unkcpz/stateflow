package flowmat

type Port struct {
  channel chan interface{}
  cache interface{}
}
