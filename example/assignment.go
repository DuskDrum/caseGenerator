package example

import "fmt"

func (e Example) AssignmentAssertionTypeFunc(req any) {
	eReq, ok := req.(Example)
	if ok {
		amount := eReq.Amount
		time := eReq.orderTime
		id := eReq.OrderId
		fmt.Printf("amount: %v; time: %v; id: %v", amount, time, id)
	}
	eReqStar, ok := req.(*Example)
	if eReqStar == nil {
		fmt.Println("example star")
		return
	}
	if ok {
		amount := eReqStar.Amount
		time := eReqStar.orderTime
		id := eReqStar.OrderId
		fmt.Printf("amount: %v; time: %v; id: %v", amount, time, id)
	}

}
