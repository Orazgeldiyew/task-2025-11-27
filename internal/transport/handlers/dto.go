package handlers

// LinksRequest используется в POST /links.
type LinksRequest struct {
	Links []string `json:"links"`
}

// LinksResponse возвращается после создания набора ссылок.
type LinksResponse struct {
	Links    map[string]string `json:"links"`
	LinksNum int64             `json:"links_num"`
}

// ReportRequest используется в POST /links/report.
type ReportRequest struct {
	LinksList []int64 `json:"links_list"`
}
