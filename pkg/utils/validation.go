package utils

import "regexp"

var urlPattern = `^((ftp|http|https):\/\/)?(www\.)?([A-Za-zА-Яа-я0-9]{1}[A-Za-zА-Яа-я0-9\-]*\.?)*\.{1}[A-Za-zА-Яа-я0-9-]{2,8}(\/([\w#!:.?+=&%@!\-\/])*)?/?`
var slugPatter = `^[A-Za-z0-9\-]*$`

var UrlRegex = regexp.MustCompile(urlPattern)
var SlugRegex = regexp.MustCompile(slugPatter)

func IsValidURL(url string) bool {
	return UrlRegex.MatchString(url)
}

func IsValidSlug(slug string) bool {
	return SlugRegex.MatchString(slug)
}
