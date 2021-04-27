package core

type XContext struct {
	RequestJson string
	Response    string
	JsonParser  Parser
}

type HandlerChain struct {
	ctx XContext
}

func (hc *HandlerChain) HandleALl() {

}

type Parser interface {
}

type HandlerAdapter interface {
	Name() string
	Handle(ctx *XContext)
	DoNext()
}
