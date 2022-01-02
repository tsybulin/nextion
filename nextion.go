package nextion

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

type CommandResult struct {
	Result byte
	Int    uint32
	Str    string
}

type nextion struct {
	port     serial.Port
	portname string
	output   chan string
	result   chan CommandResult
	pageIdx  uint8
	pages    map[uint8]Page
}

var endNextionMessage = []byte{0xFF, 0xFF, 0xFF}

func stringToHexBytes(s string) []byte {
	hx := hex.EncodeToString([]byte(s))
	decoded, err := hex.DecodeString(hx)
	if err != nil {
		log.Println("Nextion.stringToHexBytes error", err)
	}

	decoded = append(decoded, endNextionMessage...)

	return decoded
}

func (n *nextion) Send(s string, action RetAction) CommandResult {
	n.output <- s
	rslt := <-n.result
	if rslt.Result == INSTRUCTION_SUCCESS {
		//log.Printf("Nextion.Send %s result OK", s)
		return rslt
	}

	if rslt.Result == NUMERIC_DATA && action == RET_NUMERIC_DATA {
		// log.Printf("Nextion.Send %s numeric result OK: %d", s, rslt.Int)
		return rslt
	}

	if rslt.Result == STRING_DATA && action == RET_STRING_DATA {
		// log.Printf("Nextion.Send %s string result OK: %s", s, rslt.Str)
		return rslt
	}

	log.Printf("Nextion.Send %s result error: %02x", s, rslt.Result)
	return rslt
}

func (n *nextion) outHandler() {
	for s := range n.output {
		bts := stringToHexBytes(s)
		// hx := strings.ToUpper(hex.EncodeToString(bts))
		// log.Println("Nextion.outHandler \""+s+"\" :", hx)

		n, err := n.port.Write(bts)

		if err != nil {
			log.Print("Nextion.Send error:", err)
		}
		if len(bts) != n {
			log.Print("Nextion.Send len:", len(bts), "sent:", n)
		}
	}
}

func (n *nextion) inpHandler() {
	buffer := make([]byte, 255)
	total := make([]byte, 0)

	for {
		c, err := n.port.Read(buffer)
		if err != nil {
			log.Println("port read error:", err)
		}

		if c > 0 {
			total = append(total, buffer[0:c]...)
			//log.Printf("Nextion.inpHandler got %d bytes: %s", c, hex.EncodeToString(buffer[0:c]))

			for {
				i := bytes.Index(total, endNextionMessage)
				if i < 0 {
					break
				}
				s := hex.EncodeToString(total[0:i])

				if total[0] == INSTRUCTION_SUCCESS {
					// log.Println("Nextion.inpHandler ", s, "Instruction success")
					n.result <- CommandResult{Result: INSTRUCTION_SUCCESS}
				} else if total[0] == INSTRUCTION_INVALID {
					log.Println("Nextion.inpHandler ", s, "Instruction invalid")
					n.result <- CommandResult{Result: INSTRUCTION_INVALID}
				} else if total[0] == ASSIGNMENT_ERROR {
					log.Println("Nextion.inpHandler ", s, "Assignment error")
					n.result <- CommandResult{Result: ASSIGNMENT_ERROR}
				} else if total[0] == NUMERIC_DATA {
					// log.Printf("Nextion.inpHandler NUMERIC_DATA %s", s)
					n.result <- CommandResult{
						Result: NUMERIC_DATA,
						Int:    uint32(total[4])<<24 + uint32(total[3])<<16 + uint32(total[2])<<8 + uint32(total[1]),
					}
				} else if total[0] == STRING_DATA {
					// log.Printf("Nextion.inpHandler STRING_DATA %s", s)
					n.result <- CommandResult{
						Result: STRING_DATA,
						Str:    string(total[1:i]),
					}
				} else if total[0] == TOUCH_EVENT {
					pid := uint8(total[1])
					cid := uint8(total[2])

					// state := ""
					// if total[3] == 1 {
					// 	state = "press"
					// } else if total[3] == 0 {
					// 	state = "release"
					// }
					// log.Println("Nextion.inpHandler ", s, "Touch event pid", pid, "cid", cid, "state", state)

					if p := n.pages[pid]; p != nil {
						if cid == 0 {
							if t, ok := p.(Touchable); ok {
								if total[3] == 1 {
									go t.Push()
								} else if total[3] == 0 {
									go t.Pop()
								}
							}
						} else {
							if total[3] == 1 {
								if t := p.Touchables()[cid]; t != nil {
									go t.Push()
								}
							} else if total[3] == 0 {
								if t := p.Touchables()[cid]; t != nil {
									go t.Pop()
								}
							}
						}
					}
				} else {
					log.Print(s, "Unknown")
				}

				total = total[i+3:]
			}
		}
	}
}

