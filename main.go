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

	app := newApp(input)

	// looping through the yield payload from scram,
	// clear the screen every iteration so it appears to
	// descramble the phrase in place.
	for p := range app.scram() {
		display := string(p)
		clearScreen()
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

	phrase, err := getStdinPhrase()
	if err != nil {
		return nil, fmt.Errorf("error reading from stdin: %w", err)
	}

	if phrase == "" {
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

func getStdinPhrase() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		lines := make([]string, 0)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
		if len(lines) > 0 {
			return strings.Join(lines, "\n"), nil
		}
	}
	return "", nil
}

// app struct that contains the phrase, characters, and speed factor
type app struct {
	phrase      []rune
	chars       []rune
	speedFactor int
}

// newApp returns a new app struct based on the provided input arguments
func newApp(i *inputArgs) *app {
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

	return &app{
		phrase:      phrase,
		chars:       []rune(baseChars),
		speedFactor: i.speedFactor,
	}
}

// scram yields a sequence of scrambled phrases that can be iterated over until the scramble matches the phrase
func (c *app) scram() iter.Seq[[]rune] {
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
	for i := range phrase {
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
	for range phrase {
		s = append(s, randomChar(chars))
	}
	return s
}

// clearScreen clears the terminal screen by printing the ANSI escape codes
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
