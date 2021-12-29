package nextion

import (
	"errors"
	"fmt"
)

type text struct {
	nextion      Nextion
	id           uint8
	objname      string
	popHandlers  map[string]TouchHandler
	pushHandlers map[string]TouchHandler
}

func (t *text) ObjectId() uint8 {
	return t.id
}

func (t *text) SetText(s string) {
	t.nextion.Send(fmt.Sprintf("%s.txt=\"%s\"", t.objname, s), RET_ACTION_OK)
}

func (t *text) Text() (string, error) {
	ret := t.nextion.Send(fmt.Sprintf("get %s.txt", t.objname), RET_STRING_DATA)
	if ret.Result != STRING_DATA {
		return "", errors.New("command error")
	}

	return ret.Str, nil
}

func (t *text) SetVisible(visible bool) {
	v := 0
	if visible {
		v = 1
	}
	t.nextion.Send(fmt.Sprintf("vis %d,%d", t.id, v), RET_ACTION_OK)
}

func (t *text) SetBackgroundColor(bco Color565) {
	t.nextion.Send(fmt.Sprintf("%s.bco=%d", t.objname, bco), RET_ACTION_OK)
}

func (t *text) SetTextColor(pco Color565) {
	t.nextion.Send(fmt.Sprintf("%s.pco=%d", t.objname, pco), RET_ACTION_OK)
}

func (t *text) Pop() {
	for _, h := range t.popHandlers {
		h(t)
	}
}

func (t *text) Push() {
	for _, h := range t.pushHandlers {
		h(t)
	}
}

func (t *text) AttachPop(target string, handler TouchHandler) {
	t.popHandlers[target] = handler
}

func (t *text) AttachPush(target string, handler TouchHandler) {
	t.pushHandlers[target] = handler
}

type Text interface {
	Object
	Touchable
	Text() (string, error)
	SetText(s string)
	SetVisible(visible bool)
	SetBackgroundColor(bco Color565)
	SetTextColor(pco Color565)
}
