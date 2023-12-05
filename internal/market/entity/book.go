package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Order            []*Order
	Transactions     []*Transaction
	OrdersChannel    chan *Order
	OrdersChannelOut chan *Order
	Wg               *sync.WaitGroup
}

func NewBook(orderChannel chan *Order, orderChannelOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Order:            []*Order{},
		Transactions:     []*Transaction{},
		OrdersChannel:    orderChannel,
		OrdersChannelOut: orderChannelOut,
		Wg:               wg,
	}
}

func (b *Book) Trade() {
	buyOrders := NewOrderQueue()
	sellOrders := NewOrderQueue()

	heap.Init(buyOrders)
	heap.Init(sellOrders)

	for order := range b.OrdersChannel {
		if order.OrderType == "BUY" {
			buyOrders.Push(order)

			if sellOrders.Len() > 0 && sellOrders.Orders[0].Price <= order.Price {
				sellOrder := sellOrders.Pop().(*Order)

				if sellOrder.PendindShares > 0 {
					transaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price)
					b.AddTransaction(transaction, b.Wg)
					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)
					b.OrdersChannelOut <- sellOrder
					b.OrdersChannelOut <- order

					if sellOrder.PendindShares > 0 {
						sellOrders.Push(sellOrder)
					}
				}
			}
		} else if order.OrderType == "SELL" {
			sellOrders.Push(order)

			if buyOrders.Len() > 0 && buyOrders.Orders[0].Price >= order.Price {
				buyOrder := buyOrders.Pop().(*Order)

				if buyOrder.PendindShares > 0 {
					transaction := NewTransaction(order, buyOrder, order.Shares, buyOrder.Price)
					b.AddTransaction(transaction, b.Wg)
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)
					b.OrdersChannelOut <- buyOrder
					b.OrdersChannelOut <- order

					if buyOrder.PendindShares > 0 {
						buyOrders.Push(buyOrder)
					}
				}

			}
		}

	}
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	sellingShares := transaction.SellingOrder.PendindShares
	buyingShares := transaction.BuyingOrder.PendindShares

	minShares := sellingShares

	if buyingShares < sellingShares {
		minShares = buyingShares
	}

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.AddSellOrderPendingShares(-minShares)

	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, minShares)
	transaction.AddBuyOrderPendingShares(-minShares)

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)
	transaction.CloseBuyOrder()
	transaction.CloseSellOrder()

	b.Transactions = append(b.Transactions, transaction)
}
