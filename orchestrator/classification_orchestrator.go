package orchestrator
// gets multiple strategies
// uses them as it sees fit

import (
	"rradar/domain"
	"rradar/model"

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

	for _, e := range posts{
		result, err := m.classifiers[0].Classify(e)
		// maybe check what kind of error this is later
		if err != nil {
			result, _ := m.classifiers[1].Classify(e)
			classifiedPosts = append(classifiedPosts, result)
			continue
		}
		classifiedPosts = append(classifiedPosts, result)
	}

    return classifiedPosts, nil
}