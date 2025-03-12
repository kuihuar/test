package disignpattern

import "fmt"

type PaymentStrategy interface {
	Pay(amount float64) error
}

type CreditCardPayment struct{}

func (c *CreditCardPayment) Pay(amount float64) error {
	fmt.Println("Paying with credit card:", amount)
	return nil
}

type PayPalPayment struct{}

func (p *PayPalPayment) Pay(amount float64) error {
	fmt.Println("Paying with PayPal:", amount)
	return nil
}

type PaymentProcessor struct {
	paymentStrategy PaymentStrategy
}

func NewPaymentProcessor(paymentStrategy PaymentStrategy) *PaymentProcessor {
	return &PaymentProcessor{paymentStrategy: paymentStrategy}
}

func (pp *PaymentProcessor) Process(amount float64) error {
	return pp.paymentStrategy.Pay(amount)
}

func StragegyExample() {
	creditCardPay := &CreditCardPayment{}
	paymentProcess := NewPaymentProcessor(creditCardPay)

	err := paymentProcess.Process(120.00)
	fmt.Println(err)

}
