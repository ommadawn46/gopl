package bank

var deposits = make(chan int)
var balances = make(chan int)
var withdrawals = make(chan Withdrawal)

type Withdrawal struct {
	amount  int
	success chan bool
}

func Deposit(amount int) { deposits <- amount }
func Withdraw(amount int) bool {
	ch := make(chan bool)
	withdrawals <- Withdrawal{amount, ch}
	return <-ch
}
func Balance() int { return <-balances }

func teller() {
	var balance int
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case withdrawal := <-withdrawals:
			if balance < withdrawal.amount {
				withdrawal.success <- false
			} else {
				balance -= withdrawal.amount
				withdrawal.success <- true
			}
		case balances <- balance:
		}
	}
}

func init() {
	go teller()
}
