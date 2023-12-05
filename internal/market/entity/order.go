package entity

type Order struct {
	ID            string
	Investor      *Investor
	Asset         *Asset
	Shares        int
	PendindShares int
	Price         float64
	OrderType     string
	Status        string
	Transactions  []*Transaction
}

func NewOrder(orderID string, investor *Investor, asset *Asset,
	shares int,
	pendindShares int,
	price float64,
	orderType string) *Order {

	return &Order{
		ID:            orderID,
		Investor:      investor,
		Asset:         asset,
		Shares:        shares,
		PendindShares: shares,
		Price:         price,
		OrderType:     orderType,
		Status:        "OPEN",
		Transactions:  []*Transaction{},
	}

}
