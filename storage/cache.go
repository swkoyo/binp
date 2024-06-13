package storage

type CacheStore struct {
	client *LRUCache
}

func NewCache() *CacheStore {
	return &CacheStore{
		client: &LRUCache{
			length:        0,
			capacity:      100,
			head:          nil,
			tail:          nil,
			lookup:        make(map[string]*Node),
			reverseLookup: make(map[*Node]string),
		},
	}
}

type Node struct {
	val  *Snippet
	prev *Node
	next *Node
}

func createNode(val *Snippet) *Node {
	return &Node{
		val: val,
	}
}

type LRUCache struct {
	length        int
	capacity      int
	head          *Node
	tail          *Node
	lookup        map[string]*Node
	reverseLookup map[*Node]string
}

func (c *LRUCache) Get(key string) *Snippet {
	node, ok := c.lookup[key]
	if !ok {
		return nil
	}
	c.detatch(node)
	c.prepend(node)
	return node.val
}

func (c *LRUCache) Put(key string, value *Snippet) {
	node, ok := c.lookup[key]
	if ok {
		node.val = value
		c.detatch(node)
		c.prepend(node)
		return
	}
	node = createNode(value)
	c.prepend(node)
	c.length++
	c.lookup[key] = node
	c.reverseLookup[node] = key
	c.trimCache()
}

func (c *LRUCache) Delete(key string) {
	node, ok := c.lookup[key]
	if !ok {
		return
	}
	c.detatch(node)
	delete(c.lookup, key)
	delete(c.reverseLookup, node)
	c.length--
}

func (c *LRUCache) detatch(node *Node) {
	if node.next != nil {
		node.next.prev = node.prev
	}

	if node.prev != nil {
		node.prev.next = node.next
	}

	if node == c.head {
		c.head = c.head.next
	}

	if node == c.tail {
		c.tail = c.tail.prev
	}

	node.next = nil
	node.prev = nil
}

func (c *LRUCache) prepend(node *Node) {
	if c.head == nil {
		c.head = node
		c.tail = node
		return
	}

	node.next = c.head
	c.head.prev = node
	c.head = node
}

func (c *LRUCache) trimCache() {
	if c.length > c.capacity {
		node := c.tail
		key := c.reverseLookup[node]
		c.detatch(node)
		delete(c.lookup, key)
		delete(c.reverseLookup, node)
		c.length--
	}
}
