package coreraft

type Runner interface {
	Start()
	Stop()
}

type TransportAdapter interface {
	AppendEntries(AppendEntriesMessage) AppendEntriesResponse
	RequestVote(RequestVoteMessage) RequestVoteResponse
}

type TimeAdapter interface {
	Now() int64
	Advance(int64)
	Sleep(int64)
}

type StorageAdapter interface {
	appendEntryLog()
}
