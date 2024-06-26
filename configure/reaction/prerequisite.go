package reaction

import (
	cb "github.com/AleckDarcy/ContextBus/proto"
)

func NewPrerequisiteTree(tree *cb.PrerequisiteTree) *PrerequisiteTree {
	t := &PrerequisiteTree{
		PrerequisiteTree: tree,
	}

	t.Index = map[string]*PrerequisiteNode{}
	for _, node := range t.Nodes {
		if node.Type == cb.PrerequisiteNodeType_PrerequisiteMessage_ {
			t.Index[node.Message.Name] = (*PrerequisiteNode)(node)
		}
	}

	return t
}
