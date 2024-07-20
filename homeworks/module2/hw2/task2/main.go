package main

import "log"

type Node struct {
	size   int
	key    [3]int
	first  *Node
	second *Node
	third  *Node
	fourth *Node
	parent *Node
}

func NewNode(k int) *Node {
	return &Node{
		size:   1,
		key:    [3]int{k, 0, 0},
		first:  nil,
		second: nil,
		third:  nil,
		fourth: nil,
		parent: nil,
	}
}

func (n *Node) swap(x, y *int) {
	*x, *y = *y, *x
}

func (n *Node) insertToNode(k int) {
	n.key[n.size] = k
	n.size++

	if n.size == 1 {
		return
	}

	if n.size == 2 {
		if n.key[0] > n.key[1] {
			n.swap(&n.key[0], &n.key[1])
		}
	}

	if n.size == 3 {
		if n.key[0] > n.key[1] {
			n.swap(&n.key[0], &n.key[1])
		}
		if n.key[0] > n.key[2] {
			n.swap(&n.key[0], &n.key[2])
		}
		if n.key[1] > n.key[2] {
			n.swap(&n.key[1], &n.key[2])
		}
	}
}

func (n *Node) removeFromNode(k int) {
	if n.size >= 1 && n.key[0] == k {
		n.key[0] = n.key[1]
		n.key[1] = n.key[2]
		n.size--
	} else if n.size == 2 && n.key[1] == k {
		n.key[1] = n.key[2]
		n.size--
	}
}

func (n *Node) Insert(k int) *Node {
	if n == nil {
		return NewNode(k)
	}

	if n.first == nil && n.second == nil && n.third == nil {
		n.insertToNode(k)
	} else if k <= n.key[0] {
		n.first = n.first.Insert(k)
	} else if n.size == 1 || (n.size == 2 && k <= n.key[1]) {
		n.second = n.second.Insert(k)
	} else {
		n.third = n.third.Insert(k)
	}

	return n.Split(n)
}

func (n *Node) Split(item *Node) *Node {
	if item.size < 3 {
		return item
	}

	x := NewNode(item.key[0])
	y := NewNode(item.key[2])

	x.first = item.first
	x.second = item.second
	y.first = item.third
	y.second = item.fourth

	if x.first != nil {
		x.first.parent = x
	}
	if x.second != nil {
		x.second.parent = x
	}
	if y.first != nil {
		y.first.parent = y
	}
	if y.second != nil {
		y.second.parent = y
	}

	if item.parent != nil {
		item.parent.insertToNode(item.key[1])

		if item.parent.first == item {
			item.parent.first = nil
		} else if item.parent.second == item {
			item.parent.second = nil
		} else if item.parent.third == item {
			item.parent.third = nil
		}

		if item.parent.first == nil {
			item.parent.fourth = item.parent.third
			item.parent.third = item.parent.second
			item.parent.second = y
			item.parent.first = x
		} else if item.parent.second == nil {
			item.parent.fourth = item.parent.third
			item.parent.third = y
			item.parent.second = x
		} else {
			item.parent.fourth = y
			item.parent.third = x
		}

		tmp := item.parent
		item = nil
		return tmp
	} else {
		x.parent = item
		y.parent = item

		n.key[0] = item.key[1]
		n.first = x
		n.second = y
		n.third = nil
		n.fourth = nil
		n.parent = nil
		n.size = 1

		return item
	}
}

func (n *Node) Search(k int) *Node {
	if n == nil {
		return nil
	}

	found := false
	for i := 0; i < n.size; i++ {
		if n.key[i] == k {
			found = true
		}
	}

	if found {
		return n
	} else if k < n.key[0] {
		return n.first.Search(k)
	} else if n.size == 1 || (n.size == 2 && k < n.key[1]) {
		return n.second.Search(k)
	} else {
		return n.third.Search(k)
	}
}

func (n *Node) SearchMin() *Node {
	if n == nil {
		return nil
	}
	if n.first == nil {
		return n
	}
	return n.first.SearchMin()
}

