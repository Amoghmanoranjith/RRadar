// gets multiple strategies
// uses them as it sees fit

package classifier

import (
	"rradar/domain"
	modelLLM "rradar/model/llm"
	modelXML "rradar/model/xml"
)

type StrategyManager struct {
    classifiers []domain.Classifier
}

func NewStrategyManager(c []domain.Classifier) *StrategyManager {
    return &StrategyManager{
        classifiers: c,
    }
}

func (m *StrategyManager) Process(feed modelXML.Feed) (modelLLM.Feed, error) {

	relevantFeed := modelLLM.NewFeedWithEmptyEntries(feed)

	for _, e := range feed.Entries{
		result, err := m.classifiers[0].Classify(e)
		// maybe check what kind of error this is later
		if err != nil {
			result, _ := m.classifiers[1].Classify(e)
			relevantFeed.Entries = append(relevantFeed.Entries, result)
			continue
		}
		relevantFeed.Entries = append(relevantFeed.Entries, result)
	}

    return relevantFeed, nil
}