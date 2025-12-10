package diagnose

type localDiagnoseConnection struct{}

func NewLocalDiagnoseConnection() (*localDiagnoseConnection, error) {
	return &localDiagnoseConnection{}, nil
}

func (c *localDiagnoseConnection) PublishDiagnose(msg []byte) error {
	return nil
}

func (c *localDiagnoseConnection) Close() {
}
