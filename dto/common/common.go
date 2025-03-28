package common

// 此包包含所有DTO共用的结构体和常量
// 用于在不同的DTO之间共享通用的数据结构

// PaginationRequest 分页请求基础结构体
type PaginationRequest struct {
	Page     int `json:"page" form:"page" example:"1"`            // 页码，默认为1
	PageSize int `json:"page_size" form:"page_size" example:"10"` // 每页数量，默认为10
}

// GetDefaultPage 获取默认页码
func (p *PaginationRequest) GetDefaultPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetDefaultPageSize 获取默认每页数量
func (p *PaginationRequest) GetDefaultPageSize() int {
	if p.PageSize <= 0 {
		return 10
	}
	return p.PageSize
}

// PaginationSortRequest 带排序的分页请求结构体
type PaginationSortRequest struct {
	PaginationRequest
	SortBy    string `json:"sort_by" form:"sort_by" example:"created_at"` // 排序字段
	SortOrder string `json:"sort_order" form:"sort_order" example:"desc"` // 排序方向：asc-升序，desc-降序
}

// GetDefaultSortOrder 获取默认排序方向
func (p *PaginationSortRequest) GetDefaultSortOrder() string {
	if p.SortOrder == "" {
		return "desc"
	}
	return p.SortOrder
}

// PaginationResponse 分页响应结构体
type PaginationResponse struct {
	Total    int64 `json:"total"`     // 总记录数
	Page     int   `json:"page"`      // 当前页码
	PageSize int   `json:"page_size"` // 每页数量
	Pages    int   `json:"pages"`     // 总页数
}
