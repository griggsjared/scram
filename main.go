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

	//get the phrase from the command line arguments
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

	//once we have the phrase we want to work with it as a slice of runes
	rPhrase := []rune(phrase)

	//add the characters from the phrase to the baseChars so that the scrambled phrase can contain the original characters
	rChars := []rune(baseChars)
	for i := 0; i < len(rPhrase); i++ {
		char := rPhrase[i]
		if !slices.Contains(rChars, char) {
			rChars = append(rChars, char)
		}
	}

	//scramble the phrase
	var scrambled []rune
	for i := 0; i < len(rPhrase); i++ {
		scrambled = append(scrambled, randomChar(rChars))
	}

	//while the scrambled phrase does not equal the original phrase, loop through each rune in the phrase and randomly get a new character for that runes index if it did not match
	//this loop will continue until all characters in the scrambled string match the original string.
	for !slices.Equal(rPhrase, scrambled) {
		time.Sleep(time.Second / time.Duration(baseSpeedFactor))
    for i := range rPhrase {
      if rPhrase[i] != scrambled[i] {
        scrambled[i] = randomChar(rChars)
      } else {
        scrambled[i] = rPhrase[i]
      }
    }
    fmt.Printf("\r%s", string(scrambled))
	}
	fmt.Println()
}

// randomChar returns a random character from the provided slice of runes
func randomChar(fromChars []rune) rune {
	return fromChars[rand.Intn(len(fromChars))]
}
