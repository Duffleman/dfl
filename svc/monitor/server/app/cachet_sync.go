package app

// SyncWithCachet syncs jobs with cachet components
func (a *App) SyncWithCachet() error {
	components, err := a.Cachet.ListAllComponents()
	if err != nil {
		return err
	}

OUT:
	for _, job := range a.Jobs {
		for _, component := range components {
			if component.Name == job.ComponentName {
				job.ComponentID = component.ID
				continue OUT
			}
		}
	}

	return nil
}
