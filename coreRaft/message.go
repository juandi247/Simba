package coreraft

type MessageType int

const (
	HEARTBEAT MessageType = iota
	APPEND
	ACK
	REQUESTVOTE
	VOTECONFIRMATION
	UPDATEYOURDATA
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