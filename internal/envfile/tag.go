package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// TagEntry represents a key annotated with a set of tags.
type TagEntry struct {
	Key  string
	Tags []string
}

// TagOptions controls the behaviour of Tag.
type TagOptions struct {
	// Rules maps a tag name to a list of key patterns (exact or glob-style
	// prefix with '*' suffix) that should receive that tag.
	Rules map[string][]string
	// SkipUntagged excludes keys that match no rule from the result.
	SkipUntagged bool
}

// DefaultTagOptions returns sensible defaults.
func DefaultTagOptions() TagOptions {
	return TagOptions{
		Rules:        map[string][]string{},
		SkipUntagged: false,
	}
}

// Tag annotates each key in env with tags according to opts.Rules.
// The returned slice preserves the iteration order of env (sorted by key).
func Tag(env map[string]string, opts TagOptions) ([]TagEntry, error) {
	if env == nil {
		return nil, fmt.Errorf("tag: env map must not be nil")
	}

	// pre-compile patterns
	type compiled struct {
		tag string
		re  *regexp.Regexp
	}
	var patterns []compiled
	for tag, globs := range opts.Rules {
		for _, g := range globs {
			re, err := globToRegexp(g)
			if err != nil {
				return nil, fmt.Errorf("tag: invalid pattern %q for tag %q: %w", g, tag, err)
			}
			patterns = append(patterns, compiled{tag: tag, re: re})
		}
	}

	keys := sortedKeys(env)
	var result []TagEntry
	for _, k := range keys {
		var tags []string
		for _, p := range patterns {
			if p.re.MatchString(k) {
				tags = appendUnique(tags, p.tag)
			}
		}
		if opts.SkipUntagged && len(tags) == 0 {
			continue
		}
		result = append(result, TagEntry{Key: k, Tags: tags})
	}
	return result, nil
}

// globToRegexp converts a simple glob (only '*' wildcard supported) to a
// compiled regexp anchored to the full string.
func globToRegexp(pattern string) (*regexp.Regexp, error) {
	escaped := regexp.QuoteMeta(pattern)
	// restore '*' wildcard
	regexStr := "^" + strings.ReplaceAll(escaped, `\*`, `.*`) + "$"
	return regexp.Compile(regexStr)
}

func appendUnique(slice []string, s string) []string {
	for _, v := range slice {
		if v == s {
			return slice
		}
	}
	return append(slice, s)
}

// sortedKeys returns the keys of m in alphabetical order.
func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
