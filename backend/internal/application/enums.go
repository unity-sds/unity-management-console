package application

type AuditLine string

const (
	Config_Updated     AuditLine = "Config_Updated"
	Parameters_Updated           = "Parameters_Updated"

	Bootstrap_Successful = "Bootstrap_Successful"

	Bootstrap_Unsuccessful = "Bootstrap_Unsuccessful"
)

func (s AuditLine) String() string {
	return string(s)
}
