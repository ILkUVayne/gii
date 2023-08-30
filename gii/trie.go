package gii

type trieNode struct {
	passCnt int
	next    [26]*trieNode // 目前只支持a-z
	end     bool
}

type Trie struct {
	root *trieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: &trieNode{},
	}
}

func (t *Trie) Search(word string) bool {
	node := t.search(word)
	return node != nil && node.end
}

func (t *Trie) StartsWith(prefix string) bool {
	node := t.search(prefix)
	return node != nil
}

func (t *Trie) PassCnt(prefix string) int {
	if prefix == "" {
		return t.root.passCnt
	}
	node := t.search(prefix)
	if node == nil {
		return 0
	}
	return node.passCnt
}

func (t *Trie) Insert(word string) {
	// 判断是否已存在，存在则直接返回
	if t.Search(word) {
		return
	}
	move := t.root
	for _, v := range word {
		index := getIndex(v)
		// 如果节点未赋值(nil),直接赋值
		if move.next[index] == nil {
			move.next[index] = &trieNode{}
		}
		move.passCnt++
		move = move.next[index]
	}
	move.end = true
}

func (t *Trie) Del(word string) bool {
	// 不存在，直接返回
	if !t.Search(word) {
		return false
	}
	move := t.root
	for _, v := range word {
		index := getIndex(v)
		move.passCnt--
		// 如果当前节点passCnt=0 表示没有后继节点，直接重置该节点
		if move.passCnt == 0 {
			move.next[index] = nil
			return true
		}
		move = move.next[index]
	}
	move.end = false
	return true
}

func (t *Trie) search(word string) *trieNode {
	move := t.root
	for _, v := range word {
		index := getIndex(v)
		if move.next[index] == nil {
			return nil
		}
		move = move.next[index]
	}
	return move
}

func getIndex(number int32) int32 {
	index := number - 'a'
	if index < 0 || index > 25 {
		panic("trie only support character a-z")
	}
	return index
}
