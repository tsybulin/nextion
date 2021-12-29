package nextion

import (
	"errors"
	"fmt"
)

type page struct {
	nextion      Nextion
	id           uint8
	name         string
	popHandlers  map[string]TouchHandler
	pushHandlers map[string]TouchHandler
	objects      map[uint8]Object
}

func (p *page) ObjectId() uint8 {
	return p.id
}

func (p *page) SetBackgroundColor(bco Color565) {
	p.nextion.Send(fmt.Sprintf("%s.bco=%d", p.name, bco), RET_ACTION_OK)
}

func (p *page) BackgroundColor() (Color565, error) {
	ret := p.nextion.Send(fmt.Sprintf("get %s.bco", p.name), RET_NUMERIC_DATA)
	if ret.Result != NUMERIC_DATA {
		return BLACK, errors.New("command error")
	}

	return Color565(ret.Int), nil
}

func (p *page) Pop() {
	for _, h := range p.popHandlers {
		h(p)
	}
}

func (p *page) Push() {
	for _, h := range p.pushHandlers {
		h(p)
	}
}

func (p *page) AttachPop(target string, handler TouchHandler) {
	p.popHandlers[target] = handler
}

func (p *page) AttachPush(target string, handler TouchHandler) {
	p.pushHandlers[target] = handler
}

func (p *page) Touchables() map[uint8]Touchable {
	ts := make(map[uint8]Touchable)
	for _, o := range p.objects {
		if t, ok := o.(Touchable); ok {
			ts[o.ObjectId()] = t
		}
	}
	return ts
}

func (p *page) NewText(id uint8, name string) Text {
	t := &text{
		id:           id,
		objname:      name,
		nextion:      p.nextion,
		popHandlers:  make(map[string]TouchHandler),
		pushHandlers: make(map[string]TouchHandler),
	}

	p.objects[id] = t

	return t
}

func (p *page) NewNumber(id uint8, name string) Number {
	n := &number{
		id:           id,
		objname:      name,
		nextion:      p.nextion,
		popHandlers:  make(map[string]TouchHandler),
		pushHandlers: make(map[string]TouchHandler),
	}

	p.objects[id] = n

	return n
}

func (p *page) NewButton(id uint8, name string) Button {
	o := &button{
		id:           id,
		objname:      name,
		nextion:      p.nextion,
		popHandlers:  make(map[string]TouchHandler),
		pushHandlers: make(map[string]TouchHandler),
	}

	p.objects[id] = o

	return o
}

func (p *page) NewPicture(id uint8, name string) Picture {
	o := &picture{
		id:           id,
		objname:      name,
		nextion:      p.nextion,
		popHandlers:  make(map[string]TouchHandler),
		pushHandlers: make(map[string]TouchHandler),
	}

	p.objects[id] = o

	return o
}

type Page interface {
	Object
	Touchable
	Touchables() map[uint8]Touchable
	BackgroundColor() (Color565, error)
	SetBackgroundColor(bco Color565)
	NewText(id uint8, name string) Text
	NewNumber(id uint8, name string) Number
	NewButton(id uint8, name string) Button
}
