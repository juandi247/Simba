package coreraft

type MessageType int

const (
	HEARTBEAT MessageType = iota
	ELECTION              //todo: check this if there is only votation, or start election and votation separately
	APPEND
	ACK
)


type Message struct {
	Sender      int
	Receiver    int
	Term        uint64
	MessageType MessageType
	LogIndex    uint64
	// ONLY used for simulator
	DeliveryTick int64
}