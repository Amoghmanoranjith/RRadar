package llm

import (
	modelXML "rradar/model/xml"
)

type Entry struct {
	modelXML.Entry
	Interesting bool
	Confidence  float32
	Reason      string
}
