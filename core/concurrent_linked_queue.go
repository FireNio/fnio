package core

// ref from juc.ConcurrentLinkedQueue
import (
	"sync/atomic"
	"unsafe"
)

type Node struct {
	item unsafe.Pointer
	next *Node
}

type ConcurrentLinkedQueue struct {
	head *Node
	tail *Node
}

func newNode(item unsafe.Pointer) *Node {
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

/**
  public boolean offer(E e) {
      checkNotNull(e);
      final Node<E> newNode = new Node<E>(e);

      for (Node<E> t = tail, p = t;;) {
          Node<E> q = p.next;
          if (q == null) {
              // p is last node
              if (p.casNext(null, newNode)) {
                  // Successful CAS is the linearization point
                  // for e to become an element of this queue,
                  // and for newNode to become "live".
                  if (p != t) // hop two nodes at a time
                      casTail(t, newNode);  // Failure is OK.
                  return true;
              }
              // Lost CAS race to another thread; re-read next
          }
          else if (p == q)
              // We have fallen off list.  If tail is unchanged, it
              // will also be off-list, in which case we need to
              // jump to head, from which all live nodes are always
              // reachable.  Else the new tail is a better bet.
              p = (t != (t = tail)) ? t : head;
          else
              // Check for tail updates after two hops.
              p = (p != t && t != (t = tail)) ? t : q;
      }
  }
*/

func (queue *ConcurrentLinkedQueue) Offer(item unsafe.Pointer) bool {
	var newNode = newNode(item)
	var t = queue.tail
	var p = t
	for {
		var q = p.next
		if q == nil {
			// p is last node
			if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&p.next)), nil, unsafe.Pointer(newNode)) {
				// Successful CAS is the linearization point
				// for e to become an element of this queue,
				// and for newNode to become "live".
				if p != t {
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

/**
  public E poll() {
      restartFromHead:
      for (;;) {
          for (Node<E> h = head, p = h, q;;) {
              E item = p.item;

              if (item != null && p.casItem(item, null)) {
                  // Successful CAS is the linearization point
                  // for item to be removed from this queue.
                  if (p != h) // hop two nodes at a time
                      updateHead(h, ((q = p.next) != null) ? q : p);
                  return item;
              }
              else if ((q = p.next) == null) {
                  updateHead(h, p);
                  return null;
              }
              else if (p == q)
                  continue restartFromHead;
              else
                  p = q;
          }
      }
  }
*/

func (queue *ConcurrentLinkedQueue) Poll() unsafe.Pointer {
restartFromHead:
	for {
		var h = queue.head
		var p = h
		var q *Node
		for {
			var item = p.item
			if item != nil && atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&p.item)), item, nil) {
				// Successful CAS is the linearization point
				// for item to be removed from this queue.
				if p != h {
					// hop two nodes at a time
					q = p.next
					var temp *Node
					if q != nil {
						temp = q
					} else {
						temp = p
					}
					updateHead(queue, h, temp)
				}
				return item
			} else {
				q = p.next
				if q == nil {
					updateHead(queue, h, p)
					return nil
				} else if p == q {
					continue restartFromHead
				} else {
					p = q
				}
			}
		}
	}
}

/**
  final void updateHead(Node<E> h, Node<E> p) {
      if (h != p && casHead(h, p))
          h.lazySetNext(h);
  }
*/
func updateHead(queue *ConcurrentLinkedQueue, h *Node, p *Node) {
	if h != p && atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&queue.head)), unsafe.Pointer(h), unsafe.Pointer(p)) {
		h.next = h // source: h.lazySetNext(h);
	}
}

/**
  public E peek() {
      restartFromHead:
      for (;;) {
          for (Node<E> h = head, p = h, q;;) {
              E item = p.item;
              if (item != null || (q = p.next) == null) {
                  updateHead(h, p);
                  return item;
              }
              else if (p == q)
                  continue restartFromHead;
              else
                  p = q;
          }
      }
  }
*/

func (queue *ConcurrentLinkedQueue) Peek() unsafe.Pointer {
restartFromHead:
	for {
		var h = queue.head
		var p = h
		var q *Node
		for {
			var item = p.item
			q = p.next
			if item != nil || q == nil {
				updateHead(queue, h, p)
				return item
			} else if p == q {
				continue restartFromHead
			} else {
				p = q
			}
		}
	}
}

func (queue *ConcurrentLinkedQueue) IsEmpty() bool {
	return queue.Peek() == nil
}
