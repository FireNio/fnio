package core

import (
	"sync/atomic"
	"unsafe"
)

type Node struct {
	item *string
	next *Node
}

type ConcurrentLinkedQueue struct {
	head *Node
	tail *Node
}

func newNode(item *string) *Node {
	var node = Node{}
	node.item = item
	return &node
}

func NewConcurrentLinkedQueue() *ConcurrentLinkedQueue {
	var queue = ConcurrentLinkedQueue{}
	var node = newNode(nil)
	queue.head = node
	queue.tail = node
	return &queue
}

func (queue *ConcurrentLinkedQueue) Offer(item *string) bool {
	var newNode = newNode(item)
	var t = queue.tail
	var p = t
	for ; ; {
		var q = p.next
		if q == nil {
			// source: if (p.casNext(null, newNode)) {
			// p is last node
			if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&p.next)), nil, unsafe.Pointer(newNode)) {
				// Successful CAS is the linearization point
				// for e to become an element of this queue,
				// and for newNode to become "live".
				if p != t {
					// source: casTail(t, newNode);
					// hop two nodes at a time
					atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&queue.tail)), unsafe.Pointer(t), unsafe.Pointer(newNode)) // Failure is OK.
				}
				return true
			}
			// Lost CAS race to another thread; re-read next
		} else if p == q {
			// We have fallen off list.  If tail is unchanged, it
			// will also be off-list, in which case we need to
			// jump to head, from which all live nodes are always
			// reachable.  Else the new tail is a better bet.
			//t = queue.tail
			var t_temp = t
			t = queue.tail
			if t_temp != t {
				p = t
			} else {
				p = queue.head
			}
		} else {
			// Check for tail updates after two hops.
			var t_temp = t
			t = queue.tail
			if p != t_temp && t_temp != t {
				p = t
			} else {
				p = q
			}
		}
	}
}

func (queue *ConcurrentLinkedQueue) Poll() *string {
restartFromHead:
	for ; ; {
		var h = queue.head
		var p = h
		var q *Node = nil
		for ; ; {
			var item = p.item
			// source: if (item != null && p.casItem(item, null)) {
			if item != nil {
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&p.item)), unsafe.Pointer(p.item), nil) {
					// Successful CAS is the linearization point
					// for item to be removed from this queue.
					if p != h {
						// source: updateHead(h, ((q = p.next) != null) ? q : p);
						// hop two nodes at a time
						q = p.next
						var temp = q
						if q == nil {
							temp = p
						}
						if h != p && atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&queue.head)), unsafe.Pointer(h), unsafe.Pointer(temp)) {
							h.next = h // h.lazySetNext(h);
						}
					}
					return item
				}
			}
			q = p.next
			if q == nil {
				// source: updateHead(h, p);
				if h != p && atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&queue.head)), unsafe.Pointer(h), unsafe.Pointer(p)) {
					h.next = h // source: h.lazySetNext(h);
				}
				return nil
			} else if p == q {
				continue restartFromHead
			} else {
				p = q
			}
		}
	}

}
