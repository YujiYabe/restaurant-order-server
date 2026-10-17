package domain

// NewTagPrimaryList は、空の主タグ一覧を生成します。
func NewTagPrimaryList() (
	tagPrimaryList TagPrimaryList,
) {
	tagPrimaryList = TagPrimaryList{}
	return
}

// TagPrimaryList は、商品分類に使用する主タグ一覧です。
type TagPrimaryList []TagPrimary

// TagPrimary は、主タグとその配下の副タグ一覧を保持します。
type TagPrimary struct {
	ID               int
	Name             string
	TagSecondaryList TagSecondaryList
}

// TagSecondaryList は、主タグに属する副タグ一覧です。
type TagSecondaryList []TagSecondary

// TagSecondary は、副タグと所属する主タグIDを保持します。
type TagSecondary struct {
	ID           int
	TagPrimaryID int
	Name         string
}
