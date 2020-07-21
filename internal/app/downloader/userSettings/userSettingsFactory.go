package userSettings

type userSettingsConfigurationFunction func(u *UserSettings) error

func NewUserSettings(configurations ...userSettingsConfigurationFunction) (*UserSettings, error) {
	userSettings := &UserSettings{}

	for _, configuration := range configurations {
		err := configuration(userSettings)

		if err != nil {
			return nil, err
		}
	}

	return userSettings, nil
}

func NrOfConcurrentDownload(nrOfConcurrentDownload int) userSettingsConfigurationFunction {
	return func(u *UserSettings) error {
		return u.SetNrOfConcurrentDownload(nrOfConcurrentDownload)
	}
}
