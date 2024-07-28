package data

type ApiServerSettings struct {
	UrlBase            string `yaml:"url_base"`
	ApiToken           string `yaml:"api_token"`
	AuthorizationValue string `yaml:"authorization_value"`
}

type Settings struct {
	Email             string            `yaml:"email"`
	AccountId         string            `yaml:"account_id"`
	AtlassianSettings ApiServerSettings `yaml:"atlassian"`
	TempoSettings     ApiServerSettings `yaml:"tempo"`
}
