package main

import (
	"bufio"
	"flag"
	"fmt"
	"iter"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

const (
	alphaChars      string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ012345678"
	numberChars     string = "0123456789"
	specialChars    string = "!@#$%^&*(){}[]|\\;:'\",.<>/?`~"
	baseSpeedFactor int    = 20
)

func main() {
	input, err := parseInputArgs()
	if err != nil {
		fmt.Println(err)
	}

	config := newConfig(input)

	maxLen := len(string(config.phrase))
	fmt.Print("\033[H\033[2J")
	for p := range config.scram() {
		display := string(p)
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
func parseInputArgs() (*inputArgs, error) {
	alpha := flag.Bool("a", false, "Include alpha characters")
	numbers := flag.Bool("n", false, "Include numbers")
	special := flag.Bool("s", false, "Include special characters")
	all := flag.Bool("A", false, "Include all characters, will override the alpha, number, and special flags")
	custom := flag.String("c", "", "Custom characters to include, will override the all, alpha, number and special flags")
	speedFactor := flag.Int("sf", baseSpeedFactor, "Speed factor")

	flag.Parse()

	var phrase string

	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			phrase = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	} else {
		phrase = flag.Arg(0)
	}

	input := inputArgs{
		phrase:      phrase,
		incAlpha:    *alpha,
		incNumbers:  *numbers,
		incSpecial:  *special,
		incAll:      *all,
		custom:      *custom,
		speedFactor: *speedFactor,
	}

	return &input, nil
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
		chars:       []rune(baseChars),
		speedFactor: i.speedFactor,
	}
}

// scram yields a sequence of scrambled phrases that can be iterated over until the scramble matches the phrase
func (c *config) scram() iter.Seq[[]rune] {
	chars := mergeChars(c.phrase, c.chars)
	scram := scramblePhrase(c.phrase, chars)
	return func(yield func(p []rune) bool) {
		for !slices.Equal(c.phrase, scram) {
			time.Sleep(time.Second / time.Duration(c.speedFactor))

			for i := range c.phrase {
				if c.phrase[i] != scram[i] {
					scram[i] = randomChar(chars)
				} else {
					scram[i] = c.phrase[i]
				}

			}

			if !yield(scram) {
				return
			}
		}
	}
}

// randomChar returns a random character from the provided slice of runes
func randomChar(fromChars []rune) rune {
	return fromChars[rand.Intn(len(fromChars))]
}

// mergeChars returns a slice of runes that contains the base characters and any characters that are not in the base characters
func mergeChars(phrase []rune, baseChars []rune) []rune {
	chars := baseChars
	for i := 0; i < len(phrase); i++ {
		char := phrase[i]
		if !slices.Contains(chars, char) {
			chars = append(chars, char)
		}
	}
	return chars
}

// scramblePhrase returns a scrambled phrase based on the provided phrase and characters
func scramblePhrase(phrase, chars []rune) []rune {
	var s []rune
	for i := 0; i < len(phrase); i++ {
		s = append(s, randomChar(chars))
	}
	return s
}
