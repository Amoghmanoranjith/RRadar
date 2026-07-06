package orchestrator

import (
	"rradar/domain"
	"rradar/model"
)

type NotificationOrchestrator struct{
	notifier domain.Notifier
}

func NewNotificationOrchestrator(notifier domain.Notifier) *NotificationOrchestrator {
    return &NotificationOrchestrator{
        notifier: notifier,
    }
}

func (m *NotificationOrchestrator) Notify(classifiedPosts []model.ClassifiedPost ) (error error) {
	for _, e := range classifiedPosts{
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