package application

type AuditLine string

const (
	Config_Updated     AuditLine = "Config_Updated"
	Parameters_Updated           = "Parameters_Updated"
)

func (s AuditLine) String() string {
	return string(s)
}
