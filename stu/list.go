package stu

//ListNoder interface can be LISTED
type ListNoder interface{
    Dup() ListNoder
    Match() bool
}

//ListNode node of list
type ListNode struct {
    pre *ListNode
    next *ListNode
    value ListNoder
}

//DoubleEndList list with head and tail pointer
type DoubleEndList struct {
    length int  
    head *ListNode
    tail *ListNode
}

//NewDoubleEndList create DoubleEndList with no node
func NewDoubleEndList() *DoubleEndList {
    return & DoubleEndList {
        length: 0,
        head: nil,
        tail: nil,
    }
}

//Len get length of DoubleEndList
func (doubleEndList *DoubleEndList) Len() int {
    return doubleEndList.length
}

//First get first listNode of DoubleEndList
func (doubleEndList *DoubleEndList) First() *ListNode {
    return doubleEndList.head
}

//Last get last listNode of DoubleEndList
func (doubleEndList *DoubleEndList) Last() *ListNode {
    return doubleEndList.tail
}

//AddNodeHead insert node to head(unshift)
func (doubleEndList *DoubleEndList) AddNodeHead(h *ListNode) {
    h.next = doubleEndList.head
    if doubleEndList.tail == nil {
        doubleEndList.tail = h 
    }
    if doubleEndList.head != nil {
        h.pre = doubleEndList.head
    } 
    h.pre = nil
    doubleEndList.head = h
    doubleEndList.length++
}

//AddNodeTail append node to tail(push)
func (doubleEndList *DoubleEndList) AddNodeTail(t *ListNode) {
    t.pre = doubleEndList.tail
    if doubleEndList.tail != nil {
        doubleEndList.tail.next = t
    }
    if doubleEndList.head == nil {
        doubleEndList.head = t
    }
    t.next = nil
    doubleEndList.tail = t
    doubleEndList.length++
}
    
//InsertNode insert node n to pn pre if isPre, next otherwise
func (doubleEndList *DoubleEndList) InsertNode(n, pn *ListNode, isPre bool) {
    
}