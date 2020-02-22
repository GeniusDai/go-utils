package utils

type TrieNode struct {
	children map[rune]*TrieNode
	isWord   bool
}

type TrieRoot TrieNode

func (t *TrieRoot) Init() {
	t.children = make(map[rune]*TrieNode)
}

func (t *TrieRoot) Insert(s string) {
	var node *TrieNode = &TrieNode{
		children: t.children,
		isWord:   false,
	}
	for _, b := range s {
		_, ok := node.children[b]
		if !ok {
			node.children[b] = &TrieNode{
				children: make(map[rune]*TrieNode),
				isWord:   false,
			}
		}
		node = node.children[b]
	}
	node.isWord = true
	return
}

func (t *TrieRoot) HasWord(s string) bool {
	node := t.findNode(s)
	if node == nil {
		return false
	}
	return node.isWord
}

func (t *TrieRoot) HasPrefix(s string) bool {
	node := t.findNode(s)
	return node != nil
}

func (t *TrieRoot) Delete(s string) {
	// First update isWord field
	if node := t.findNode(s); node != nil {
		node.isWord = false
	} else {
		return
	}
	// Then delete children field if it's the only path.
	// Previous operation ensures existence of the path.
	var node *TrieNode
	for i := len(s) - 1; i >= 0; i-- {
		node = t.findNode(s[:i])
		if len(node.children[rune(s[i])].children) == 0 {
			delete(node.children, rune(s[i]))
		} else {
			break
		}
	}
}

func (t *TrieRoot) findNode(s string) *TrieNode {
	var node *TrieNode = &TrieNode{
		children: t.children,
		isWord:   false,
	}
	for _, b := range s {
		children, ok := node.children[b]
		if !ok {
			return nil
		}
		node = children
	}
	return node
}
