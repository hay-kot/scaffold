package models

// {{ .Each.Item }} model
type {{ .Each.Item | toPascalCase }} struct {
	ID int
}
