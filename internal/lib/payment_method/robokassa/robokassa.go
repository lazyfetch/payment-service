package robokassa

// Request нужен для валидации в json и дальнейшей отправки robokassa api
// Для получения ссылки на оплату, также содержит все поля о магазине и пр. информации
type Request struct {
}

// Response нужен для также валидации в json ответа от webhook'a robokassa api
// Вид взаим POST Robokassa -> Response OK от ResultURL = они прекращают слать вебхук
// Нужно чтобы было понятно что с нашим сервисом все окей, и мы получили хук
type Response struct {
}

func GeneratePaymentURL() {

}

func CheckResultURL() {

}
