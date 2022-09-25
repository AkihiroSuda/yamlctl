package yamlutil

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
)

// splitCommentOnlyHeaderFooter is needed because goccy/go-yaml does not work well with comment-only headers/footers.
//
// TODO: fix goccy/go-yaml upstream so that it works with comment-only headers/footers
func splitCommentOnlyHeaderFooter(b []byte) (header, body, footer []byte, err error) {
	r := bufio.NewReader(bytes.NewReader(b))
	var (
		lines                     [][]byte
		firstNonCommentBlockBegin = -1
		currentCommentBlockBegin  = -1
	)
	for i := 0; ; i++ {
		var line []byte
		line, err = r.ReadBytes(byte('\n')) // line contains the delim
		isEOF := errors.Is(err, io.EOF)
		if isEOF {
			err = nil
			// we still have a valid line (without the delim)
		}
		if err != nil {
			return
		}

		trimmedLine := strings.TrimSpace(string(line))
		if commentOrEmpty := trimmedLine == "" || strings.HasPrefix(trimmedLine, "#"); commentOrEmpty {
			if currentCommentBlockBegin == -1 {
				currentCommentBlockBegin = i
			}
		} else {
			if firstNonCommentBlockBegin == -1 {
				firstNonCommentBlockBegin = i
			}
			currentCommentBlockBegin = -1
		}
		lines = append(lines, line)
		if isEOF {
			break
		}
	}
	for i := 0; i < firstNonCommentBlockBegin; i++ {
		header = append(header, lines[i]...)
	}
	if firstNonCommentBlockBegin >= 0 {
		for i := firstNonCommentBlockBegin; i < len(lines); i++ {
			if currentCommentBlockBegin >= 0 && i >= currentCommentBlockBegin {
				footer = append(footer, lines[i]...)
			} else {
				body = append(body, lines[i]...)
			}
		}
	}
	return
}
