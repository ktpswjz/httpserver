package main

type Target struct {
}

func (s *Target) GetUrl(schema, domain string) string {
	return "http://172.16.99.181"
}
