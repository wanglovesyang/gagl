package lex

import (
	"fmt"
	//"math"
	"sort"
	//"strings"
)

const (
	StatusPass = iota
	StatusKeep
	StatusEnd
)

var ErrNothingMatch = fmt.Errorf("Nothing matches in flex tree")

type NodeTag struct {
	BackOffset  int32
	FrontOffset int32
	Tag         int32
	Name        string
}

type FlexNode struct {
	Status   int32
	Tag      NodeTag
	SubRunes []rune
	Children []*FlexNode
}

func (n *FlexNode) Empty() bool {
	return len(n.Children) == 0 || len(n.SubRunes) == 0 || len(n.SubRunes) != len(n.Children)
}

func (n *FlexNode) Insert(r rune, status int32, tag *NodeTag) (ret *FlexNode, reterr error) {
	defer func() {
		if reterr == nil {
			if n.Status == StatusEnd {
				n.Status = StatusKeep
			}

			if tag != nil {
				ret.Tag = *tag
			}
		}
	}()

	ret = &FlexNode{Status: status}
	if len(n.SubRunes) == 0 {
		n.SubRunes = []rune{r}
		n.Children = []*FlexNode{ret}
	} else {
		ind := sort.Search(len(n.SubRunes), func(i int) bool {
			return n.SubRunes[i] >= r
		})

		if ind >= 0 && ind < len(n.SubRunes) {
			if n.SubRunes[ind] == r {
				reterr = fmt.Errorf("Rune %v already exist in SubRunes", r)
				return
			}
		}

		if ind == 0 {
			n.SubRunes = append([]rune{r}, n.SubRunes...)
			n.Children = append([]*FlexNode{ret}, n.Children...)
		} else if ind >= len(n.SubRunes) {
			n.SubRunes = append(n.SubRunes, r)
			n.Children = append(n.Children, ret)
		} else {
			/*nl := make([]rune, len(n.SubRunes)+1)
			copy(nl[0:ind], n.SubRunes[0:ind])
			nl[ind] = r
			copy(nl[ind+1:], n.SubRunes[ind:])
			n.SubRunes = nl

			nll := make([]*FlexNode, len(n.Children)+1)
			copy(nll[0:ind], n.Children[0:ind])
			nll[ind] = ret
			copy(nll[ind+1:], n.Children[ind:])
			n.Children = nll*/
			n.SubRunes = append(n.SubRunes, r)
			copy(n.SubRunes[ind+1:len(n.SubRunes)], n.SubRunes[ind:len(n.SubRunes)-1])
			n.SubRunes[ind] = r

			n.Children = append(n.Children, ret)
			copy(n.Children[ind+1:len(n.Children)], n.Children[ind:len(n.Children)-1])
			n.Children[ind] = ret
		}
	}

	return
}

func (n *FlexNode) FindChild(r rune) (ret *FlexNode, suc bool) {
	if n.Empty() {
		suc = false
		return
	}

	ind := sort.Search(len(n.SubRunes), func(i int) bool {
		return n.SubRunes[i] >= r
	})

	//log.Printf("rune:%v, ind:%d", r, ind)

	if ind >= len(n.SubRunes) {
		suc = false
	} else if n.SubRunes[ind] == r {
		suc = true
		ret = n.Children[ind]
	} else {
		suc = false
	}

	return
}

type FlexTree struct {
	Root *FlexNode
}

func BuildFlexTree(keys map[string]NodeTag) (ret *FlexTree, reterr error) {
	ret = &FlexTree{}
	reterr = ret.buildFromKeySets(keys)
	return
}

func (n *FlexTree) Recognize(src string) (pattern string, tag NodeTag, reterr error) {
	if n.empty() {
		reterr = fmt.Errorf("Fail to recognize when the tree is empty")
	}

	p := n.Root
	cur := 0
	srcc := []rune(src)
	for i := 0; i < len(srcc); i++ {
		r := srcc[i]
		if child, suc := p.FindChild(r); suc {
			//log.Printf("rune %v hit", r)
			if child.Status != StatusPass {
				//log.Printf("rune %v hit last, tag = %v", r, child.Tag)
				cur = i + 1
				tag = child.Tag
			}
			p = child
		} else {
			//log.Printf("rune %v not hit, break", r)
			break
		}
	}

	if cur > 0 {
		pattern = string(srcc[tag.FrontOffset : int32(cur)-tag.BackOffset])
	} else {
		reterr = ErrNothingMatch
	}

	return
}

func (n *FlexTree) empty() bool {
	if n.Root == nil {
		return true
	}

	return n.Root.Empty()
}

func (n *FlexTree) buildFromKeySets(keys map[string]NodeTag) (reterr error) {
	n.Root = &FlexNode{Status: -1}
	for k, tag := range keys {
		kk := []rune(k)
		if err := n.insert(kk, tag); err != nil {
			reterr = fmt.Errorf("Error occurs in insert, %v", err)
			return
		}
	}

	return
}

func (n *FlexTree) insert(key []rune, tag NodeTag) (reterr error) {
	//log.Printf("Insert: %s -- %v", string(key), key)
	p := n.Root
	for i := 0; i < len(key); i++ {
		r := key[i]
		var c *FlexNode
		var t *NodeTag
		st := int32(StatusPass)
		if i == len(key)-1 {
			st = StatusEnd
			t = &tag
		}

		if child, suc := p.FindChild(r); suc {
			c = child
			if st != StatusPass {
				c.Status = st
			}

			if t != nil {
				c.Tag = *t
			}
		} else {
			if child, err := p.Insert(r, st, t); err != nil {
				reterr = fmt.Errorf("Error orccurs in node insert, %v", err)
				break
			} else {
				/*if p == n.Root {
					log.Printf("rune %c insert to root", r)
					log.Printf("root.SubRunes = %v", p.SubRunes)
				}*/

				c = child
			}
		}

		p = c
	}

	return
}
