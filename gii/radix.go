package gii

import (
	"strings"
)

const (
	static int = iota
	param
)

type radixNode struct {
	path      string
	fullPath  string
	indices   string
	passCnt   int
	children  []*radixNode
	wildChild bool
	handlers  HandlersChain
}

type Radix struct {
	root *radixNode
}

func NewRadix() *Radix {
	return &Radix{
		root: &radixNode{},
	}
}

type methodTree struct {
	method string
	root   *Radix
}

type methodTrees []methodTree

// ----------------------- Radix ---------------------------------

func (r *Radix) Search(word string, mode int) bool {
	node := r.root.search(word, mode)
	//return node != nil && node.end == true && node.fullPath == word
	return node != nil && node.handlers != nil
}

func (r *Radix) GetHandles(word string) HandlersChain {
	return r.root.search(word, param).handlers
}

func (r *Radix) StartWith(prefix string) bool {
	node := r.root.search(prefix, static)
	return node != nil && strings.HasPrefix(node.fullPath, prefix)
}

func (r *Radix) PassCnt(prefix string) int {
	node := r.root.search(prefix, static)
	if node == nil || !strings.HasPrefix(node.fullPath, prefix) {
		return 0
	}
	return node.passCnt
}

func (r *Radix) Insert(word string, Handlers HandlersChain) {
	r.root.insert(word, Handlers)
}

func (r *Radix) Del(word string) bool {
	if !r.Search(word, static) {
		return false
	}
	return r.root.del(word)
}

// ----------------------- Radix Node ---------------------------------

func (rn *radixNode) addChild(child *radixNode) {
	if rn.wildChild && len(rn.children) > 0 {
		wildChild := rn.children[len(rn.children)-1]
		rn.children = append(rn.children[:len(rn.children)-1], child, wildChild)
		return
	}
	rn.children = append(rn.children, child)
}

func (rn *radixNode) insert(word string, handlers HandlersChain) {
	fullPath := word
	rn.passCnt++

	// 空树，直接添加
	if rn.fullPath == "" && len(rn.children) == 0 {
		rn.insertWord(word, fullPath, handlers)
		return
	}
walk:
	for {
		// 获取公共前缀长度 commonLen
		cl := rn.commonPrefixLen(word, rn.path)
		// 公共长度小于path,拆分path公共前缀
		if cl < len(rn.path) {
			// 创建需要拆分的子节点（非公共前缀）
			children := &radixNode{
				path:      rn.path[cl:],
				fullPath:  rn.fullPath,
				indices:   rn.indices,
				passCnt:   rn.passCnt - 1,
				children:  rn.children,
				handlers:  rn.handlers,
				wildChild: rn.wildChild,
			}

			// 续接上拆分的子节点,调整父节点
			rn.children = []*radixNode{children}
			rn.indices = string(rn.path[cl])
			rn.fullPath = rn.fullPath[:len(rn.fullPath)-(len(rn.path)-cl)]
			rn.path = rn.path[:cl]
			rn.wildChild = false
			rn.handlers = nil
		}
		// 公共长度小于word
		if cl < len(word) {
			// 去除公共前缀
			word = word[cl:]
			// word 首字母
			c := word[0]
			// 判断是否还存在公共前缀
			for i := 0; i < len(rn.indices); i++ {
				if rn.indices[i] == c {
					rn = rn.children[i]
					rn.passCnt++
					continue walk
				}
			}

			if c == ':' && rn.wildChild {
				rn = rn.children[len(rn.children)-1]
				if len(word) > len(rn.path) && word[:len(rn.path)] == rn.path && word[len(rn.path)] == '/' {
					rn.passCnt++
					continue walk
				}

				panic("'" + strings.SplitN(word, "/", 2)[0] +
					"' in new path '" + fullPath +
					"' conflicts with existing wildcard '" + rn.path +
					"' in existing prefix '" + fullPath[:strings.Index(fullPath, word)] + rn.path +
					"'")
			}

			// 没有公共前缀了
			if c != ':' {
				rn.indices += string(c)
				child := &radixNode{passCnt: 1}
				rn.addChild(child)
				rn = child
			}

			rn.insertWord(word, fullPath, handlers)
			return
		}
		// 刚好匹配path
		rn.handlers = handlers
		return
	}
}

