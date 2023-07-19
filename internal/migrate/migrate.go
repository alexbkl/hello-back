/*
Package migrate provides database schema migrations.
*/

package migrate

import (
	"sync"

	"github.com/Hello-Storage/hello-back/internal/event"
)

var log = event.Log
var once sync.Once

// Values is a shortcut for map[string]interface{}
type Values map[string]interface{}
