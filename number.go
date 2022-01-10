package nextion

import (
	"errors"
	"fmt"
)

type number struct {
	nextion      Nextion
	id           uint8
	objname      string
	popHandlers  map[string]TouchHandler
	pushHandlers map[string]TouchHandler
}

func (n *number) ObjectId() uint8 {
	return n.id
}

func (n *number) SetValue(v int32) {
	n.nextion.Send(fmt.Sprintf("%s.val=%d", n.objname, v), RET_ACTION_OK)
}

func (n *number) Value() (int32, error) {
	ret := n.nextion.Send(fmt.Sprintf("get %s.val", n.objname), RET_NUMERIC_DATA)
	if ret.Result != NUMERIC_DATA {
		return 0, errors.New("command error")
	}

	return ret.Int, nil
}

func (n *number) SetVisible(visible bool) {
	v := 0
	if visible {
		v = 1
	}
	n.nextion.Send(fmt.Sprintf("vis %d,%d", n.id, v), RET_ACTION_OK)
}

func (n *number) SetBackgroundColor(bco Color565) {
	n.nextion.Send(fmt.Sprintf("%s.bco=%d", n.objname, bco), RET_ACTION_OK)
}

func (n *number) SetTextColor(pco Color565) {
	n.nextion.Send(fmt.Sprintf("%s.pco=%d", n.objname, pco), RET_ACTION_OK)
}

func (n *number) Pop() {
	for _, h := range n.popHandlers {
		h(n)
	}
}

func (n *number) Push() {
	for _, h := range n.pushHandlers {
		h(n)
	}
}

func (n *number) AttachPop(target string, handler TouchHandler) {
	n.popHandlers[target] = handler
}

func (n *number) AttachPush(target string, handler TouchHandler) {
	n.pushHandlers[target] = handler
}

type Number interface {
	Object
	Touchable
	SetValue(v int32)
	Value() (int32, error)
	SetVisible(visible bool)
	SetBackgroundColor(bco Color565)
	SetTextColor(pco Color565)
}
