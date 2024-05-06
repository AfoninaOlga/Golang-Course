package service

import (
	"slices"
)

type FoundComic struct {
	Id    int
	Url   string
	Count int
}

func (xs *XkcdService) indexSearch(keywords []string) map[int]int {
	counts := make(map[int]int)
	for _, k := range keywords {
		for _, id := range xs.db.GetIndex()[k] {
			counts[id]++
		}
	}
	return counts
}

func (xs *XkcdService) dbSearch(keywords []string) map[int]int {
	counts := make(map[int]int)
	for _, k := range keywords {
		for id, c := range xs.db.GetAll() {
			if _, contains := slices.BinarySearch(c.Keywords, k); contains {
				counts[id]++
			}
		}
	}
	return counts
}

func (xs *XkcdService) GetTopN(keywords []string, n int) []FoundComic {
	found := make([]FoundComic, 0, xs.db.Size())
	comics := xs.db.GetAll()
	var counts map[int]int

	counts = xs.indexSearch(keywords)

	for id, cnt := range counts {
		found = append(found, FoundComic{Id: id, Count: cnt, Url: comics[id].Url})
	}

	slices.SortFunc(found, func(a, b FoundComic) int {
		return b.Count - a.Count
	})
	if len(found) < n {
		return found
	}
	return found[:n]
}
