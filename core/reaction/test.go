package reaction

import (
	cb "github.com/AleckDarcy/ContextBus/proto"
)

// Tree0 (EventA) && (EventB = 1)
var Tree0 = NewPrerequisiteTree(cb.Test_PrerequisiteTree0)

// Tree1 ((EventA) && (EventB = 1)) || (1 < EventC < 4)
var Tree1 = NewPrerequisiteTree(cb.Test_PrerequisiteTree1)
