package main

import "encoding/binary"

type Item struct {
		key   []byte
		value []byte
}

type Node struct {
		*dal

		pageNum    pgnum
		items      []*Item
		childNodes []pgnum
}

func NewEmptyNode() *Node {
		return &Node{}
}

func newItem(key, value []byte) *Item {
		return &Item{
				key:   key,
				value: value,
		}
}

func (n *Node) isLeaf() bool {
		return len(n.childNodes) == 0
}

func (n *Node) serialize(buf []byte) []byte {
		leftPos := 0
		rightPos := len(buf) - 1

		isLeaf := n.isLeaf()

		var bitSetVar uint64
		if isLeaf {
				bitSetVar = 1
		}

		buf[leftPos] = byte(bitSetVar)
		leftPos += 1

		binary.LittleEndian.PutUint16(buf[leftPos:], uint16(len(n.items)))

		leftPos += 2

		for i:=0; i < len(n.items); i++ {
				items := n.items[i]

				if !isLeaf {
						childNode := n.childNodes[i]

						binary.LittleEndian.PutUint64(buf[leftPos:], uint64(childNode))
						leftPos += pageNumSize
				}

				klen := len(items.key)
				vlen := len(items.value)
				
				offset := leftPos - klen - vlen - 2
				binary.LittleEndian.PutUint16(buf[offset:], uint16(klen))

				leftPos += 2

				rightPos -= vlen
				copy(buf[leftPos:], items.value)

				rightPos -= 1
				buf[rightPos] = byte(vlen)

				rightPos -= klen
				copy(buf[rightPos:], items.key)

				rightPos -= 1
				buf[rightPos] = byte(klen)
		}

		if !isLeaf {
				lastChildNode := n.childNodes[len(n.childNodes)-1]

				binary.LittleEndian.PutUint64(buf[leftPos:], uint64(lastChildNode))	
		}

		return buf
}

func (n *Node) deserialize(buf []byte) {		
		leftPos := 0

		isLeaf := uint16(buf[0])

		itemsCount := int(binary.LittleEndian.Uint16(buf[1:3]))
		leftPos += 3

		for i:=0; i < itemsCount; i++ {
				if isLeaf == 0 {
						pageNum := binary.LittleEndian.Uint64(buf[leftPos:])
						leftPos += pageNumSize

						n.childNodes = append(n.childNodes, pgnum(pageNum))
				}

				offset := binary.LittleEndian.Uint16(buf[leftPos:])
				leftPos += 2

				klen := uint16(buf[int(offset)])
				offset += 1

				key := buf[offset : offset+klen]
				offset += klen

				vlen := uint16(buf[int(offset)])
				offset += 1

				value := buf[offset : offset+vlen]
				offset += vlen

				n.items = append(n.items, newItem(key, value))
		}

		if isLeaf == 0 {
				pageNum := pgnum(binary.LittleEndian.Uint64(buf[leftPos:]))
				n.childNodes = append(n.childNodes, pageNum)
		}
}

