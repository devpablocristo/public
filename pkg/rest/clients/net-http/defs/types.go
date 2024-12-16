package defs

type GenericTokenResponse struct {
	TokenData map[string]any
}

func (g *GenericTokenResponse) GetAccessToken() string {
	if token, ok := g.TokenData["access_token"].(string); ok {
		return token
	}
	return ""
}
