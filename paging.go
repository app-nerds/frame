package frame

import "math"

/*
AdjustPage decrements the value of "page" because we want to use
zero-based pages for the math. Make sure page is never less than
zero.
*/
func (fa *FrameApplication) AdjustPage(page int) int {
	page--

	if page < 0 {
		page = 0
	}

	return page
}

/*
HasNextPage returns true when the result of page multiplied by
pageSize is less than the total recordCount.
*/
func (fa *FrameApplication) HasNextPage(page, pageSize, recordCount int) bool {
	return ((page * pageSize) + pageSize) < recordCount
}

/*
TotalPages returns how many pages are available in a paged result
based pageSize and the total recordCount.
*/
func (fa *FrameApplication) TotalPages(pageSize, recordCount int) int {
	return int(math.Ceil(float64(recordCount) / float64(pageSize)))
}
