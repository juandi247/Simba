package adapters


import (
	raft "simba/raft"
)

type Runner interface {
	Start()
	Stop()
}

type TransportAdapter interface {
	SendMessage([]raft.Message)
}

type TimeAdapter interface {
	Now() int64
	Advance(int64)
	Sleep(int64)
}

type StorageAdapter interface {
	appendEntryLog()
}
