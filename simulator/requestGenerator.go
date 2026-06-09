package simulator

import (
	"math/rand"
	"simba/raft"
)


const numberOfRequests = 1000
const maxEntryLength = 10
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")


func GenerateRequests(rng *rand.Rand) []SimMessage {

	arr:= make([]SimMessage, numberOfRequests)

	for i:=0; i<1000; i++{
	arr[i] = SimMessage{
			id: i,
			Message: raft.NewEntry{
				Command: GenerateRandomString(rng),
			},
			//minimum 20 and max the ticks
			DeliveryTick: 10 + rng.Intn(maxTicks-10 +1) ,
		}

	}
	return arr

}


func GenerateRandomString(rng *rand.Rand)string {
	arr:= make([]rune, len(letters))

	for i:=range arr{
 		number:= rng.Intn(len(letters))
		arr[i] = letters[number]
	}

return string(arr)	
}
