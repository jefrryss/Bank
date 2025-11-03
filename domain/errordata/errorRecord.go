package errordata

//Структура для записи строк и их ошибок при пасринге
type ErrorRecord struct {
	Line string
	Err  error
}
