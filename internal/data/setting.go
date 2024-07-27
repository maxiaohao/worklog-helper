package data

type ApiServerSettings struct {
	UrlBase  string `yaml:"url_base"`
	ApiToken string `yaml:"api_token"`
}

type Settings struct {
	Email             string            `yaml:"email"`
	AtlassianSettings ApiServerSettings `yaml:"atlassian"`
	TempoSettings     ApiServerSettings `yaml:"tempo"`
	AccountId         int               `yaml:"account_id"`
}
