package app

import (
	"slices"
)

type FoundComic struct {
	Id    int
	Url   string
	Count int
}

func (a *App) indexSearch(keywords []string) map[int]int {
	counts := make(map[int]int)
	for _, k := range keywords {
		for _, id := range a.db.GetIndex()[k] {
			counts[id]++
		}
	}
	return counts
}

func (a *App) dbSearch(keywords []string) map[int]int {
	counts := make(map[int]int)
	for _, k := range keywords {
		for id, c := range a.db.GetAll() {
			if _, contains := slices.BinarySearch(c.Keywords, k); contains {
				counts[id]++
			}
		}
	}
	return counts
}

func (a *App) GetTopN(keywords []string, n int, useIndex bool) []FoundComic {
	found := make([]FoundComic, 0, a.db.Size())
	comics := a.db.GetAll()
	var counts map[int]int

	if useIndex {
		counts = a.indexSearch(keywords)
	} else {
		counts = a.dbSearch(keywords)
	}

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