func (n *nextion) Init(baud int) error {
	mode := serial.Mode{
		BaudRate: baud,
		DataBits: 8,
		StopBits: serial.OneStopBit,
		Parity:   serial.NoParity,
	}

	port, err := serial.Open(n.portname, &mode)
	if err != nil {
		return err
	}

	n.port = port

	port.ResetInputBuffer()
	port.ResetOutputBuffer()

	err = port.SetReadTimeout(time.Millisecond * 200)
	if err != nil {
		return err
	}

	go n.outHandler()
	go n.inpHandler()

	n.Send("bkcmd=3", RET_ACTION_OK)

	return nil
}

func (n *nextion) Close() error {
	close(n.output)
	close(n.result)
	return n.port.Close()
}

func (n *nextion) SetBaud(baud int) error {
	if n.port == nil {
		return errors.New("not initialized")
	}

	n.Send(fmt.Sprintf("baud=%d", baud), RET_ACTION_OK)

	mode := serial.Mode{
		BaudRate: baud,
		DataBits: 8,
		StopBits: serial.OneStopBit,
		Parity:   serial.NoParity,
	}

	return n.port.SetMode(&mode)
}

func (n *nextion) SetDim(dim int) {
	n.Send(fmt.Sprintf("dim=%d", dim), RET_ACTION_OK)
}

func (n *nextion) SetCurrentDateTime() error {
	now := time.Now()
	if rslt := n.Send(fmt.Sprintf("rtc0=%d", now.Year()), RET_ACTION_OK); rslt.Result != INSTRUCTION_SUCCESS {
		return errors.New("year set error")
	}

	if rslt := n.Send(fmt.Sprintf("rtc1=%d", now.Month()), RET_ACTION_OK); rslt.Result != INSTRUCTION_SUCCESS {
		return errors.New("month set error")
	}

	if rslt := n.Send(fmt.Sprintf("rtc2=%d", now.Day()), RET_ACTION_OK); rslt.Result != INSTRUCTION_SUCCESS {
		return errors.New("day set error")
	}

	if rslt := n.Send(fmt.Sprintf("rtc3=%d", now.Hour()), RET_ACTION_OK); rslt.Result != INSTRUCTION_SUCCESS {
		return errors.New("hour set error")
	}

	if rslt := n.Send(fmt.Sprintf("rtc4=%d", now.Minute()), RET_ACTION_OK); rslt.Result != INSTRUCTION_SUCCESS {
		return errors.New("minute set error")
	}
	return nil
}

func (n *nextion) NewPage(id uint8, name string) Page {
	p := &page{
		nextion:      n,
		id:           id,
		name:         name,
		popHandlers:  make(map[string]TouchHandler),
		pushHandlers: make(map[string]TouchHandler),
		objects:      map[uint8]Object{},
	}

	n.pages[id] = p

	return p
}

func (n *nextion) ShowPage(id uint8) {
	if p := n.pages[id]; p != nil {
		n.pageIdx = p.ObjectId()
		n.Send(fmt.Sprintf("page %d", n.pageIdx), RET_ACTION_OK)
	}
}

type Nextion interface {
	Init(baud int) error
	Close() error
	SetBaud(baud int) error
	SetDim(dim int)
	SetCurrentDateTime() error
	Send(s string, action RetAction) CommandResult
	NewPage(id uint8, name string) Page
	ShowPage(id uint8)
}

func NewNextion(port string) Nextion {
	nextion := &nextion{
		portname: port,
		output:   make(chan string),
		result:   make(chan CommandResult, 10),
		pageIdx:  0,
		pages:    map[uint8]Page{},
	}

	return nextion
}
