package types

import (
	"fmt"
	"log"
)

var (
	AnyClass *Class = nil
)

// delayedReady holds types waiting to be intialised
var delayedReady []*Class

// TypeDelayReady stores the list of types to initialise
//
// Call MakeReady when all initialised
func TypeDelayReady(t *Class) {
	delayedReady = append(delayedReady, t)
}

// TypeMakeReady readies all the types
func TypeMakeReady() (err error) {
	for _, t := range delayedReady {
		err = t.Ready()
		if err != nil {
			return fmt.Errorf("error initialising go type %s: %v", t.Name, err)
		}
	}

	delayedReady = nil
	return nil
}

func init() {
	err := TypeMakeReady()
	if err != nil {
		log.Fatal(err)
	}
}
