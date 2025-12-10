package diagnose

type DiagnoseService interface {
	ServiceFanoutTest() (string, error)
}
