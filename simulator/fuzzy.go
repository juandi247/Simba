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

const minCrashNodeDowntime= 10
const maxCrashNodeDowntime=50

const minLatencyDelay= 10
const maxLatencyDelay=50


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


func (fc *FuzzyConfig) determineCrashingOfNode()(bool, int64){

	randomNumber:= fc.rand.Float64()

	if randomNumber > fc.NodeCrashProb {
		comeBackToLiveTick:=  minCrashNodeDowntime + fc.rand.Int63n(maxCrashNodeDowntime - minCrashNodeDowntime + 1)
		return true, comeBackToLiveTick
	}
	return false, 0
}


/*
this functions takes the configuration, and returns
bool: if the message its dropped or not
int: the number of delay ticks

*/
func (fc *FuzzyConfig) RandomizeNetwork() (bool, int64){
	randomNumber:= fc.rand.Float64()

	if randomNumber > fc.MessageLostProb {
		// The message was determinted to be LOST
		return true, 0
	}

	delayTicks:=  minLatencyDelay + fc.rand.Int63n(maxLatencyDelay - minLatencyDelay + 1)
	return false, delayTicks

}
