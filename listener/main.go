package listener

type Endpoint struct {
	Address string
	Port    int
}

var listen Endpoint

type Listener interface {
	Create(Endpoint, *chan string)
}
