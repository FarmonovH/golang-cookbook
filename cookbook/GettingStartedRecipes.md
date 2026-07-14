# Golangda boshlash uchun retseptlar 

Golangda tashqi 3-taraf package ishlatish uchun 

masalan 

```go
package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
)

func main(){
	var number uint64 = 1234567890
	fmt.Println("Size of file is ", humanize.Bytes(number)) 
}
```

biz ushbu codeni yurg'izmoqchi bo'lsak *go run main.go* qilsak 
human.go:6:2: no required module provides package github.com/dustin/go-humanize:
go.mod file not found in the current directory or any parent directory; see
'go help modules'
shunday error ga duch kelamiz ya'ni tashqi 3 taraf codini ishlatish uchun bizga o'zimizning packagega qo'shishimiz kerak.
Buning uchun bizga dastlab *go mod init package_name* qilib keying *go get github.com/dustin/go-humanize* qilsak bo'ladi. 3-taraf code qo'shildi endi ishlatsak bo'ladi agar ortiqcha ishlatilmayotgan packagelar bo'lsa ularni tozalash uchun *go mod tidy* qilsak yetarli.

## Test yozish 

Biz yozgan codimizni ishlashini isbotlash uchun test yozishimiz kerak.
*test yozish bizga berilgan masalani to'g'ri tushunishimiz to'g'ri fikrlashni o'rgatadi.*

misol 
```go 
package main

import "strconv"


func main(){}

func conv(str string)(num int64, err error){
	num, err = strconv.ParseInt(str, 2, 64)
	return 
}
```
bizda stringni songa o'giradigan funksya bor shuni test qilmoqchimiz 

test yozish uchun *filename_test.go* ko'rinishida bo'ladi ya'ni **_test.go** aniq bo'lishi kerak 

```go 
package main

import "testing"

func TestConv(t *testing.T){
	num, err := conv("123345334")
	if err != nil {
		t.Fatal(err)
	}

	if num != 123345334 {
		t.Fatal("number don't match")
	}

}

func TestFailConv(t *testing.T) {
	_, err := conv("")
	if err == nil {
		t.Fatal(err)
	}
}
```

shu testing kutubxonasidan foydalanib kichik test case yozamiz 
testni run qilishimiz uchun *go test -v* -v verbose name

shu ko'rinishda test qilinadi hali bu masalaga yana to'xtalamiz

## Uchinchi taraf modulini olish 

*go get github.com/example/package@v1.2.3*
yoki oxirgisi kerak bo'lsa 
***go get github.com/example/package@latest***


