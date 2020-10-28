package accounts

type Account struct {
	Id          string
	Name        string
	Environment string
	Role        string
	Region      string `default:"eu-west-1"`
}

func (a Account) asMap() map[string]string {
	return map[string]string{
		"Id":          a.Id,
		"Name":        a.Name,
		"Environment": a.Environment,
		"Role":        a.Role,
		"Region":      a.Region,
	}
}

// Arn returns the aws ARN formatted string
func (a Account) Arn() string {
	return "arn:aws:iam::" + a.Id + ":role/" + a.Role
}

// Get returns property
func (a Account) Get(prop string) string {
	return a.asMap()[prop]
}
