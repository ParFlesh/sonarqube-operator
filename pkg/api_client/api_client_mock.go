package api_client

type APIClientMock struct {
	PingError error
}

func (r *APIClientMock) New(URL string) APIReader {
	return r
}

func (r *APIClientMock) Ping() error {
	return r.PingError
}
