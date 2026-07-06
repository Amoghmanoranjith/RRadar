package domain

import (
	"rradar/model"
)

type Classifier interface{
	Classify(post model.Post)( model.ClassifiedPost, error)
}