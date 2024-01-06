package internal

import (
	"regexp"
	"strconv"
)

type constError string

func (err constError) Error() string {
	return string(err)
}

// Errors raised by the internal package
const (
	ErrLinkLocalHostNotFound = constError("no Link Local host found")
)

type LabelParser interface {
	FindLinkLocalHostName(labels map[string]string) (string, error)
	IsTraefikEnabled(labels map[string]string) bool
}

type regexLabelParser struct {
	labelRe  *regexp.Regexp
	domainRe *regexp.Regexp
}

func (p *regexLabelParser) FindLinkLocalHostName(labels map[string]string) (string, error) {
	for key, value := range labels {
		if p.labelRe.Match([]byte(key)) {
			match := p.domainRe.FindStringSubmatch(value)
			if len(match) > 0 {
				hostname := match[0]
				return hostname, nil
			}
		}
	}
	return "", ErrLinkLocalHostNotFound
}

func (p *regexLabelParser) IsTraefikEnabled(labels map[string]string) bool {
	var err error
	value, ok := labels["traefik.enable"]
	if !ok {
		// We have no `traefik.enable` at all so it is not active
		return ok
	}
	// Let's check the value
	ok, err = strconv.ParseBool(value)
	if err != nil {
		// Unable to parse the value so expecting false
		return false
	}
	return ok
}

func NewLabelParser() LabelParser {
	parser := &regexLabelParser{}

	parser.labelRe = regexp.MustCompile(`traefik\.http\.routers\.(.*)\.rule`) // error if regexp invalid
	parser.domainRe = regexp.MustCompile(`(?P<domain>[^\x60]*?\.local)`)      // error if regexp invalid
	return parser
}
