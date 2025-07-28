package utils

func ParseNumber(from string) string {
	if len(from) >= 3 && from[:3] == "549" {
		from = from[:2] + from[3:]
	}
	return from
}
