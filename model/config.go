package model

type Properties struct {
	Mail     Mail
	Database Database
	Redis    Redis
}

type Mail struct {
	Dialer struct {
		Host     string
		Port     int
		Username string
		Password string
	}

	Receiver []string
}

type Database struct {
	Url string
}

type Redis struct {
	BigKey struct {
		Memory    int
		Len       int
		Priority  string
		Separator []string
	}
}
