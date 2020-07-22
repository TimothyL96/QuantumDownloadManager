package setting

type settingConfigOption func(u *Setting) error

// NewSetting returns a new instance of Setting.
func NewSetting(configurations ...settingConfigOption) (*Setting, error) {
	setting := &Setting{}

	for _, configuration := range configurations {
		if err := configuration(setting); err != nil {
			return nil, err
		}
	}

	return setting, nil
}

// Functional options functions:

// NrOfConcurrentConnection helps set the value of concurrent connection during instance creation
func NrOfConcurrentConnection(nrOfConcurrentConnection int) settingConfigOption {
	return func(s *Setting) error {
		return s.SetNrOfConcurrentConnection(nrOfConcurrentConnection)
	}
}
