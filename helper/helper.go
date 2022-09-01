package helper

import "net/url"

func IsValidUrl(str string) bool {
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	u, err := url.Parse(str)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func Pop[V int | int64 | float64 | string | byte | rune](slice []V) (value V, newSlice []V) {
	return slice[len(slice)-1], slice[:len(slice)-1]
}
