package main

import "strings"

func filterBadWords(msg string) string {
	filteredValue := "****"

	badWords := map[string]string{
		"kerfuffle": filteredValue,
		"sharbert":  filteredValue,
		"fornax":    filteredValue,
	}

	words := strings.Split(msg, " ")
	for idx, word := range words {
		filteredWord, ok := badWords[strings.ToLower(word)]
		if ok {
			words[idx] = filteredWord
		}
	}
	return strings.Join(words, " ")
}
