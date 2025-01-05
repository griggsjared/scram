package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode/utf8"
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

	//add the characters from the phrase to the baseChars so that the scrambled phrase can contain the original characters
	combinedChars := baseChars
	for i := 0; i < utf8.RuneCountInString(phrase); i++ {
		char := string([]rune(phrase)[i])
		if !strings.Contains(combinedChars, string(char)) {
			combinedChars += string(char)
		}
	}

	//scramble the phrase
	var scrambled string
	for i := 0; i < utf8.RuneCountInString(phrase); i++ {
		scrambled += string(randomChar(combinedChars))
	}

	//while the scrambled phrase does not equal the original phrase, loop through each rune in the phrase and randomly get a new character for that runes index if it did not match
	//this loop will continue until all characters in the scrambled string match the original string.
	for scrambled != phrase {
		time.Sleep(time.Second / time.Duration(baseSpeedFactor))
		var n []rune
		cs1 := []rune(phrase)
		cs2 := []rune(scrambled)
		for i := 0; i < utf8.RuneCountInString(string(scrambled)); i++ {
			if cs1[i] != cs2[i] {
				n = append(n, randomChar(combinedChars))
			} else {
				n = append(n, cs1[i])
			}
		}
		scrambled = string(n)
		fmt.Printf("\r%s", scrambled)
	}
	fmt.Println()
}

// randomChar returns a random character from the provided string
func randomChar(fromChars string) rune {
	cs := []rune(fromChars)
	return cs[rand.Intn(len(cs))]
}
