package response

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Error   string `json:"error"`
	Meta    *Meta  `json:"meta"`
}

type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Meta    *Meta  `json:"meta"`
	Error   string `json:"error"`
}

func NewResponse(data any) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

func NewPaginatedResponse(data any, page, perPage int, total int64) PaginatedResponse {
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return PaginatedResponse{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}

func NewErrorResponse(err string) Response {
	return Response{
		Success: false,
		Error:   err,
	}
}
