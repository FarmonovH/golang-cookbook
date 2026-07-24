# Error

Golang dasturlash tilida error qanday ishlaydi. Error handling ko'p dasturlash tillarida xar xil ishlaydi masalan python, java va boshqa tillarda error handling try catch yoki except orqali handling qilinadi va maxsus class orqali generatsiya qilinadi. 
    Bu bizga katta code blokini bitta try catch ichida ishlashi aynan qaysi ishlamay qolganligini bilish noqulaylik tug'duradi.
Shuning alohida try catch yozilsa code o'qilishi(readable) ligi tushib ketadi va uzunlashib ketadi.

Golang dasturlash tili esa errorni value sifatida qaraydi va uni funksya orqali return qilinadi funksya chaqirilganda error nil ekanligi tekshiriladi 

```go 
package main

import "errors"

func checkNum(num int)(bool, error) {
    if num < 0 {
        err := errors.New("son noldan kichik")
        return false, err
    }
    return true, nil 
}

/*
ushbu ko'rinishda handling qilinadi 
*/

func main(){
    val, err := checkNum(-323)
    if err != nil {
        panic(err)
    }  
    fmt.Println(val)
}
```
## Custom error yaratish

Goda error yaratishning bir nechta usullari bor 

```go 
var err1 error = errors.New("your error message") /* bu usul error message uchun */
var err2 error = fmt.Errorf("your error message and %w", err1) 
/* fmt orqali yaratilsa bunda errors new dan imkoniyati ko'p bu haqida batafsil gaplashmiz */
```
endi qanday qilib custom error struct yaratishni ko'rib chiqamiz
*error* aslida interface 
```go 
type error interface {
    Error() string
} 
```
shuni qaysi struct implement qilsa error ga aylanadi 

masalan
```go 
package main

import (
	"fmt"
)

/* bizga kerak fields bilan to'ldiramiz masalan backend status code etc*/
type ConnError struct {
	Line int
	Message string
	Code int
}

/* Error methodini implement qiladi */
func (c ConnError)Error() string {
	return fmt.Sprintf("Message: %s, Line: %d, Code: %d", c.Message, c.Line, c.Code)
}

/* bu yerda custom error pointer return qilamiz error bizda interface shuning uchun & bilan return qilsa bo'ladi */
func connError() error {
	return &ConnError{
		Line: 32,
		Message: "some error",
		Code: 500,
	}
}

/* qanday ishlatilishi */
func main(){
	if err := connError(); err != nil {
        err, ok := err.(*ConnError)
		if ok {
			fmt.Println(err.Error())
		} else {
			panic(err)
		}
	}
}
```

## Wrapping errors (errorlarni wrap qilish)

Wrapping degani qandaydir errorni olib ustiga qo'shimcha yozib yana qaytarish 
biz buni fmt.Errorf yordamida amalga oshiramiz.

```go 
func connDatabase() error {
    return errors.New("err conn to database")
}

func executeQuery() error {
    if err := connDatabase(); err != nil {
        return err
    }
}

func httpHandler() error {
    if err := executeQuery(); err != nil {
        return err
    }
}

```
shu ko'rinishda yozsak codeni bu xatoni topish qiyinlashadi masalan error httpHandlerda bo'lsa ham ***err conn to database*** chiqadi, bu esa qaysi qatlamda error bo'lganini bilaolmaymiz 

endi keyingi codega nazar soling 
```go 
func connDatabase() error {
    return errors.New("err conn to database")
}

func executeQuery() error {
    if err := connDatabase(); err != nil {
        return fmt.Errorf("execute query: %w", err)
    }
}

func httpHandler() error {
    if err := executeQuery(); err != nil {
        return fmt.Errorf("http handler: %w", err)
    }
}
```

bu yerda errorni wrap qildik qaysi qatlamda error chiqsa aniq bila olamiz qaysi qatlamda error borligini  
execute query: err conn to database  deb chiqadi shu error wrapping deyiladi.

## Errorlarni tekshirish (Inspecting errors)

Errorlarni tekshirish ushbu error aynan bu errormi degan savolga javob beradi.

Bizda ayrim paytlar errorlarni solishtirishimizga to'g'ri keladi masalan

```go 
func main() {
    err1 := errors.New("some...")
    err2 := errors.New("some...")
    /* == bu operator bilan tenglashtirsa bo'ladi lekin bitta joyi xato bo'lsa kutilmagan vaziyatlarga olib keladi*/
    if err1.Error() == "some.." {
        //do someting
    }
}
```
lekin wrapper errorlarda eng asosiy error bilan solishtirganda == operatori ishlamaydi 
errors.Is(err, target) bu solishtiradi va wrapped bo'lsa unwrap qilib rekrusiv tekshiradi 
errors.At(error, &type_error) bu error turini tekshiradi

## Errorlar Panic bilan 

Panic bu qaysi function ishlayotgan bo'lsa ushbu gorutine ni to'xtatadi qaysidur funksya ichida ishayotgan bo'lsa ushbu funksyani to'xtatadi 
faqat panicdan keyin defer ishlaydi bo'ldi.

```go 
func A(){
    defer func(){
         // some code 
    }()

    panic("panic")
}
```

Bu codeda brinchi panic beradi va keyin defer ichidagi code ishlaydi 

```go 
package main

import (
	"fmt"
)

func  A(){
	defer fmt.Println("defered A")
	B()
}

func B(){
	defer func(){
		if err := recover(); err != nil {
			fmt.Println("panic ", err)
		} else {
			fmt.Println("panic not found")
		}
	}()
	C()
}

func C() {
	defer fmt.Println("defer C")
	raisePanic()
}

func raisePanic(){
	panic("raise panic function")
}

func main() {
	A()
}
```

bu yerda error ushlab qolinib qayta ishlanyabdi 



