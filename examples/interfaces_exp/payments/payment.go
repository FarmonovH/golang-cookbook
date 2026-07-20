package payments


type PaymentMethods interface {
	Pay(string, float64)int
	Cancel(int)
}


type PaymentModule struct {
	paymentMethods PaymentMethods
	info map[int]Info
}

func NewPaymentModule(paymentMethods PaymentMethods) *PaymentModule {
	return &PaymentModule{
		paymentMethods: paymentMethods,
		info: make(map[int]Info),
	}
}


func (p PaymentModule) Pay(description string, price float64) int {
	id := p.paymentMethods.Pay(description, price)
	p.info[id] = Info{
		Description: description, 
		Price: price, 
		Canceled: false,
	}
	return id
}


func (p PaymentModule) Cancel(id int) {
	info, ok := p.info[id]
	if !ok {
		return 
	}
	info.Canceled = true
	p.info[id] = info
}

func (p PaymentModule) GetInfo(id int)(Info, bool) {
	info, ok := p.info[id]
	if !ok {
		return Info{}, false
	}
	return info, true
}

func (p PaymentModule) GetAllInfo() map[int]Info {
	temp := make(map[int]Info, len(p.info))
	for key, val := range p.info {
		temp[key] = val
	}
	return temp
}
