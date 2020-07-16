package core

import "unsafe"

type IntPtrMap struct {
	load_factor float32
	cap         int32
	keys        []int64
	values      []unsafe.Pointer
	scan_size   int32
	size        int32
	mask        int64
	scan_index  int32
	limit       int32
}

func NewIntPtrMap(size int32, loadFactor float32) *IntPtrMap {
	var c = ClothCoverInt32(size)
	var mmap = IntPtrMap{
		cap:         c,
		mask:        int64(c - 1),
		load_factor: loadFactor,
		keys:        make([]int64, c),
		values:      make([]unsafe.Pointer, c),
		limit:       int32(float32(c) * loadFactor),
	}
	FillInt64(mmap.keys[:], -1)
	return &mmap
}

func (m *IntPtrMap) Scan() {
	m.scan_size = 0
	m.scan_index = -1
}

func (m *IntPtrMap) PutIfAbsent(key int64, value unsafe.Pointer) unsafe.Pointer {
	var v = m.Get(key)
	if v == nil {
		v = m.Put(key, value)
	}
	return v
}

func (m *IntPtrMap) Put(key int64, value unsafe.Pointer) unsafe.Pointer {
	var values = m.values
	var keys = m.keys
	var mask = m.mask
	var start_index = key & mask
	var index = start_index
	for ; ; {
		var _key = keys[index]
		if _key == -1 {
			keys[index] = key
			values[index] = value
			m.grow()
			return nil
		} else if _key == key {
			var old = values[index]
			values[index] = value
			return old
		}
		index = next_key(index, mask)
		if index == start_index {
			panic("failed to put")
		}
	}
}

func next_key(key int64, mask int64) int64 {
	return (key + 1) & mask
}

func put(keys []int64, values []unsafe.Pointer, key int64, value unsafe.Pointer, mask int64) {
	var start_index = key & mask
	var index = start_index
	for {
		var _key = keys[index]
		if _key == -1 {
			keys[index] = key
			values[index] = value
			break
		} else if _key == key {
			values[index] = value
			break
		}
		index = next_key(index, mask)
		if index == start_index {
			panic("failed to put")
		}
	}
}

func (m *IntPtrMap) HasNext() bool {
	return m.Next() != -1
}

func (m *IntPtrMap) Next() int32 {
	if m.scan_size < m.size {
		var keys = m.keys
		var index = m.scan_index + 1
		var cap = m.cap
		for ; index < cap; index++ {
			if keys[index] != -1 {
				break
			}
		}
		m.scan_size++
		m.scan_index = index
		return index
	}
	return -1
}

func (m *IntPtrMap) Key() int64 {
	return m.keys[m.scan_index]
}

func (m *IntPtrMap) Value() unsafe.Pointer {
	return m.values[m.scan_index]
}

func (m *IntPtrMap) IndexKey(index int32) int64 {
	return m.keys[index]
}

func (m *IntPtrMap) IndexValue(index int32) unsafe.Pointer {
	return m.values[index]
}

func (m *IntPtrMap) grow() {
	m.size++
	if m.size > m.limit {
		var cap = ClothCoverInt32(m.cap + 1)
		var mask = int64(cap - 1)
		var keys = make([]int64, cap)
		var values = make([]unsafe.Pointer, cap)
		var limit = (int32)(float32(cap) * m.load_factor)
		FillInt64(keys, -1)
		m.Scan()
		var size int32 = 0
		for ; ; {
			var index = m.Next()
			if index == -1 {
				break
			}
			size++
			put(keys, values, m.IndexKey(index), m.IndexValue(index), mask)
		}
		if size != m.size {
			panic("IntPtrMap error")
		}
		m.cap = cap
		m.mask = mask
		m.keys = keys
		m.values = values
		m.limit = limit
	}
}

func (m *IntPtrMap) Get(key int64) unsafe.Pointer {
	var keys = m.keys
	var mask = m.mask
	var start_index = key & mask
	var index = start_index
	for {
		var _key = keys[index]
		if _key == -1 {
			return nil
		} else if _key == key {
			return m.values[index]
		}
		index = next_key(index, mask)
		if index == start_index {
			return nil
		}
	}
}
// ref from java IdentityHashMap
func remove_at(keys []int64, values []unsafe.Pointer, index int64, mask int64) {
	// Adapted from Knuth Section 6.4 Algorithm R
	// Look for items to swap into newly vacated slot
	// starting at index immediately following deletion,
	// and continuing until a null slot is seen, indicating
	// the end of a run of possibly-colliding keys.
	var i = next_key(index, mask)
	var key = keys[i]
	var next = index
	for ; key != -1; i = next_key(i, mask) {
		// The following test triggers if the item at slot i (which
		// hashes to be at slot r) should take the spot vacated by d.
		// If so, we swap it in, and then continue with d now at the
		// newly vacated i.  This process will terminate when we hit
		// the null slot at the end of this run.
		// The test is messy because we are using a circular table.
		var r = key & mask
		if (i < r && (r <= next || next <= i)) || (r <= next && next <= i) {
			keys[next] = key
			values[next] = values[i]
			keys[i] = -1
			values[i] = nil
			next = i
		}
		key = keys[i]
	}
}

func (m *IntPtrMap) Remove(key int64) unsafe.Pointer {
	var keys = m.keys
	var mask = m.mask
	var start_index = key & mask
	var index = start_index
	for {
		var _key = keys[index]
		if _key == -1 {
			return nil
		} else if _key == key {
			m.size--
			var value = m.values[index]
			keys[index] = -1
			m.values[index] = nil
			remove_at(keys, m.values, index, mask)
			return value
		}
		index = next_key(index, mask)
		if index == start_index {
			return nil
		}
	}
}

func (m *IntPtrMap) IsEmpty() bool {
	return m.size == 0
}

func (m *IntPtrMap) Size() int32 {
	return m.size
}

func (m *IntPtrMap) Clear() {
	FillInt64(m.keys, -1)
	FillPtrNil(m.values)
	m.size = 0
}
