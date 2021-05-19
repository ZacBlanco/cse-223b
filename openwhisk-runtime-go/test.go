package main

import "fmt"

func main() {
	var state interface{}
	var obj map[string]interface{}

	for i := 0; i < 10; i++ {
		fmt.Println(action(obj, &state))
	}
	fmt.Println(state)
}

type MyState struct {
	state int
}

func action(obj map[string]interface{}, state *interface{}) map[string]interface{} {
	if *state == nil {
		tmp := MyState{}
		*state = &tmp
	}
	myState := (*state).(*MyState)
	myState.state += 1
	return map[string]interface{}{
		"state": myState.state,
	}
}