func (n *Node) Remove(k int) *Node {
	item := n.Search(k)
	if item == nil {
		return n
	}

	var min *Node
	if item.key[0] == k {
		min = item.second.SearchMin()
	} else {
		min = item.third.SearchMin()
	}

	if min != nil {
		var z *int
		if k == item.key[0] {
			z = &item.key[0]
		} else {
			z = &item.key[1]
		}
		item.swap(z, &min.key[0])
		item = min
	}

	item.removeFromNode(k)
	return n.Fix(item)
}

func (n *Node) Fix(leaf *Node) *Node {
	if leaf.size == 0 && leaf.parent == nil {
		leaf = nil
		return nil
	}
	if leaf.size != 0 {
		if leaf.parent != nil {
			return n.Fix(leaf.parent)
		}
		return leaf
	}

	parent := leaf.parent
	if parent.first.size == 2 || parent.second.size == 2 || parent.size == 2 {
		leaf = n.Redistribute(leaf)
	} else if parent.size == 2 && parent.third.size == 2 {
		leaf = n.Redistribute(leaf)
	} else {
		leaf = n.Merge(leaf)
	}

	return n.Fix(leaf)
}

func (n *Node) Merge(leaf *Node) *Node {
	parent := leaf.parent

	if parent.first == leaf {
		parent.second.insertToNode(parent.key[0])
		parent.second.third = parent.second.second
		parent.second.second = parent.second.first

		if leaf.first != nil {
			parent.second.first = leaf.first
		} else if leaf.second != nil {
			parent.second.first = leaf.second
		}

		if parent.second.first != nil {
			parent.second.first.parent = parent.second
		}

		parent.removeFromNode(parent.key[0])
		leaf = nil
		parent.first = nil
	} else if parent.second == leaf {
		parent.first.insertToNode(parent.key[0])

		if leaf.first != nil {
			parent.first.third = leaf.first
		} else if leaf.second != nil {
			parent.first.third = leaf.second
		}

		if parent.first.third != nil {
			parent.first.third.parent = parent.first
		}

		parent.removeFromNode(parent.key[0])
		leaf = nil
		parent.second = nil
	}

	if parent.parent == nil {
		var tmp *Node
		if parent.first != nil {
			tmp = parent.first
		} else {
			tmp = parent.second
		}
		tmp.parent = nil
		parent = nil
		return tmp
	}
	return parent
}

