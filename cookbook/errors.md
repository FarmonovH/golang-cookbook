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

Wrapping degani 

