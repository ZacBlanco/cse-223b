package main

type ActorState struct {
	numInvocations int64
}

func Main(params map[string]interface{}, state *interface{}) map[string]interface{} {
	if *state == nil {
		*state = &ActorState{
			numInvocations: 1,
		}
	}
	x := (*state).(*ActorState)
	x.numInvocations += 1
	return map[string]interface{}{
		"state": x.numInvocations,
	}
}
