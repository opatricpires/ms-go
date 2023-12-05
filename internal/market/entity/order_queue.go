package entity

type OrderQueue struct {
	Orders []*Order
}

// LESS
func (oq *OrderQueue) Less(i, j int) bool {
	return oq.Orders[i].Price < oq.Orders[j].Price
}

// SWAP
func (oq *OrderQueue) Swap(i, j int) {
	oq.Orders[i], oq.Orders[j] = oq.Orders[j], oq.Orders[i]
}

// LEN
func (oq *OrderQueue) Len() int {
	return len(oq.Orders)
}

// PUSH
func (oq *OrderQueue) Push(x interface{}) {
	oq.Orders = append(oq.Orders, x.(*Order))
}

// POP
func (oq *OrderQueue) Pop() interface{} {
	old := oq.Orders
	n := len(old)
	item := old[n-1]
	oq.Orders = old[0 : n-1]
	return item
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{}
}
