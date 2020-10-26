package accounts

type Account struct {
	Id          string
	Name        string
	Environment string
	Role        string
	Region      string `default:"eu-west-1"`
}

func (a Account) Arn() string {
	return "arn:aws:iam::" + a.Id + ":role/" + a.Role
}
