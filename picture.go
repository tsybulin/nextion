package nextion

import "fmt"

type picture struct {
	nextion      Nextion
	id           uint8
	objname      string
	popHandlers  map[string]TouchHandler
	pushHandlers map[string]TouchHandler
}

func (o *picture) ObjectId() uint8 {
	return o.id
}

func (o *picture) SetVisible(visible bool) {
	v := 0
	if visible {
		v = 1
	}
	o.nextion.Send(fmt.Sprintf("vis %d,%d", o.id, v), RET_ACTION_OK)
}

func (o *picture) SetPicture(pic uint8) {
	o.nextion.Send(fmt.Sprintf("%s.pic=%d", o.objname, pic), RET_ACTION_OK)
}

func (o *picture) Pop() {
	for _, h := range o.popHandlers {
		h(o)
	}
}

func (o *picture) Push() {
	for _, h := range o.pushHandlers {
		h(o)
	}
}

func (o *picture) AttachPop(target string, handler TouchHandler) {
	o.popHandlers[target] = handler
}

func (o *picture) AttachPush(target string, handler TouchHandler) {
	o.pushHandlers[target] = handler
}

type Picture interface {
	Object
	Touchable
	SetVisible(visible bool)
	SetPicture(pic uint8)
}
