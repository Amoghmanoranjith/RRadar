package notifier

import (
	"rradar/domain"
    modelLLM "rradar/model/llm"
)

type StrategyManager struct{
	notifier domain.Notifier
}

func NewStrategyManager(notifier domain.Notifier) *StrategyManager {
    return &StrategyManager{
        notifier: notifier,
    }
}

func (m *StrategyManager) Process(feed modelLLM.Feed) ( error error) {


	for _, e := range feed.Entries{
        if e.Interesting == false{
            continue
        }
		err := m.notifier.Notify(e)
		// maybe check what kind of error this is later
		if err != nil {
			return err
		}
		
	}

    return nil
}