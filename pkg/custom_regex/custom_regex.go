package custom_regex

import "regexp"

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var NameRegex = regexp.MustCompile(`^[\p{L}'’\-]{1,40}$`)
var PasswordRegex = regexp.MustCompile(`^.{3,}$`)
