package handlers
// LinksRequest — тело запроса для POST /links.
// Пример:
// {
//   "links": ["google.com", "yandex.ru"]
// }
//
// @Description Запрос на создание набора ссылок для проверки.	
type LinksRequest struct {
		// Список URL для проверки.
	// Может содержать одну или несколько ссылок.
	Links []string `json:"links"`
}
// LinksResponse — тело ответа для POST /links.
//
// @Description Ответ с присвоенным номером набора и начальными статусами ссылок.
type LinksResponse struct {
		// Статусы по каждой ссылке: pending / available / not available.	
	Links    map[string]string `json:"links"`
	// Уникальный номер набора ссылок.
	LinksNum int64             `json:"links_num"`
}
// ReportRequest — тело запроса для POST /links/report.
// Пример:
// {
//   "links_list": [1, 2, 3]
// }
//
// @Description Запрос на формирование PDF-отчёта по списку ранее созданных наборов.
type ReportRequest struct {
	// Список идентификаторов наборов ссылок (links_num),
	// для которых нужно сформировать отчёт.
	LinksList []int64 `json:"links_list"`
}
