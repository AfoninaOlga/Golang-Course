package app

import (
	"slices"
)

type FoundComic struct {
	Id    int
	Url   string
	Count int
}

func (a *App) IndexSearch(keywords []string) []FoundComic {
	res := make([]FoundComic, 0, a.db.Size())
	counts := make(map[int]int)
	comics := a.db.GetAll()
	for _, k := range keywords {
		for _, id := range a.db.GetIndex()[k] {
			counts[id]++
		}
	}
	for id, cnt := range counts {
		res = append(res, FoundComic{Id: id, Count: cnt, Url: comics[id].Url})
	}
	return res
}

func (a *App) DBSearch(keywords []string) []FoundComic {
	res := make([]FoundComic, 0, a.db.Size())
	counts := make(map[int]int)
	comics := a.db.GetAll()
	for _, k := range keywords {
		for id, c := range comics {
			if _, contains := slices.BinarySearch(c.Keywords, k); contains {
				counts[id]++
			}
		}
	}
	for id, cnt := range counts {
		res = append(res, FoundComic{Id: id, Count: cnt, Url: comics[id].Url})
	}
	return res
}

func (a *App) GetTopN(keywords []string, n int, useIndex bool) []FoundComic {
	var found []FoundComic
	if useIndex {
		found = a.IndexSearch(keywords)
	} else {
		found = a.DBSearch(keywords)
	}

	slices.SortFunc(found, func(a, b FoundComic) int {
		return b.Count - a.Count
	})
	if len(found) < n {
		return found
	}
	return found[:n]
}
