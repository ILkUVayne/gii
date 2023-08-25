package tree

import (
	"reflect"
	"testing"
)

func TestNewTrie(t *testing.T) {
	node := NewTrie()
	ok := reflect.DeepEqual(node, &Trie{root: &trieNode{}})
	if !ok {
		t.Error("NewTrie() cannot create trie")
	}
}

func TestTrie_Insert(t *testing.T) {
	root := NewTrie()
	insertWord := "hello"
	insertWord2 := "hellow"
	root.Insert(insertWord)
	root.Insert(insertWord2)
	if !root.Search(insertWord) {
		t.Errorf("%s cannot be insert into trie", insertWord)
	}
	if !root.Search(insertWord2) {
		t.Errorf("%s cannot be insert into trie", insertWord2)
	}
}

func TestTrie_Del(t *testing.T) {
	root := NewTrie()
	insertWord := "hello"
	insertWord2 := "hellow"
	root.Insert(insertWord)
	root.Insert(insertWord2)
	res := root.Del(insertWord2)
	if !res {
		t.Errorf("del result %v", res)
	}
	if root.Search(insertWord2) {
		t.Errorf("cannot del %s", insertWord2)
	}
}

func TestTrie_StartsWith(t *testing.T) {
	root := NewTrie()
	insertWord := "hello"
	insertWord2 := "hefsfe"
	root.Insert(insertWord)
	root.Insert(insertWord2)
	if !root.StartsWith("he") {
		t.Errorf("cannot be match by prfix %s", "he")
	}

	if !root.StartsWith("hef") {
		t.Errorf("cannot be match by prfix %s", "hef")
	}
}

func TestTrie_PassCnt(t *testing.T) {
	root := NewTrie()
	insertWord := "hello"
	insertWord2 := "hefsfe"
	insertWord3 := "helsfe"
	insertWord4 := "helsada"
	root.Insert(insertWord)
	root.Insert(insertWord2)
	root.Insert(insertWord3)
	root.Insert(insertWord4)
	if root.PassCnt("hell") != 1 {
		t.Errorf("prefix \"hell\" PassCnt != %d", 1)
	}
	if root.PassCnt("hel") != 3 {
		t.Errorf("prefix \"hel\" PassCnt != %d", 3)
	}
	if root.PassCnt("hels") != 2 {
		t.Errorf("prefix \"hels\" PassCnt != %d", 2)
	}
}
