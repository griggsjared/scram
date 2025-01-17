package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"golang.org/x/exp/slices"
)

const (
	//baseChars is the base characters that can be used to scramble the phrase
	baseChars string = "!@#$%^&*(){}[]abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	//baseSpeedFactor is the base speed factor for the scramble speed.
	baseSpeedFactor int = 20
)

func main() {

	phrase := getPhrase()
	chars := getChars(phrase, baseChars)
  sPhrase := getScrambledPhrase(phrase, chars)

	//while the scrambled phrase does not equal the original phrase, loop through each rune in the phrase and randomly get a new character for that runes index if it did not match
	//this loop will continue until all characters in the scrambled string match the original string.
	for !slices.Equal(phrase, sPhrase) {
		time.Sleep(time.Second / time.Duration(baseSpeedFactor))
		for i := range phrase {
			if phrase[i] != sPhrase[i] {
				sPhrase[i] = randomChar(chars)
			} else {
				sPhrase[i] = phrase[i]
			}
		}
		fmt.Printf("\r%s", string(sPhrase))
	}
	fmt.Println()
}

// randomChar returns a random character from the provided slice of runes
func randomChar(fromChars []rune) rune {
	return fromChars[rand.Intn(len(fromChars))]
}

// getPhrase gets a provided phrase from stdin or uses a default phrase if no phrase is provided
func getPhrase() []rune {
	var phrase string
	ioArgs := os.Args[1:]
	if len(ioArgs) > 0 {
		phrase = ioArgs[0]
	}

	//if no phrase is provided, use a default phrase
	if len(phrase) == 0 {
		defaultPhrases := []string{
			"Hello, World!",
			"Hello, from scram",
			"This is a test",
			"Cheers! \U0001F37B",
		}
		phrase = defaultPhrases[rand.Intn(len(defaultPhrases))]
	}

	return []rune(phrase)
}

// getChars returns a slice of runes that contains the base characters and any characters that are not in the base characters
func getChars(phrase []rune, baseChars string) []rune {
	chars := []rune(baseChars)
	for i := 0; i < len(phrase); i++ {
		char := phrase[i]
		if !slices.Contains(chars, char) {
			chars = append(chars, char)
		}
	}
	return chars
}

// getScrambledPhrase returns a scrambled phrase based on the provided phrase and characters
func getScrambledPhrase(phrase, chars []rune) []rune {
	var s []rune
	for i := 0; i < len(phrase); i++ {
		s = append(s, randomChar(chars))
	}
	return s
}
