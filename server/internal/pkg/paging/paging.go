// Package paging 提供分页参数的归一化处理。
package paging

// NormalizePage 将分页参数规约到合理范围：page 最小 1，pageSize 范围 [1, 100]。
func NormalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
