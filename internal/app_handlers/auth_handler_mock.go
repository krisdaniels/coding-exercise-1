package app_handlers

type AuthHandlerMock struct {
	HandleCalled bool
	LastRequest  AuthRequest
	NextResponse *AuthResponse
	NextError    error
}

func (m *AuthHandlerMock) Handle(request AuthRequest) (*AuthResponse, error) {
	m.HandleCalled = true
	m.LastRequest = request
	return m.NextResponse, m.NextError
}
