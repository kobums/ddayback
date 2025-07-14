package dday

const (
	CategoryPersonal = "개인"
	CategoryStudy    = "학업"
	CategoryWork     = "업무"
	CategoryOther    = "기타"
)

var Categories = []string{
	CategoryPersonal,
	CategoryStudy,
	CategoryWork,
	CategoryOther,
}

func IsValidCategory(category string) bool {
	for _, c := range Categories {
		if c == category {
			return true
		}
	}
	return false
}

func GetDefaultCategory() string {
	return CategoryPersonal
}