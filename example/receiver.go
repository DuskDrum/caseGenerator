package example

type Receiver struct {
}

// ReceiverSample 简单的receiver
func (a Receiver) ReceiverSample() {

}

// ReceiverStarSample 指针的receiver
func (a *Receiver) ReceiverStarSample() {

}

// ReceiverStarSample 指针的receiver
//func (a dict.ReceiverDict) ReceiverStarSample() {
//
//}

func (a *Receiver) ReceiverEllipsis() {

}
