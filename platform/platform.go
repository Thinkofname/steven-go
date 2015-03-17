package platform

func Init(handler Handler) {
	run(handler)
}

type Handler struct {
	Start func()
	Draw  func()
}
