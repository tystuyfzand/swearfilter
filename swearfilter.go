package swearfilter

import (
	"golang.org/x/text/runes"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Filter is the swear filter implementation
type Filter interface {
	Check(msg string) ([]Match, error)
	Add(words ...string)
	Delete(words ...string)
	Words() []string
}

// Match is a matched filtered word
// The index is the index in the final message, not the source (YET)
type Match struct {
	Word  string
	Index int
}

var unicodeSet = unicodeSpaceSet{}

type Option func(s *swearFilter)

// DisableNormalize Disables normalization of alphabetic characters if set to true (ex: à -> a)
func DisableNormalize() Option {
	return func(s *swearFilter) {
		s.disableNormalize = true
	}
}

// DisableSpacedTab Disables converting tabs to singular spaces (ex: [tab][tab] -> [space][space])
func DisableSpacedTab() Option {
	return func(s *swearFilter) {
		s.disableSpacedTab = true
	}
}

// DisableMultiWhitespaceStripping Disables stripping down multiple whitespaces (ex: hello[space][space]world -> hello[space]world)
func DisableMultiWhitespaceStripping() Option {
	return func(s *swearFilter) {
		s.disableMultiWhitespaceStripping = true
	}
}

// DisableZeroWidthStripping Disables stripping zero-width spaces
func DisableZeroWidthStripping() Option {
	return func(s *swearFilter) {
		s.disableZeroWidthStripping = true
	}
}

// EnableSpacedBypass Disables testing for spaced bypasses (if hell is in filter, look for occurrences of h and detect only alphabetic characters that follow; ex: h[space]e[space]l[space]l[space] -> hell)
func EnableSpacedBypass() Option {
	return func(s *swearFilter) {
		s.enableSpacedBypass = true
	}
}

//swearFilter contains settings for the swear filter
type swearFilter struct {
	//Options to tell the swear filter how to operate
	disableNormalize                bool
	disableSpacedTab                bool
	disableMultiWhitespaceStripping bool
	disableZeroWidthStripping       bool
	enableSpacedBypass              bool

	//A list of words to check against the filters
	BadWords map[string]struct{}
	mutex    sync.RWMutex
}

// New returns an initialized swearFilter struct to check messages against
func New(filterWords []string, opts ...Option) Filter {
	filter := &swearFilter{
		BadWords: make(map[string]struct{}),
	}

	for _, opt := range opts {
		opt(filter)
	}

	for _, word := range filterWords {
		filter.BadWords[word] = struct{}{}
	}

	return filter
}

//Check will return any words that trip an enabled swear filter, an error if any, or nothing if you've removed all the words for some reason
func (filter *swearFilter) Check(msg string) ([]Match, error) {
	filter.mutex.RLock()
	defer filter.mutex.RUnlock()

	if filter.BadWords == nil || len(filter.BadWords) == 0 {
		return nil, nil
	}

	message := strings.ToLower(msg)

	//Normalize the text
	if !filter.disableNormalize {
		normalize := transform.Chain(norm.NFD, runes.Remove(unicodeSet), norm.NFC)

		bytes := make([]byte, len(message))

		_, _, err := normalize.Transform(bytes, []byte(message), true)
		if err != nil {
			return nil, err
		}
		message = string(bytes)
	}

	//Turn tabs into spaces
	if !filter.disableSpacedTab {
		message = strings.Replace(message, "\t", " ", -1)
	}

	//Get rid of zero-width spaces
	if !filter.disableZeroWidthStripping {
		message = strings.Replace(message, "\u200b", "", -1)
	}

	//Convert multiple re-occurring whitespaces into a single space
	if !filter.disableMultiWhitespaceStripping {
		regexLeadCloseWhitepace := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
		message = regexLeadCloseWhitepace.ReplaceAllString(message, "")
		regexInsideWhitespace := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
		message = regexInsideWhitespace.ReplaceAllString(message, "")
	}

	// TODO Matches won't properly capture the index of the original message
	// This needs to be redone to keep track of what's changed, and keep the start/end of the replaced items
	trippedWords := make([]Match, 0)
	checkSpace := false
	for swear := range filter.BadWords {
		if swear == " " {
			checkSpace = true
			continue
		}

		if idx := strings.Index(message, swear); idx != -1 {
			trippedWords = append(trippedWords, Match{
				Word:  swear,
				Index: idx,
			})

			continue
		}

		if filter.enableSpacedBypass {
			nospaceMessage := strings.Replace(message, " ", "", -1)

			if idx := strings.Index(nospaceMessage, swear); idx != -1 {
				trippedWords = append(trippedWords, Match{
					Word:  swear,
					Index: idx,
				})
			}
		}
	}

	if checkSpace && message == "" {
		trippedWords = append(trippedWords, Match{
			Word:  " ",
			Index: 0,
		})
	}

	return trippedWords, nil
}

//Add appends the given word to the uhohwords list
func (filter *swearFilter) Add(badWords ...string) {
	filter.mutex.Lock()
	defer filter.mutex.Unlock()

	if filter.BadWords == nil {
		filter.BadWords = make(map[string]struct{})
	}

	for _, word := range badWords {
		filter.BadWords[word] = struct{}{}
	}
}

//Delete deletes the given word from the uhohwords list
func (filter *swearFilter) Delete(badWords ...string) {
	filter.mutex.Lock()
	defer filter.mutex.Unlock()

	for _, word := range badWords {
		delete(filter.BadWords, word)
	}
}

//Words return the uhohwords list
func (filter *swearFilter) Words() (activeWords []string) {
	filter.mutex.RLock()
	defer filter.mutex.RUnlock()

	if filter.BadWords == nil {
		return nil
	}

	for word := range filter.BadWords {
		activeWords = append(activeWords, word)
	}
	return
}

type unicodeSpaceSet struct {
}

func (u unicodeSpaceSet) Contains(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}
