package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
)

type Word struct {
	FrenchWord  string `csv:"frenchWord"`
	EnglishWord string `csv:"englishWord"`
}

func main() {

	file, err := os.Open("words.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	wordFile, err := os.OpenFile("words.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	wordList := []*Word{}
	wasHappend := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ":") {
			wasHappend = false

			lineParts := strings.Split(scanner.Text(), ":")

			englishWord := strings.ToLower(strings.TrimSpace(lineParts[0]))
			frenchWord := strings.ToLower(strings.TrimSpace(lineParts[1]))

			if strings.Contains(englishWord, ",") {
				englishWord := multipleWords(englishWord)

				if strings.Contains(frenchWord, ",") {
					frenchWord = fmt.Sprintf("A REMPLACER(%v)", strings.Replace(frenchWord, ",", "/", -1))
				}

				for _, word := range englishWord {
					wordList = append(wordList, &Word{FrenchWord: frenchWord, EnglishWord: word})
				}

				wasHappend = true
			} else if strings.Contains(frenchWord, ",") {
				frenchWord := multipleWords(frenchWord)

				if strings.Contains(englishWord, ",") {
					englishWord = fmt.Sprintf("A REMPLACER(%v)", strings.Replace(englishWord, ",", "/", -1))
				}

				for _, word := range frenchWord {
					wordList = append(wordList, &Word{FrenchWord: word, EnglishWord: englishWord})
				}

				wasHappend = true
			}

			if !wasHappend {
				wordList = append(wordList, &Word{FrenchWord: frenchWord, EnglishWord: englishWord})
			}
		}
	}

	err = gocsv.MarshalFile(&wordList, wordFile)
	if err != nil {
		log.Fatal(err)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func multipleWords(wordsLine string) []string {
	var words []string

	for _, word := range strings.Split(wordsLine, ",") {
		words = append(words, strings.TrimSpace(word))
	}

	return words
}
