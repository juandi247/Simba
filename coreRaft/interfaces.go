package coreraft

type Runner interface {
	Start()
	Stop()
}

type TransportAdapter interface {
	SendMessage(Message)
}

type TimeAdapter interface {
	Now() int64
	Advance(int64)
	Sleep(int64)
}
