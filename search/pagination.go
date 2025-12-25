package search

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
)

type Paginator struct {
	currentPage int
	maxPages    int
	hasMore     bool
}

func NewPaginator(maxPages int) *Paginator {
	return &Paginator{
		currentPage: 1,
		maxPages:    maxPages,
		hasMore:     true,
	}
}

func (p *Paginator) CurrentPage() int {
	return p.currentPage
}

func (p *Paginator) HasMore() bool {
	return p.hasMore && p.currentPage < p.maxPages
}

func (p *Paginator) Next() {
	p.currentPage++
}

func (p *Paginator) MarkComplete() {
	p.hasMore = false
}

func (p *Paginator) DetectEndOfResults(page *rod.Page, noResultsSelectors []string) bool {
	for _, selector := range noResultsSelectors {
		has, _, err := page.Has(selector)
		if err == nil && has {
			return true
		}
	}

	return false
}

func (p *Paginator) WaitForResults(page *rod.Page, resultSelector string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		has, _, err := page.Has(resultSelector)
		if err == nil && has {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for search results")
}

func (p *Paginator) ScrollToLoadMore(page *rod.Page, scrolls int) error {
	for i := 0; i < scrolls; i++ {
		_, err := page.Eval(`window.scrollBy(0, window.innerHeight * 0.8)`)
		if err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}
