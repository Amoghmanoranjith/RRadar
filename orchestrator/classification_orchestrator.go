package orchestrator

// gets multiple strategies
// uses them as it sees fit

import (
	"errors"
	"math/rand/v2"
	"rradar/domain"
	customError "rradar/error"
	"rradar/model"
	"time"
)

type ClassificationOrchestrator struct {
	classifiers []domain.Classifier
}

func NewClassificationOrchestrator(c []domain.Classifier) *ClassificationOrchestrator {
	return &ClassificationOrchestrator{
		classifiers: c,
	}
}

func (m *ClassificationOrchestrator) Classify(posts []model.Post) ([]model.ClassifiedPost, error) {
	classifiedPosts := []model.ClassifiedPost{}

	for _, e := range posts {
		// init with success = false
		success := false
		// maybe check what kind of error this is later
		var rateErr *customError.RateLimitError
		// only 4 retries for now
		for i := 0; i < 6; i++ {
			result, err := m.classifiers[0].Classify(e)
			if errors.As(err, &rateErr) {
				// now back off beween requests
				val := time.Duration(rand.IntN(1<<(i+1)) + 1)
				// sleep for 1 to 2^(i+1)
				time.Sleep(val * time.Second)
				continue
			}
			if err != nil {
				break
			}
			classifiedPosts = append(classifiedPosts, result)
			success = true
			break
		}

		if !success {
			// try with fallback
			for i := 0; i < 6; i++ {
				result, err := m.classifiers[1].Classify(e)
				if errors.As(err, &rateErr) {
					// now back off beween requests
					val := time.Duration(rand.IntN(1<<(i+1)) + 1)
					// sleep for 1 to 2^(i+1)
					time.Sleep(val * time.Second)
					continue
				}
				if err != nil {
					break
				}
				classifiedPosts = append(classifiedPosts, result)
				success = true
				break
			}
		}

	}

	return classifiedPosts, nil
}

/*
take a classifier
put classifier in RetryHandler and let it call
make a call
if response is 429 then backoff

*/
