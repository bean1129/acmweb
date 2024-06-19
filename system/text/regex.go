package text

import (
	"regexp"
	"sync"
)

type ZRegex struct {
	regexMu  sync.RWMutex
	regexMap map[string]*regexp.Regexp
}

func NewRegex() *ZRegex {
	return &ZRegex{
		regexMu: sync.RWMutex{},
		// Cache for regex object.
		// Note that:
		// 1. It uses sync.RWMutex ensuring the concurrent safety.
		// 2. There's no expiring logic for this map.
		regexMap: make(map[string]*regexp.Regexp),
	}
}

func (c *ZRegex) MatchString(pattern string, src string) ([]string, error) {
	if r, err := c.getRegexp(pattern); err == nil {
		return r.FindStringSubmatch(src), nil
	} else {
		return nil, err
	}
}

func (c *ZRegex) getRegexp(pattern string) (regex *regexp.Regexp, err error) {
	// Retrieve the regular expression object using reading lock.
	c.regexMu.RLock()
	regex = c.regexMap[pattern]
	c.regexMu.RUnlock()
	if regex != nil {
		return
	}
	// If it does not exist in the cache,
	// it compiles the pattern and creates one.
	regex, err = regexp.Compile(pattern)
	if err != nil {
		return
	}
	// Cache the result object using writing lock.
	c.regexMu.Lock()
	c.regexMap[pattern] = regex
	c.regexMu.Unlock()
	return
}
