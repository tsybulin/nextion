package nextion

type TouchHandler func(Object)

type Touchable interface {
	AttachPop(string, TouchHandler)
	AttachPush(string, TouchHandler)
	Pop()
	Push()
}
