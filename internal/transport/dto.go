package transport

type LinksRequest struct {
	Links []string `json:"links"`
}

type LinksResponse struct {
	Links    map[string]string `json:"links"`
	LinksNum int64             `json:"links_num"`
}

type ReportRequest struct {
	LinksList []int64 `json:"links_list"`
}
