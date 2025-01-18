package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

const (
	alphaChars      string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ012345678u"
	numberChars     string = "0123456789"
	specialChars    string = "!@#$%^&*(){}[]|\\;:'\",.<>/?`~"
	baseSpeedFactor int    = 20
)

func main() {

	input := parseInputArgs()
	config := newConfig(input)
	sPhrase := getScrambledPhrase(config.phrase, config.chars)

	//while the scrambled phrase does not equal the original phrase, loop through each rune in the phrase and randomly get a new character for that runes index if it did not match
	//this loop will continue until all characters in the scrambled string match the original string.
	phrase := config.phrase
	maxLen := len(string(phrase)) //will get the length of the string we are displaying
	fmt.Print("\033[H\033[2J") //clear the screen
	for !slices.Equal(phrase, sPhrase) {
		time.Sleep(time.Second / time.Duration(config.speedFactor))
		for i := range phrase {
			if phrase[i] != sPhrase[i] {
				sPhrase[i] = randomChar(config.chars)
			} else {
				sPhrase[i] = phrase[i]
			}
		}
		display := string(sPhrase)
		if len(display) > maxLen {
			maxLen = len(display)
		}
		fmt.Printf("\r%s", strings.Repeat(" ", maxLen)) //will write out max found length of display with blank chars
		fmt.Printf("\r%s", display)
	}
	fmt.Println()
}

// inputArgs struct that contains the input arguments
type inputArgs struct {
	phrase      string
	incAlpha    bool
	incNumbers  bool
	incSpecial  bool
	incAll      bool
	custom      string
	speedFactor int
}

// parseInputArgs parses the input arguments and returns a pointer to an inputArgs struct
func parseInputArgs() *inputArgs {
	alpha := flag.Bool("a", false, "Include alpha characters")
	numbers := flag.Bool("n", false, "Include numbers")
	special := flag.Bool("s", false, "Include special characters")
	all := flag.Bool("A", false, "Include all characters, will override the alpha, number, and special flags")
	custom := flag.String("c", "", "Custom characters to include, will override the all, alpha, number and special flags")
	speedFactor := flag.Int("sf", baseSpeedFactor, "Speed factor")

	flag.Parse()

	phrase := flag.Arg(0)

	// TODO: get the phrase from a piped input

	input := inputArgs{
		phrase:      phrase,
		incAlpha:    *alpha,
		incNumbers:  *numbers,
		incSpecial:  *special,
		incAll:      *all,
		custom:      *custom,
		speedFactor: *speedFactor,
	}

	return &input
}

// config struct that contains the phrase, characters, and speed factor
type config struct {
	phrase      []rune
	chars       []rune
	speedFactor int
}

// newConfig returns a new config struct based on the provided input arguments
func newConfig(i *inputArgs) *config {

	//get the phrase from the input args, if there wasnt one we can get one from a list of premades
	phrase := []rune(i.phrase)
	if len(phrase) == 0 {
		defaultPhrases := []string{
			"Hello, World!",
			"Hello, from scram",
			"This is a test",
			"Cheers! \U0001F37B",
		}
		phrase = []rune(defaultPhrases[rand.Intn(len(defaultPhrases))])
	}

	//get the base caracters from the input args.
	var baseChars string
	if i.incAll {
		baseChars = alphaChars + numberChars + specialChars
	} else if i.custom != "" {
		baseChars = i.custom
	} else {
		if i.incAlpha {
			baseChars += alphaChars
		}
		if i.incNumbers {
			baseChars += numberChars
		}
		if i.incSpecial {
			baseChars += specialChars
		}
	}

	return &config{
		phrase:      phrase,
		chars:       getChars(phrase, baseChars),
		speedFactor: i.speedFactor,
	}
}

// randomChar returns a random character from the provided slice of runes
func randomChar(fromChars []rune) rune {
	return fromChars[rand.Intn(len(fromChars))]
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
