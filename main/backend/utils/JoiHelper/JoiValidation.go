package JoiHelper

import "github.com/softbrewery/gojoi/pkg/joi"

var UsernameSchema = joi.String().Min(3).Max(30).Regex("^[a-zA-Z][a-zA-Z._\\s]+[a-zA-Z]$")
var PasswordSchema = joi.String().Min(8).Max(50).Regex("[a-zA-Z0-9._!?#\\-\\s\\*]+")

// https://regex101.com/r/7wRJOy/2
var UserInput = joi.String().Min(3).Max(70).Regex("^[a-zA-Z0-9!#$€%&()\\-_=:+.\\s]+$")
