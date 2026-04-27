package domain

type ListId struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

var ValidProductTypes = []string{"STANDARD", "WEIGHTABLE", "KIT", "SERVICE"}

func IsValidProductType(productType string) bool {
	for _, t := range ValidProductTypes {
		if t == productType {
			return true
		}
	}
	return false
}
