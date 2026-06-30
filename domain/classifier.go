package domain

import (
	modelLLM "rradar/model/llm"
	modelXML "rradar/model/xml"
)

type Classifier interface{
	Classify(entry modelXML.Entry)(entryRelevance modelLLM.Entry, error error)
}