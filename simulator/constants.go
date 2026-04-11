package simulator


const TickFrequency = 1
const maxTicks = 1000000
const maxQueueSize = 100000
const maxInboxSize = 10000

const LeaderHeartbeatFreq = 10
const ElectionTimeout = 10
const MinFollowerTimeout = 30
const MaxFollowerTimeout = 50

/*
Assertions for ranges inside the constants, to have it in compile time. 
*/
func init() {
	if LeaderHeartbeatFreq >= MinFollowerTimeout || LeaderHeartbeatFreq >= MaxFollowerTimeout {
		panic("the leader heart beat MUST be smaller than both minfollower and maxfollower timeouts")
	}

	if MinFollowerTimeout >= MaxFollowerTimeout {
		panic("the minfollowre must be smaller than the max")
	}

	if minCrashNodeDowntime >= maxCrashNodeDowntime {
		panic("the minCrashNodeDowntime must be smaller than the max")
	}

	if minLatencyDelay >= maxCrashNodeDowntime {
		panic("the minCrashNodeDowntime must be smaller than the max")
	}
}

