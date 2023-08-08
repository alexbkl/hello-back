package constant

import (
	"fmt"
)

var LoginMessage = "Greetings from joinhello\nSign this message to log into joinhello\nnonce: "

func BuildLoginMessage(nonce string) []byte {
	return []byte(fmt.Sprintf("%s%s", LoginMessage, nonce))
}
