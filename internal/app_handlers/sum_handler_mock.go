package app_handlers

type SumHandlerMock struct {
	HandleCalled bool
	LastRequest  SumRequest
	NextResponse *SumResponse
	NextError    error
}

func (m *SumHandlerMock) Handle(req SumRequest) (*SumResponse, error) {
	m.HandleCalled = true
	m.LastRequest = req
	return m.NextResponse, m.NextError
}
