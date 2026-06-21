// this should give a reusable connection
package http

import (
	"net/http"
	"time"
)

var Client = &http.Client{
	Timeout: 10 * time.Second,
}