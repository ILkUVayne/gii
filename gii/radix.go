package gii

import (
	"strings"
)

type radixNode struct {
	path     string
	fullPath string
	indices  string
	passCnt  int
	children []*radixNode
	end      bool
	handlers HandlersChain
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

func (r *Radix) Search(word string) bool {
	node := r.root.search(word)
	return node != nil && node.end == true && node.fullPath == word
}

func (r *Radix) GetHandles(word string) HandlersChain {
	return r.root.search(word).handlers
}

func (r *Radix) StartWith(prefix string) bool {
	node := r.root.search(prefix)
	return node != nil && strings.HasPrefix(node.fullPath, prefix)
}

func (r *Radix) PassCnt(prefix string) int {
	node := r.root.search(prefix)
	if node == nil || !strings.HasPrefix(node.fullPath, prefix) {
		return 0
	}
	return node.passCnt
}

func (r *Radix) Insert(word string, Handlers HandlersChain) {
	if r.Search(word) {
		return
	}
	r.root.insert(word, Handlers)
}

func (r *Radix) Del(word string) bool {
	if !r.Search(word) {
		return false
	}
	return r.root.del(word)
}

// ----------------------- Radix Node ---------------------------------

func (rn *radixNode) insert(word string, handlers HandlersChain) {
	fullPath := word

	// 空树，直接添加
	if rn.fullPath == "" && len(rn.children) == 0 {
		rn.insertWord(word, fullPath, handlers)
		return
	}
walk:
	for {
		// 获取公共前缀长度 commonLen
		cl := rn.commonPrefixLen(word, rn.path)
		// 公共前缀长度大于0时，必定经过该节点
		if cl > 0 {
			rn.passCnt++
		}
		// 公共长度小于path,拆分path公共前缀
		if cl < len(rn.path) {
			// 创建需要拆分的子节点（非公共前缀）
			children := &radixNode{
				path:     rn.path[cl:],
				fullPath: rn.fullPath,
				end:      rn.end,
				indices:  rn.indices,
				passCnt:  rn.passCnt - 1,
				children: rn.children,
				handlers: rn.handlers,
			}
			// 调整父节点,续接上拆分的子节点
			rn.indices = string(rn.path[cl])
			rn.fullPath = rn.fullPath[:len(rn.fullPath)-(len(rn.path)-cl)]
			rn.path = rn.path[:cl]
			rn.end = false
			rn.children = []*radixNode{children}
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
					continue walk
				}
			}
			// 没有公共前缀了
			rn.indices += string(c)
			children := &radixNode{}
			children.insertWord(word, fullPath, handlers)
			rn.children = append(rn.children, children)
			return
		}
		// 刚好匹配path
		rn.end = true
		rn.handlers = handlers
		return
	}
}

func (rn *radixNode) search(word string) *radixNode {
walk:
	for {
		prefix := rn.path
		if len(word) > len(prefix) {
			// 去除公共前缀
			word = word[len(prefix):]
			// 获取首字母
			c := word[0]
			// 遍历首字母集，确定子节点
			for i := 0; i < len(rn.indices); i++ {
				if rn.indices[i] == c {
					rn = rn.children[i]
					continue walk
				}
			}
		}
		// 和当前节点精准匹配上了
		if word == prefix {
			return rn
		}
		// 走到这里意味着 len(word) <= len(prefix) && word != prefix
		return rn
	}
}

func (rn *radixNode) del(word string) bool {
	// root 直接精准命中了
	if rn.fullPath == word {
		// 如果一个孩子都没有
		if len(rn.indices) == 0 {
			rn.path = ""
			rn.fullPath = ""
			rn.end = false
			rn.passCnt = 0
			return true
		}

		// 如果只有一个孩子
		if len(rn.indices) == 1 {
			rn.children[0].path = rn.path + rn.children[0].path
			*rn = *rn.children[0]
			return true
		}

		// 如果有多个孩子
		rn.passCnt--
		rn.end = false
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
		for i := 0; i < len(move.indices); i++ {
			if move.indices[i] != c {
				continue
			}

			// 精准命中但是他仍有后继节点
			if move.children[i].path == word && move.children[i].passCnt > 1 {
				move.children[i].end = false
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
			if !move.end && len(move.indices) == 1 {
				// 合并自己与唯一的孩子
				move.path += move.children[0].path
				move.fullPath = move.children[0].fullPath
				move.end = move.children[0].end
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
	rn.path, rn.fullPath = path, fullPath
	rn.end = true
	rn.passCnt++
	rn.handlers = handlers
}
