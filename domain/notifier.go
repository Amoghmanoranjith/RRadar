package domain

import "rradar/model"

type Notifier interface {
	Notify(entry model.ClassifiedPost) (error error)
}
