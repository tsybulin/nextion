package main

import (
	"log"
	"sync"
	"time"

	"github.com/tsybulin/nextion"
)

const PORT = "/dev/tty.usbmodem401"

// const PORT = "/dev/serial0"

var (
	nxt        nextion.Nextion
	numCopies  nextion.Number
	txtQuality nextion.Text
	numDPI     nextion.Number
)

var gauge nextion.Number
var q uint32
var busy = false

func progress(b nextion.Button) {
	busy = true
	ticker := time.NewTicker(100 * time.Millisecond)
	q = 0
	for range ticker.C {
		q = q + 1
		if q < 360 {
			gauge.SetValue(q)
		} else {
			break
		}
	}

	ticker.Stop()
	q = 0
	nxt.Send("g0.en=0", nextion.RET_ACTION_OK)
	nxt.Send("vis g0,0", nextion.RET_ACTION_OK)
	busy = false
}

func btnPrintDidPush(o nextion.Object) {
	if busy {
		return
	}

	go progress(o.(nextion.Button))
}

func btnUpDidPush(o nextion.Object) {
	n, err := numCopies.Value()
	if err == nil {
		n = n + 1
		if n > 99 {
			n = 99
		}
		numCopies.SetValue(n)
	}
}

func btnDownDidPush(o nextion.Object) {
	n, err := numCopies.Value()
	if err == nil {
		n = n - 1
		if n < 1 {
			n = 1
		}
		numCopies.SetValue(n)
	}
}

func cbGreyDidPush(o nextion.Object) {
	n, err := o.(nextion.Number).Value()
	if err == nil {
		if n == 0 {
			txtQuality.SetText("DRAFT")
		} else {
			txtQuality.SetText("FINE")
		}
	}
}

func cbVerticalDidPush(o nextion.Object) {
	n, err := o.(nextion.Number).Value()
	if err == nil {
		if n == 0 {
			numDPI.SetValue(300)
		} else {
			numDPI.SetValue(600)
		}
	}
}

func main() {
	nxt = nextion.NewNextion(PORT)

	err := nxt.Init(9600)
	if err != nil {
		log.Fatal(err)
	}

	defer nxt.Close()

	nxt.SetDim(33)

	err = nxt.SetBaud(57600)
	if err != nil {
		log.Print("SetBaud error: ", err)
	}

	page := nxt.NewPage(0, "mainPage")

	gauge = page.NewNumber(11, "z0")
	btnPrint := page.NewButton(2, "btnPrint")
	btnPrint.AttachPush("main", btnPrintDidPush)

	numCopies = page.NewNumber(10, "numCopies")
	page.NewButton(12, "btnUp").AttachPush("main", btnUpDidPush)
	page.NewButton(13, "btnDown").AttachPush("main", btnDownDidPush)

	txtQuality = page.NewText(3, "txtQuality")
	page.NewNumber(5, "cbGrey").AttachPush("main", cbGreyDidPush)
	numDPI = page.NewNumber(4, "numDPI")
	page.NewNumber(6, "cbVertical").AttachPush("main", cbVerticalDidPush)

	nxt.ShowPage(0)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
