package CharBalancer

import (
	"errors"
	"sync"
)

// A CharacterGroup is a structure that contains the opening and closing Characters of a common Character.
type CharacterGroup struct {
	Opening rune
	Closing rune
}

// A CharacterEntry contains a Character and the CharacterGroup it belongs to.
type CharacterEntry struct {
	Char  rune
	Group *CharacterGroup
}

// A CharacterGroupStack is a thread-safe stack implementation for CharacterGroup.
type CharacterGroupStack struct {
	mux   *sync.Mutex
	array []*CharacterEntry
}

// NewCharacterGroup returns a CharacterGroup pointer with the provided
// opening and closing Characters (of type rune).
func NewCharacterGroup(opening, closing rune) *CharacterGroup {
	return &CharacterGroup{Opening: opening, Closing: closing}
}

// NewCharacterEntry returns a CharacterEntry with the provided Character and CharacterGroup.
func NewCharacterEntry(Char rune, Group **CharacterGroup) *CharacterEntry {
	return &CharacterEntry{Char: Char, Group: *Group}
}

// NewCharacterGroupStack returns a new CharacterGroupStack.
func NewCharacterGroupStack() *CharacterGroupStack {
	return &CharacterGroupStack{mux: &sync.Mutex{}, array: make([]*CharacterEntry, 0)}
}

// IsClosingCharacter returns true if the Character in the CharacterEntry is the closing Character
// in the corresponding CharacterGroup.
func (entry *CharacterEntry) IsClosingCharacter() bool {
	return entry.Char == entry.Group.Closing
}

// IsBalanced returns true if all the Characters are balanced, false otherwise.
func (stack *CharacterGroupStack) IsBalanced() bool {
	return len(stack.array) == 0
}

// Push adds the provided CharacterEntry to the stack.
func (stack *CharacterGroupStack) Push(entry **CharacterEntry) {
	if stack.array == nil {
		stack.mux.Lock()
		stack.array = make([]*CharacterEntry, 0)
		stack.mux.Unlock()
	}

	length := len(stack.array)

	// Check if we have a pair of opening/closing Characters.
	if (*entry).IsClosingCharacter() {
		prev := stack.array[length-1]

		// Matching!
		if length >= 1 && !prev.IsClosingCharacter() && prev.Group == (*entry).Group {
			_, err := stack.Pop()

			if err != nil {
				panic(err)
			}

			return
		}
	}

	stack.mux.Lock()
	defer stack.mux.Unlock()

	stack.array = append(stack.array, *entry)
}

// Pop removes the latest element from the Stack and returns it.
func (stack *CharacterGroupStack) Pop() (*CharacterEntry, error) {
	stack.mux.Lock()
	defer stack.mux.Unlock()

	if stack.array == nil {
		stack.array = make([]*CharacterEntry, 0)
	}

	length := len(stack.array)

	if length == 0 {
		return nil, errors.New("underlying array in stack is empty")
	}

	// Get last value
	poppedValue := stack.array[length-1]

	// Set array to slice containing all except last value
	if length <= 1 {
		stack.array = make([]*CharacterEntry, 0)
	} else {
		stack.array = stack.array[:length-1]
	}

	return poppedValue, nil
}
