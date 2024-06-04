package hashsum

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

type ChecksumFormat int

const (
	ChecksumFormatBSD ChecksumFormat = iota
	ChecksumFormatGNU
)

var (
	bsdLineRex = regexp.MustCompile(`^SHA256 +\((?P<name>[^)]+)\) += +(?P<checksum>[a-z0-9]+)$`)
	gnuLineRex = regexp.MustCompile(`^(?P<checksum>[0-9a-z]+) +(?P<name>\S+)$`)
)

func Parse(stream io.Reader, format ChecksumFormat) (map[string]string, error) {
	// The regex should match the group as a whole and both subgroups
	const expectedMatchCount = 3

	var lineRex *regexp.Regexp
	switch format {
	case ChecksumFormatGNU:
		lineRex = gnuLineRex
	case ChecksumFormatBSD:
		lineRex = bsdLineRex
	default:
		panic("unknown checksum type")
	}

	checksums := map[string]string{}
	s := bufio.NewScanner(stream)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if matches := lineRex.FindStringSubmatch(line); len(matches) == expectedMatchCount {
			name := ""
			checksum := ""
			for i, groupName := range lineRex.SubexpNames() {
				if groupName == "name" {
					name = matches[i]
					if format == ChecksumFormatGNU && strings.HasPrefix(name, "*") {
						name = name[1:]
					}
				} else if groupName == "checksum" {
					checksum = matches[i]
				}
			}
			checksums[name] = checksum
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return checksums, nil
}
