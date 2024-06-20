package main

import (
	"fmt"
	"sort"
)

type Candidate struct {
	Name  string
	Votes int
}

func countVotes(candidates []string) []Candidate {
	votesMap := make(map[string]int)

	for _, candidate := range candidates {
		votesMap[candidate]++
	}

	var result []Candidate
	for candidate, votes := range votesMap {
		result = append(result, Candidate{Name: candidate, Votes: votes})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Votes > result[j].Votes
	})

	return result
}

func main() {
	votes := []string{"Ann", "Kate", "Peter", "Kate", "Ann", "Ann", "Helen"}
	fmt.Println(countVotes(votes))
}
