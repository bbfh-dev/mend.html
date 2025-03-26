package tags

// Represents a regular HTML paired node
type CustomIfNode struct {
	*pairedNode
	Value  string
	Expect bool
}

func NewCustomIfNode(value string, expect bool) *CustomIfNode {
	return &CustomIfNode{
		pairedNode: newPairedNode(),
		Value:      value,
		Expect:     expect,
	}
}

func (node *CustomIfNode) check() bool {
	return (node.Value == "true") == node.Expect
}

func (node *CustomIfNode) Render(out writer, indent int) {
	if node.check() {
		node.renderMinimal(out, indent)
	}
}

func (node *CustomIfNode) Visible() bool {
	return node.check()
}
