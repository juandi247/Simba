package simulator

import "math/rand"


type FuzzyLevel int

const (
	LOW FuzzyLevel = iota 
	MEDIUM
	HIGH
)

type FuzzyConfig struct {
	rand *rand.Rand
	LatencyProb     float64
	MessageLostProb float64
	NodeCrashProb 	float64
}

var FuzzyConfigMap = map[FuzzyLevel]FuzzyConfig{
	LOW: {
		LatencyProb:     0.01,
		MessageLostProb: 0.01,
		NodeCrashProb: 0.01,
	},
	MEDIUM: {
		LatencyProb:     0.05,
		MessageLostProb: 0.05,
		NodeCrashProb: 2.00,
	},
	HIGH: {
		LatencyProb:     0.10,
		MessageLostProb: 0.10,
		NodeCrashProb: 5.00,
	},
}

func FuzzyConfiguration(seed int64, fuzzyLevel FuzzyLevel) (FuzzyConfig) {

	source := rand.NewSource(seed)
	rand := rand.New(source)

	fuzzyConfig, exists := FuzzyConfigMap[fuzzyLevel]
	if !exists {
		panic("Assertion for variable of fuzzy level")
	}
	fuzzyConfig.rand=rand

	return fuzzyConfig
}

/*
this functions takes the configuration, and returns
int: the number of ticks
bool: if the message its dropped or not

*/


func (fc *FuzzyConfig) RandomizeNetwork(){
	/*
		There its the log from the fuzzy using the RAND, to obtain the probabilities of a message being dropped.
		TODO:
		  ?- (i think this is the bestone) Pass the message as parameter and return it so that its added to the QUEUE
		  ?- ?Or dont receive any data and return the parametres of latency, message drop as return values, and then after assing them
	*/
}
