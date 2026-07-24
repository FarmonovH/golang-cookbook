# String

Ko‘pgina dasturlar matn bilan u yoki bu tarzda ishlaydi — xoh foydalanuvchi bilan bevosita muloqot qilishda, xoh turli tizimlar yoki qurilmalar o‘rtasida ma’lumot almashishda. Matn deyarli barcha tizimlarda qo‘llaniladigan universal ma’lumot turi hisoblanadi va string ko‘rinishidagi ma’lumotlar deyarli hamma joyda uchraydi. Shu sababli, satrlarni samarali boshqarish va qayta ishlash dasturchining eng muhim ko‘nikmalaridan biridir.

Go dasturlash tilida satrlar bilan ishlash uchun bir nechta paketlar mavjud:

- strconv paketi satrlarni boshqa ma’lumot turlariga o‘zgartirish yoki aksincha, boshqa turlarni satrga aylantirish uchun ishlatiladi.
fmt paketi C tilidagi kabi formatlash belgilaridan (verbs) foydalanib satrlarni formatlash imkonini beradi.
unicode/utf8 va unicode/utf16 paketlari Unicode kodlashidan foydalanadigan satrlar bilan ishlash uchun kerakli funksiyalarni taqdim etadi.
strings paketi esa satrlar ustida bajariladigan eng ko‘p uchraydigan amallarni o‘z ichiga oladi. Agar kerakli funksiyani qaysi paketdan izlashni bilmasangiz, avvalo strings paketini tekshirish eng to‘g‘ri tanlov bo‘ladi


String yaratish:
- '' bittalik qo'shtirnoq bilan 
- "" 2 talik qo'shtirnoq bilan 
- `` bilan bir nechta qatorlar bilan 


