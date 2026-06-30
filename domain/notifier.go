package domain

import (
	modelLLM "rradar/model/llm"
)

type Notifier interface{
	Notify(entry modelLLM.Entry)
}