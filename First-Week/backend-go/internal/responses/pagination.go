package responses

type Pagination struct {
    CurrentPage   int `json:"current_page"`
    ItemsPerPage  int `json:"items_per_page"`
    TotalItems    int `json:"total_items"`
    TotalPages    int `json:"total_pages"`
}