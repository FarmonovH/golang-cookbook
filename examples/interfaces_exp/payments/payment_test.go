package payments

import (
	"testing"
)


type MockPayment struct {
	PayCalled bool
	CancelCalled bool 
	
	LastDescription string
	LastPrice float64
	LastCancelId int

	ReturnedId int
}

func (m *MockPayment) Pay(description string, price float64) int {
	m.PayCalled = true 
	m.LastDescription = description
	m.LastPrice = price
	return m.ReturnedId
}

func (m *MockPayment) Cancel(id int) {
	m.CancelCalled = true 
	m.LastCancelId = id
}

func TestNewPaymentModule(t *testing.T) {
	mock := &MockPayment{}
	pm := NewPaymentModule(mock)

	if pm == nil {
		t.Fatal("excepted payment module")
	}

	if len(pm.info) !=  0 {
		t.Fatal("info map should be empty")
	} 
}

func TestPay(t *testing.T) {
	mock := &MockPayment{
		ReturnedId: 13,
	}

	pm := NewPaymentModule(mock)
	
	id := pm.Pay("Put Jedan", 160)
	if id != 13 {
		t.Fatalf("excepted id 13 god %d", id)
	}

	if !mock.PayCalled {
		t.Fatal("Pay method was not called")
	}

	if mock.LastDescription != "Put Jedan" {
		t.Fatal("Wrong description")
	}

	if mock.LastPrice != 160 {
		t.Fatal("Wrong price")
	}
}


func TestPayInfo(t *testing.T) {
	mock := &MockPayment{}

	pm := NewPaymentModule(mock)
	id := pm.Pay("put Jedan", 230)
	
	if _, ok := pm.GetInfo(id); !ok {
		t.Fatal("info not found")
	} 

	info, _ := pm.GetInfo(id)
	if info.Canceled {
		t.Fatalf("wrong canceled")
	}

	if info.Price != mock.LastPrice {
		t.Fatal("wrong price")
	}

	if info.Description != mock.LastDescription {
		t.Fatal("wrong description")
	}
}

func TestGetAllInfo(t *testing.T) {
	mock := &MockPayment{}

	pm := NewPaymentModule(mock)

	pm.Pay("put Jedan", 12)
	
	if len(pm.info) != 1 {
		t.Fatal("excepted one element")
	}

	delete(pm.info, 0)

	if len(pm.info) != 0 {
		t.Fatal("original map changed")
	}
}


func TestCancelUnknownId(t *testing.T) {
	mock := &MockPayment{}
	pm := NewPaymentModule(mock)

	pm.Cancel(990)
}



