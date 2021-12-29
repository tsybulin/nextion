package nextion

const (
	INSTRUCTION_SUCCESS byte = 0x01
	INSTRUCTION_INVALID byte = 0x00
	ASSIGNMENT_ERROR    byte = 0x1c
	TOUCH_EVENT         byte = 0x65
	STRING_DATA         byte = 0x70
	NUMERIC_DATA        byte = 0x71
)

type RetAction byte

const (
	RET_ACTION_OK    RetAction = iota
	RET_NUMERIC_DATA RetAction = iota + 1
	RET_STRING_DATA  RetAction = iota + 2
)
