package main

import (
	"regexp"
	"strings"
)

// Generate all possible permutations of the slice of strings
// We need this because it may be that the last name could come before the
// first or we could have an unusual arangement
func rearrange(pieces []string) [][]string {
	if len(pieces) == 1 {
		return [][]string{pieces}
	}

	var results [][]string
	for i := 0; i < len(pieces); i++ {
		choice := pieces[i]
		var remaining []string
		remaining = append(remaining, pieces[:i]...)
		remaining = append(remaining, pieces[i+1:]...)
		permuations := rearrange(remaining)

		for _, perm := range permuations {
			completePermuation := append(perm, choice)
			results = append(results, completePermuation)
		}
	}
	return results

}

// This function takes a name and generates possible emails
func generateEmails(name string, domain string) []string {
	name = strings.ToLower(name)
	pieces := strings.Split(name, " ")

	var usernames []string
	arrangements := rearrange(pieces)

	for _, arrangement := range arrangements {
		usernames = append(usernames, generateNamesWithArrangement(arrangement)...)
	}

	var emails []string
	for _, username := range usernames {
		emails = append(emails, username+"@"+domain)
	}

	return emails
}

// This function takes a name as an array of words (e.g. ["laura" "nyro"])
// and attempts to find names that could be email usernames
// (e.g. "laura.nyro", "lnyro", "laura", etc)
func generateNamesWithArrangement(pieces []string) []string {

	// Pick one string from each of the slices and join them together and
	// you'll have a potential email.
	var likelyForms [][]string

	for _, piece := range pieces {
		likelyForms = append(likelyForms, getLikelyForms(piece))
		likelyForms = append(likelyForms, []string{"", ".", "+", "_"})
	}

	total := 1
	for _, form := range likelyForms {
		total *= len(form)
	}

	var superset []string
	for i := 0; i < total; i++ {
		remaining := i
		var possible string
		for _, form := range likelyForms {
			possible += form[remaining%len(form)]
			remaining /= len(form)
		}
		superset = append(superset, possible)
	}

	// Filter out very unlikely usernames
	var filters []*regexp.Regexp
	filterStrings := []string{"(\\.|\\+|_)(\\.|\\+|_)", "^(\\.|\\+|_)", "(\\.|\\+|_)$", "^$"}
	for _, filterString := range filterStrings {
		filters = append(filters, regexp.MustCompile(filterString))
	}

	var results []string
	for _, possible := range superset {
		matches := false
		for _, filter := range filters {
			if len(filter.FindStringIndex(possible)) != 0 {
				matches = true
			}
		}

		if !matches {
			results = append(results, possible)
		}
	}

	return results
}

// andrew => ["andrew", "", "a", "an", "and"]
func getLikelyForms(piece string) []string {
	var result []string

	result = append(result, piece)

	// TODO: allow customizing max chars of a name used
	for i := 0; i < len(piece) && i <= 1; i++ {
		result = append(result, piece[0:i])
	}

	return result
}
