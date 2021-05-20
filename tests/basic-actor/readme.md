## Simple Actor Test

This test demonstrates passing state quickly from one invocation to the
next. `actor.go` is the example file which allows the state passing:

`actor.go`

```go
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
```

The result of this function will increase a counter for every activation
that comes to this particular action. When OpenWhisk finishes the
execution of the action, the state will stick around until the action is
destroyed (via timeout, or _eventually_ checkpointed)

To deploy and test:

- download and install the openwhisk CLI `wsk` and put it on your `$PATH`
- start the `StandaloneOpenWhisk` process (via vscode or `./gradlew :core:standalone:bootRun`)

Run the following commands:

```console
$ wsk -i --apihost http:;172.17.0.1:3233 action create actor-test actor.go --kind go:1.15-actor
$ wsk -i --apihost http://172.17.0.1:3233 action invoke actor-test -r
{
    "state": 2
}
$ wsk -i --apihost http://172.17.0.1:3233 action invoke actor-test -r
{
    "state": 3
}
```

You now have a serverless function with (kind of) persistent state!


