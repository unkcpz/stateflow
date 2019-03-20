package flowmat

type port struct {
  channel chan interface{}
  cache interface{}
}
