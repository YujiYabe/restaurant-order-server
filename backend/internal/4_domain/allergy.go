package domain

type (
	// Allergy は、アレルゲン識別子と表示名の対応を保持します。
	Allergy map[string]string
	// AllergyList は、言語などが異なる複数のアレルギー定義を保持します。
	AllergyList []Allergy
)

// NewAllergyDefault は、日本語の既定アレルゲン名を含むマップを生成します。
func NewAllergyDefault() (
	allergy Allergy,
) {
	allergy = newAllergyListJa()
	return
}

// newAllergyListJa は、日本語ロケールと表示対象のアレルゲン名を定義します。
// 参考: https://www.relief.jp/docs/001403.html
func newAllergyListJa() (
	allergy Allergy,
) {
	allergy = Allergy{
		"lang_type": "1041", // 日本語

		// 乳卵
		"egg":  "卵",
		"milk": "乳",

		// 魚介
		"shrimp":     "えび",
		"crab":       "かに",
		"squid":      "いか",
		"salmon":     "さけ",
		"mackerel":   "さば",
		"salmon_roe": "いくら",
		"abalone":    "あわび",

		// 種子
		"peanut":     "ピーナッツ",
		"walnut":     "くるみ",
		"soybean":    "大豆",
		"almond":     "アーモンド",
		"cashew_nut": "カシューナッツ",
		"sesame":     "ごま",

		// 穀物
		"wheat":     "小麦",
		"buckwheat": "そば",

		// フルーツ
		"apple":      "りんご",
		"orange":     "オレンジ",
		"banana":     "バナナ",
		"peach":      "もも",
		"kiwi_fruit": "キウイフルーツ",

		// 肉
		"pork":    "豚肉",
		"chicken": "鶏肉",
		"beef":    "牛肉",

		// その他
		"wild_yam":           "やまいも",
		"matsutake_mushroom": "まつたけ",
		"gelatin":            "ゼラチン",
	}

	return
}
