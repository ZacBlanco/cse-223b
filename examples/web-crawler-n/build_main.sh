#!/bin/bash

NUM_ACTORS=$1
TEMPLATE_FILE="template.go"
NEW_FILE="main_$NUM_ACTORS.go"

function printMainChildActor()  {
    echo "func Mainchild$1(args map[string]interface{}, state *interface{}) map[string]interface{} {"
    echo "	const actorId = $1"
    echo "	seeds := GetSeedsFromArgs(args)"
    echo "	rand.Seed(time.Now().UnixNano())"
    echo "	rand.Shuffle(len(seeds), func(i, j int) { seeds[i], seeds[j] = seeds[j], seeds[i] })"
    echo "	actor := NewActorWithId(actorId)"
    echo "	actor.UseLatestState(state)"
    echo "	ret := actor.StartWebCrawlerAndReturnWebCrawlerState(seeds)"
    echo "	*state = actor.State"
    echo "	fmt.Println(\"Save state:\", *state)"
    echo "	return ret"
    echo "}"
    echo
}

function build() {
    cp $TEMPLATE_FILE $NEW_FILE

    # print constants
    echo "const NUM_ACTORS = $NUM_ACTORS" >> $NEW_FILE
    echo >> $NEW_FILE

    # prints NUM_ACTOR number of actor functions
    for i in $(seq 1 $NUM_ACTORS); do
        ID=$(($i-1))
        printMainChildActor $ID >> $NEW_FILE
    done

    echo "Created $NEW_FILE"

}

build
