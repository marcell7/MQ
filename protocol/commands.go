package protocol

type Command int

const (
	CMD_PUBREG Command = iota
	CMD_SUBREG
	CMD_PUB
	CMD_SUB
	CMD_RECV
	CMD_RESP
	CMD_OK
	CMD_ERROR
)
