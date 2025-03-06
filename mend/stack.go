package mend

type Stack[T comparable] struct {
	value []T
}

func NewStack[T comparable]() *Stack[T] {
	return &Stack[T]{}
}

func (stack *Stack[T]) Length() int {
	return len(stack.value)
}

// Appends the element to stack
func (stack *Stack[T]) Add(value T) *Stack[T] {
	stack.value = append(stack.value, value)
	return stack
}

// Pops the last element and returns it.
// Panics if the stack is empty
func (stack *Stack[T]) Pop() T {
	if stack.Length() == 0 {
		panic("Trying to pop an empty stack")
	}

	item := stack.value[len(stack.value)-1]
	stack.value = stack.value[:len(stack.value)-1]
	return item
}
