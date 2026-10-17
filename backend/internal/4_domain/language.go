package domain

type (
	// Language は、言語設定をドメインモデル化するための予約型です。
	// 現在はフィールドを持たず、有効言語の検証はStock.VerifyJANCodeListが担当します。
	Language struct {
		// defaultLangCode     int
		// isValidLangCodeList []int
		// isValidLangCodeMap  map[int]string
	}
)
