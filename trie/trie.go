package trie

import (
	"sync"

	"github.com/Zubayear/sonic/stack"
)

// Node represents a single node in the Trie.
// Each node contains a map of children and a boolean flag indicating
// whether the node marks the end of a valid word.
type Node struct {
	children map[rune]*Node
	isEnd    bool
}

// NewTrieNode creates and returns a new Trie node.
func NewTrieNode() *Node {
	return &Node{make(map[rune]*Node), false}
}

// Trie represents the Trie data structure.
// It contains a root node and a size that tracks the number of words.
type Trie struct {
	root  *Node
	size  int
	mutex sync.RWMutex
}

// NewTrie creates and returns an empty Trie.
func NewTrie() *Trie {
	return &Trie{NewTrieNode(), 0, sync.RWMutex{}}
}

// Size returns the total number of words stored in the Trie.
// Time Complexity: O(1)
func (t *Trie) Size() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.size
}

// IsEmpty checks if the Trie contains any words.
// Time Complexity: O(1)
func (t *Trie) IsEmpty() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.size == 0
}

// Insert adds a word into the Trie.
// If the word already exists, it still marks it as a word end, but does not increment size twice.
// Time Complexity: O(N), where N = length of the word
// Space Complexity: O(N), for new nodes if none exist along the path
func (t *Trie) Insert(word string) {
	if len(word) == 0 {
		return
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()
	current := t.root
	for _, ch := range word {
		if current.children[ch] == nil {
			current.children[ch] = NewTrieNode()
		}
		current = current.children[ch]
	}
	if !current.isEnd {
		current.isEnd = true
		t.size++
	}
}

// Search checks if a word exists in the Trie.
// Returns true if the word is found and is a complete word.
// Time Complexity: O(N), where N = length of the word
func (t *Trie) Search(word string) bool {
	if len(word) == 0 {
		return false
	}
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	current := t.root
	for _, ch := range word {
		if current.children[ch] == nil {
			return false
		}
		current = current.children[ch]
	}
	return current.isEnd
}

// StartsWith checks if any word in the Trie starts with the given prefix.
// Returns true if such a prefix exists, even if it is not a complete word.
// Time Complexity: O(K), where K = length of the prefix
func (t *Trie) StartsWith(prefix string) bool {
	if len(prefix) == 0 {
		return false
	}
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	current := t.root
	for _, ch := range prefix {
		if current.children[ch] == nil {
			return false
		}
		current = current.children[ch]
	}
	return true
}

// dfs performs a depth-first search starting from the given node
// and collects all words that stem from the current prefix.
// Time Complexity: O(M * L), where M = number of words from node, L = average word length
func (t *Trie) dfs(node *Node, prefix string) []string {
	var result []string
	var dfs func(node *Node, prefix string)
	dfs = func(node *Node, prefix string) {
		if node.isEnd {
			result = append(result, prefix)
		}
		for ch, child := range node.children {
			dfs(child, prefix+string(ch))
		}
	}
	dfs(node, prefix)
	return result
}

// findNodeForPrefix returns the node corresponding to the last character of the given prefix.
// If the prefix does not exist in the Trie, it returns nil.
// Time Complexity: O(K), where K = length of the prefix
func (t *Trie) findNodeForPrefix(prefix string) *Node {
	current := t.root
	for _, ch := range prefix {
		if current.children[ch] == nil {
			return nil
		}
		current = current.children[ch]
	}
	return current
}

// GetWordsWithPrefix retrieves all words in the Trie that start with the given prefix.
// Returns an empty slice if the prefix does not exist.
// Time Complexity: O(K + M * L), where K = length of prefix, M = number of matching words
func (t *Trie) GetWordsWithPrefix(prefix string) []string {
	if len(prefix) == 0 {
		return nil
	}
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	var result []string
	current := t.findNodeForPrefix(prefix)
	if current == nil {
		return result
	}
	return t.dfs(current, prefix)
}

// Remove deletes a word from the Trie if it exists.
// Returns true if the word was removed, false if it was not present.
// It also removes unnecessary nodes to keep the Trie compact.
// Time Complexity: O(N), where N = length of the word
// Space Complexity: O(N) for the stack to track nodes
func (t *Trie) Remove(word string) bool {
	if len(word) == 0 {
		return false
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()
	current := t.root
	type Pair struct {
		node *Node
		ch   rune
	}

	s := stack.NewStack[Pair]()
	for _, ch := range word {
		next := current.children[ch]
		if next == nil {
			return false
		}
		_, _ = s.Push(Pair{current, ch})
		current = next
	}
	if !current.isEnd {
		return false
	}
	current.isEnd = false

	for !s.IsEmpty() {
		val, _ := s.Pop()
		parent := val.node
		ch := val.ch
		child := parent.children[ch]
		if len(child.children) == 0 && !child.isEnd {
			delete(parent.children, ch)
		} else {
			break
		}
	}
	t.size--
	return true
}
