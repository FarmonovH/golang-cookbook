# Go Interface — Amaliy mashq

Interfeyslarni chuqur tushunish uchun bosqichma-bosqich murakkablashib boruvchi topshiriq. Har bir bosqichni o'zingiz yozib chiqing, keyin tekshirib oling.

## 📋 Topshiriq: "To'lov tizimi" (Payment System)

Turli xil to'lov usullarini (karta, naqd pul, kripto) bitta umumiy interfeys orqali boshqaradigan tizim yozing.

---

### 1-bosqich: Interfeys va asosiy turlar

```go
package main

import "fmt"

// TODO: Payer interfeysini yarating
// U ikkita metodga ega bo'lishi kerak:
//   Pay(amount float64) error
//   Name() string

// TODO: CardPayment struct yarating (maydonlar: Number string, Balance float64)
// TODO: CashPayment struct yarating (maydonlar: Owner string)
// TODO: CryptoPayment struct yarating (maydonlar: Wallet string, Balance float64)

// Har uch struct uchun Payer interfeysini implement qiling:
// - CardPayment.Pay(): agar Balance yetarli bo'lmasa, error qaytaring
// - CashPayment.Pay(): har doim muvaffaqiyatli (naqd pul cheklanmagan deb hisoblang)
// - CryptoPayment.Pay(): agar Balance yetarli bo'lmasa, error qaytaring

func main() {
	// TODO: []Payer slice yarating, uchala turdan ham qo'shing
	// TODO: har birini for-range bilan aylanib, Pay(100) chaqiring
	// TODO: xatolikni tekshiring va chop eting
}
```

**Maqsad:** oddiy interfeys implementatsiyasi va polimorfizmni his qilish — bitta `[]Payer` slice ichida turli struct'lar bir xil interfeys orqali ishlaydi.

---

### 2-bosqich: Type assertion qo'shish

`main()` ichida har bir to'lovni qayta ishlaganingizdan keyin:

- Agar to'lov turi **`*CardPayment`** bo'lsa — uning qolgan balansini chop eting.
- Agar **`*CryptoPayment`** bo'lsa — "⚠️ kripto tranzaksiya tasdiqlanishi kutilmoqda" degan xabar chop eting.
- Buning uchun **type assertion** (`, ok` formasi bilan) ishlating, `type switch` emas — buni keyingi bosqichda qilamiz.

---

### 3-bosqich: `type switch` bilan qayta yozish

2-bosqichda yozgan bir nechta `if _, ok := ...` bloklarini **bitta `switch v := payer.(type)`** ga aylantiring. Natijada kod qanchalik toza bo'lganini solishtiring.

---

### 4-bosqich: Ikkinchi interfeys qo'shish (interfeys kompozitsiyasi)

```go
type Refundable interface {
	Refund(amount float64) error
}
```

- `CardPayment` uchun `Refund` metodini implement qiling.
- `CashPayment` va `CryptoPayment` uchun **implement qilmang** (qasddan).
- `main()` ichida, har bir `Payer` uchun tekshiring: agar u **ham `Refundable` interfeysini** qanoatlantirsa (`type assertion` orqali `Payer` dan `Refundable`ga), `Refund(50)` chaqiring; aks holda "bu to'lov turi qaytarilmaydi" deb yozing.

> 💡 Bu qism eng muhim qism — chunki bu yerda siz **"interfeysdan interfeysga" assertion** qilishni o'rganasiz, bu `ConnError` misolidan bir qadam murakkabroq.

---

### 5-bosqich (qiyinroq, ixtiyoriy): Custom error turi + `errors.As`

```go
type InsufficientFundsError struct {
	Payer   string
	Missing float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("%s uchun %.2f yetishmayapti", e.Payer, e.Missing)
}
```

- `CardPayment.Pay()` va `CryptoPayment.Pay()` balans yetarli bo'lmaganda oddiy `errors.New(...)` o'rniga aynan shu `*InsufficientFundsError` ni qaytarsin.
- `main()` da har bir `Pay()` xatoligini `errors.As()` orqali tekshiring va agar bu aynan `InsufficientFundsError` bo'lsa, qancha pul yetishmayotganini alohida chop eting.

---

## ✅ Nimalarni mustahkamlaysiz

| Bosqich | O'rganiladigan tushuncha |
|---|---|
| 1 | Interfeys e'lon qilish, implement qilish, polimorfizm |
| 2 | Type assertion (`, ok` formasi) |
| 3 | `type switch` |
| 4 | Interfeysdan interfeysga assertion, interfeys kompozitsiyasi |
| 5 | Custom error turlari, `errors.As` |

---

## Qanday ishlatish kerak

Har bir bosqichni alohida yozib chiqing va build/run qiling:

```bash
go run main.go
```

Har bir bosqichni tugatgach, kodni ko'rib chiqib xato yoki noto'g'ri tushunilgan joylarni topish uchun ChatGPT/Claude yoki mentoringizga ko'rsating — to'liq yechimni so'ramasdan, faqat yo'naltirishni so'rang, shunda o'zingiz topib chiqasiz.