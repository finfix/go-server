package service

// PublicURL возвращает канонический URL MCP
func (s *OAuthService) PublicURL() string {
	return s.config.PublicURL
}