func (rn *radixNode) search(word string, mode int) *radixNode {
walk:
	for {
		prefix := rn.path
		if len(word) > len(prefix) {
			// 前缀不匹配时，表示路由不存在
			if word[:len(prefix)] != prefix {
				return nil
			}
			// 去除公共前缀
			word = word[len(prefix):]
			// 获取首字母
			c := word[0]
			// 遍历首字母集，确定子节点
			for i := 0; i < len(rn.indices); i++ {
				if rn.indices[i] == c {
					// 判断是否可能为动态节点
					if mode == param && rn.wildChild {
						if len(word) < len(rn.children[i].path) {
							break
						}

						if len(word) == len(rn.children[i].path) && rn.children[i].path != word {
							break
						}

						if len(word) > len(rn.children[i].path) {
							if strings.IndexAny(word, "/") == -1 {
								return rn.children[len(rn.children)-1]
							}

							if rn.children[i].path != word[:len(rn.children[i].path)] {
								break
							}
						}
					}
					// 静态节点
					rn = rn.children[i]
					continue walk
				}
			}
			if mode == static && c == ':' && rn.wildChild && len(rn.children) > 0 {
				rn = rn.children[len(rn.children)-1]
				continue walk
			}
			// 动态路由节点
			if mode == param && rn.wildChild {
				wordSplit := strings.SplitN(word, "/", 2)
				rn = rn.children[len(rn.children)-1]
				if len(wordSplit) == 1 {
					return rn
				}
				word = rn.path + "/" + wordSplit[1]
				continue walk
			}
		}
		// 和当前节点精准匹配上了
		if word == prefix {
			return rn
		}
		// 走到这里意味着 len(word) <= len(prefix) && word != prefix
		return nil
	}
}

func (rn *radixNode) del(word string) bool {
	// root 直接精准命中了
	if rn.fullPath == word {
		// 如果一个孩子都没有
		if len(rn.children) == 0 {
			rn.path = ""
			rn.fullPath = ""
			rn.handlers = nil
			rn.passCnt = 0
			return true
		}

		// 如果只有一个孩子
		if len(rn.children) == 1 {
			rn.children[0].path = rn.path + rn.children[0].path
			*rn = *rn.children[0]
			return true
		}

		// 如果有多个孩子
		rn.passCnt--
		rn.handlers = nil
		return true
	}

	// 确定 word 存在的情况下
	move := rn
	// root 单独作为一个分支处理
	// 其他情况下，需要对孩子进行处理
walk:
	for {
		move.passCnt--
		prefix := move.path
		word = word[len(prefix):]
		c := word[0]

		if c == ':' && move.wildChild {
			if move.children[len(move.children)-1].passCnt == 1 {
				move.children = move.children[:len(move.children)-1]
				return true
			}
			move = move.children[len(move.children)-1]
			continue walk
		}

		for i := 0; i < len(move.indices); i++ {
			if move.indices[i] != c {
				continue
			}

			// 精准命中但是他仍有后继节点
			if move.children[i].path == word && move.children[i].passCnt > 1 {
				move.children[i].handlers = nil
				move.children[i].passCnt--
				return true
			}

			// 找到对应的 child 了
			// 如果说后继节点的 passCnt = 1，直接干掉
			if move.children[i].passCnt > 1 {
				move = move.children[i]
				continue walk
			}

			move.children = append(move.children[:i], move.children[i+1:]...)
			move.indices = move.indices[:i] + move.indices[i+1:]
			// 如果干掉一个孩子后，发现只有一个孩子了，并且自身 end 为 false 则需要进行合并
			if move.handlers == nil && len(move.indices) == 1 {
				// 合并自己与唯一的孩子
				move.path += move.children[0].path
				move.fullPath = move.children[0].fullPath
				move.handlers = move.children[0].handlers
				move.indices = move.children[0].indices
				move.children = move.children[0].children
			}
			return true
		}
	}
}

func (rn *radixNode) commonPrefixLen(word, path string) int {
	commonLen := 0
	for commonLen < len(word) && commonLen < len(path) && word[commonLen] == path[commonLen] {
		commonLen++
	}
	return commonLen
}

func (rn *radixNode) insertWord(path, fullPath string, handlers HandlersChain) {
	prefixLen := len(fullPath[:len(fullPath)-len(path)])
	for {
		wildCard, i, valid := findWildCard(path)
		if i < 0 {
			break
		}
		if !valid {
			panic("only one wildcard per path segment is allowed, has: '" +
				wildCard + "' in path '" + fullPath + "'")
		}
		if len(wildCard) < 2 {
			panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
		}

		if i > 0 {
			rn.path = path[:i]
			prefixLen += len(rn.path)
			rn.fullPath = fullPath[:prefixLen]
		}

		prefixLen += len(wildCard)
		child := &radixNode{
			path:     wildCard,
			fullPath: fullPath[:prefixLen],
			passCnt:  1,
		}
		rn.addChild(child)
		rn.wildChild = true
		rn = child

		path = path[i:]
		if len(wildCard) < len(path) {
			path = path[len(wildCard):]
			child := &radixNode{passCnt: 1}
			rn.addChild(child)
			rn = child
			continue
		}
		rn.handlers = handlers
		return
	}

	// no wildCard
	rn.path, rn.fullPath = path, fullPath
	rn.handlers = handlers
}

func findWildCard(word string) (wildCard string, i int, valid bool) {
	for start, v := range []byte(word) {
		if v != ':' {
			continue
		}

		valid = true
		for end, v := range []byte(word[start+1:]) {
			switch v {
			case '/':
				return word[start : start+1+end], start, valid
			case ':':
				valid = false
			}
		}
		return word[start:], start, valid
	}
	return "", -1, false
}