func (n *Node) Redistribute(leaf *Node) *Node {
	parent := leaf.parent
	first := parent.first
	second := parent.second
	third := parent.third

	if parent.size == 2 && first.size < 2 && second.size < 2 && third.size < 2 {
		if first == leaf {
			parent.first = parent.second
			parent.second = parent.third
			parent.third = nil
			parent.first.insertToNode(parent.key[0])
			parent.first.third = parent.first.second
			parent.first.second = parent.first.first

			if leaf.first != nil {
				parent.first.first = leaf.first
			} else if leaf.second != nil {
				parent.first.first = leaf.second
			}

			if parent.first.first != nil {
				parent.first.first.parent = parent.first
			}

			parent.removeFromNode(parent.key[0])
			leaf = nil
		} else if second == leaf {
			first.insertToNode(parent.key[0])
			parent.removeFromNode(parent.key[0])
			if leaf.first != nil {
				first.third = leaf.first
			} else if leaf.second != nil {
				first.third = leaf.second
			}

			if first.third != nil {
				first.third.parent = first
			}

			parent.second = parent.third
			parent.third = nil
			leaf = nil
		} else if third == leaf {
			second.insertToNode(parent.key[1])
			parent.third = nil
			parent.removeFromNode(parent.key[1])
			if leaf.first != nil {
				second.third = leaf.first
			} else if leaf.second != nil {
				second.third = leaf.second
			}

			if second.third != nil {
				second.third.parent = second
			}
			leaf = nil
		}
	} else if parent.size == 2 && (first.size == 2 || second.size == 2 || third.size == 2) {
		if third == leaf {
			if leaf.first != nil {
				leaf.second = leaf.first
				leaf.first = nil
			}

			leaf.insertToNode(parent.key[1])
			if second.size == 2 {
				parent.key[1] = second.key[1]
				second.removeFromNode(second.key[1])
				leaf.first = second.third
				second.third = nil
				if leaf.first != nil {
					leaf.first.parent = leaf
				}
			} else if first.size == 2 {
				parent.key[1] = second.key[0]
				leaf.first = second.second
				second.second = second.first
				if leaf.first != nil {
					leaf.first.parent = leaf
				}

				second.key[0] = parent.key[0]
				parent.key[0] = first.key[1]
				first.removeFromNode(first.key[1])
				second.first = first.third
				if second.first != nil {
					second.first.parent = second
				}
				first.third = nil
			}
		} else if second == leaf {
			if third.size == 2 {
				if leaf.first == nil {
					leaf.first = leaf.second
					leaf.second = nil
				}
				second.insertToNode(parent.key[1])
				parent.key[1] = third.key[0]
				third.removeFromNode(third.key[0])
				second.second = third.first
				if second.second != nil {
					second.second.parent = second
				}
				third.first = third.second
				third.second = third.third
				third.third = nil
			} else if first.size == 2 {
				if leaf.second == nil {
					leaf.second = leaf.first
					leaf.first = nil
				}
				second.insertToNode(parent.key[0])
				parent.key[0] = first.key[1]
				first.removeFromNode(first.key[1])
				second.first = first.third
				if second.first != nil {
					second.first.parent = second
				}
				first.third = nil
			}
		} else if first == leaf {
			if leaf.first == nil {
				leaf.first = leaf.second
				leaf.second = nil
			}
			first.insertToNode(parent.key[0])
			if second.size == 2 {
				parent.key[0] = second.key[0]
				second.removeFromNode(second.key[0])
				first.second = second.first
				if first.second != nil {
					first.second.parent = first
				}
				second.first = second.second
				second.second = second.third
				second.third = nil
			} else if third.size == 2 {
				parent.key[0] = second.key[0]
				second.key[0] = parent.key[1]
				parent.key[1] = third.key[0]
				third.removeFromNode(third.key[0])
				first.second = second.first
				if first.second != nil {
					first.second.parent = first
				}
				second.first = second.second
				second.second = third.first
				if second.second != nil {
					second.second.parent = second
				}
				third.first = third.second
				third.second = third.third
				third.third = nil
			}
		}
	} else if parent.size == 1 {
		leaf.insertToNode(parent.key[0])

		if first == leaf && second.size == 2 {
			parent.key[0] = second.key[0]
			second.removeFromNode(second.key[0])

			if leaf.first == nil {
				leaf.first = leaf.second
			}

			leaf.second = second.first
			second.first = second.second
			second.second = second.third
			second.third = nil
			if leaf.second != nil {
				leaf.second.parent = leaf
			}
		} else if second == leaf && first.size == 2 {
			parent.key[0] = first.key[1]
			first.removeFromNode(first.key[1])

			if leaf.second == nil {
				leaf.second = leaf.first
			}

			leaf.first = first.third
			first.third = nil
			if leaf.first != nil {
				leaf.first.parent = leaf
			}
		}
	}
	return parent
}

func main() {
	var root *Node

	root = root.Insert(10)
	if root == nil || root.key[0] != 10 {
		log.Fatalf("Ошибка вставки узла %v\n", root.key[0])
	}
	root = root.Insert(20)
	root = root.Insert(5)
	root = root.Insert(6)
	root = root.Insert(15)

	result := root.Search(10)
	if result == nil || result.key[0] != 10 {
		log.Fatalf("Ошибка поиска узла %v\n", result.key[0])
	}

	root = root.Remove(10)
	if root.Search(10) != nil {
		log.Fatalf("Узел не удалилился")
	}
}
