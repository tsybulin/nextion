package nextion

import "fmt"

type button struct {
	nextion      Nextion
	id           uint8
	objname      string
	popHandlers  map[string]TouchHandler
	pushHandlers map[string]TouchHandler
}

func (o *button) ObjectId() uint8 {
	return o.id
}

func (o *button) SetText(s string) {
	o.nextion.Send(fmt.Sprintf("%s.txt=\"%s\"", o.objname, s), RET_ACTION_OK)
}

func (o *button) SetVisible(visible bool) {
	v := 0
	if visible {
		v = 1
	}
	o.nextion.Send(fmt.Sprintf("vis %d,%d", o.id, v), RET_ACTION_OK)
}

func (o *button) SetBackgroundColor(bco Color565) {
	o.nextion.Send(fmt.Sprintf("%s.bco=%d", o.objname, bco), RET_ACTION_OK)
}

func (o *button) SetBackgroundColorPushed(bco Color565) {
	o.nextion.Send(fmt.Sprintf("%s.bco2=%d", o.objname, bco), RET_ACTION_OK)
}

func (o *button) SetTextColor(pco Color565) {
	o.nextion.Send(fmt.Sprintf("%s.pco=%d", o.objname, pco), RET_ACTION_OK)
}

func (o *button) SetTextColorPushed(pco Color565) {
	o.nextion.Send(fmt.Sprintf("%s.pco2=%d", o.objname, pco), RET_ACTION_OK)
}

func (o *button) SetPicture(pic uint8) {
	o.nextion.Send(fmt.Sprintf("%s.pic=%d", o.objname, pic), RET_ACTION_OK)
}

func (o *button) SetPicturePushed(pic uint8) {
	o.nextion.Send(fmt.Sprintf("%s.pic2=%d", o.objname, pic), RET_ACTION_OK)
}

func (o *button) Pop() {
	for _, h := range o.popHandlers {
		h(o)
	}
}

func (o *button) Push() {
	for _, h := range o.pushHandlers {
		h(o)
	}
}

func (o *button) AttachPop(target string, handler TouchHandler) {
	o.popHandlers[target] = handler
}

func (o *button) AttachPush(target string, handler TouchHandler) {
	o.pushHandlers[target] = handler
}

type Button interface {
	Object
	Touchable
	SetText(s string)
	SetVisible(visible bool)
	SetBackgroundColor(bco Color565)
	SetTextColor(pco Color565)
	SetBackgroundColorPushed(bco Color565)
	SetTextColorPushed(pco Color565)
	SetPicture(pic uint8)
	SetPicturePushed(pic uint8)
}
