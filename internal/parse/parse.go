package parse

import "hw1taskmanager/internal/domain"

func ParseCommandLine(s string) ([]string, error) {
	res := make([]string, 0)

	quote := false
	currentPart := ""
	for _, c := range s {
		if quote {
			if c == '"' {
				res = append(res, currentPart)
				currentPart = ""
				quote = false
			} else {
				currentPart += string(c)
			}
		} else {
			if c == '"' {
				quote = true
			} else if c == ' ' {
				if currentPart != "" {
					res = append(res, currentPart)
					currentPart = ""
				}
			} else {
				currentPart += string(c)
			}
		}
	}
	if quote {
		return []string{}, domain.ErrQuoteDoesntClosed
	}

	if currentPart != "" {
		res = append(res, currentPart)
	}

	return res, nil
}
