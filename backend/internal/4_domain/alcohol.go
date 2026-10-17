package domain

import (
	"context"
	"fmt"
	"github.com/YujiYabe/restaurant-order-server/backend/internal/pkg"
	"strconv"
	"strings"
	"time"
)

// AlcoholList は、酒類商品のJANコード一覧です。
type AlcoholList []int

// TwentyList は、利用者UUIDとその有効期限を保持するマップです。
// AllowedTwentyListとPendingFaceListは、この型を別々のマップとして使用します。
type TwentyList map[string]time.Time

var japanLocation = time.FixedZone("Asia/Tokyo", 9*60*60)

// NewAlcoholList は、空の酒類商品JANコード一覧を生成します。
func NewAlcoholList() (
	alcoholList AlcoholList,
) {
	alcoholList = AlcoholList{}
	return
}

// NewTwentyList は、空の利用者UUID・有効期限マップを生成します。
func NewTwentyList() (
	twentyList TwentyList,
) {
	twentyList = make(TwentyList)
	return
}

// IsIncludeTwentyList は、利用者UUIDが成人確認済みのAllowedTwentyListに含まれ、
// かつ有効期限内であるかを返します。空のUUIDと期限切れのUUIDはfalseとなり、
// 期限切れのエントリはAllowedTwentyListから削除されます。
func (receiver *Domain) IsIncludeTwentyList(
	ctx context.Context,
	alcohol Alcohol,
) (
	ok bool,
) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()
	if alcohol.UUID == "" {
		return
	}

	expiresAt, exists := receiver.AllowedTwentyList[alcohol.UUID]
	if !exists {
		return
	}
	if !time.Now().Before(expiresAt) {
		delete(receiver.AllowedTwentyList, alcohol.UUID)
		return
	}

	ok = true
	return
}

// AddTwentyList は、顔認証を完了した利用者UUIDをAllowedTwentyListへ追加します。
// 有効期限は日本時間の翌日0時です。空のUUIDは追加しません。
func (receiver *Domain) AddTwentyList(
	ctx context.Context,
	alcohol Alcohol,
) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()
	if alcohol.UUID == "" {
		return
	}
	if receiver.AllowedTwentyList == nil {
		receiver.AllowedTwentyList = NewTwentyList()
	}

	now := time.Now().In(japanLocation)
	expiresAt := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	receiver.AllowedTwentyList[alcohol.UUID] = expiresAt
}

// AllowFaceRegistration は、書類による成人確認を完了した利用者UUIDに対して、
// 10分間有効な1回限りの顔画像登録許可を発行します。空のUUIDには発行しません。
func (receiver *Domain) AllowFaceRegistration(
	_ context.Context,
	alcohol Alcohol,
) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()
	if alcohol.UUID == "" {
		return
	}
	if receiver.PendingFaceList == nil {
		receiver.PendingFaceList = NewTwentyList()
	}
	receiver.PendingFaceList[alcohol.UUID] = time.Now().Add(10 * time.Minute)
}

// ConsumeFaceRegistration は、利用者UUIDに発行された顔画像登録許可を1回だけ消費します。
// 有効な許可を消費した場合はtrueを返します。許可が存在しない場合や期限切れの場合は
// falseを返し、期限切れのエントリをPendingFaceListから削除します。
func (receiver *Domain) ConsumeFaceRegistration(
	_ context.Context,
	alcohol Alcohol,
) bool {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()
	expiresAt, exists := receiver.PendingFaceList[alcohol.UUID]
	if !exists || !time.Now().Before(expiresAt) {
		delete(receiver.PendingFaceList, alcohol.UUID)
		return false
	}
	delete(receiver.PendingFaceList, alcohol.UUID)
	return true
}

// IsAlcohol は、注文のJANコードがAlcoholListに登録された酒類商品かどうかを返します。
func (receiver *Domain) IsAlcohol(
	ctx context.Context,
	order Order,
) (
	ok bool,
) {
	receiver.mutex.RLock()
	defer receiver.mutex.RUnlock()
	for _, alcohol := range receiver.AlcoholList {
		if alcohol == order.JANCode {
			ok = true
			return
		}
	}
	return
}

// VerifyTwenty は、OCRから取得した和暦の生年月日を解析し、現在日時時点の年齢を計算します。
// 明治、大正、昭和、平成、令和に対応します。形式が不正な場合、未対応の元号の場合、
// または20歳未満の場合はエラーを返します。
func (receiver *Domain) VerifyTwenty(
	ctx context.Context,
	birthdayString string,
) (
	age int,
	birthday time.Time,
	err error,
) {
	birthday = time.Now()

	// 不要なスペース削除
	birthdayString = strings.Join(
		strings.Fields(birthdayString),
		"",
	)

	enaNameMap := map[string]int{
		"明治": 1873,
		"大正": 1912,
		"昭和": 1926,
		"平成": 1989,
		"令和": 2019,
	}

	enaName := ""
	calendar := 0
	found := false

	for key, value := range enaNameMap {
		if strings.Contains(birthdayString, key) {
			enaName = key
			calendar = value
			found = true
			break
		}
	}

	if !found {
		err = fmt.Errorf("unknown enaName")
		pkg.Logging(ctx, err)
		return
	}

	const splitString = "："
	replaceCandidate := strings.NewReplacer(
		enaName, "",
		"年", splitString,
		"月", splitString,
		"日", "",
		"生", "",
	)

	var splitYYMMDD = strings.Split(
		replaceCandidate.Replace(birthdayString),
		splitString,
	)
	if len(splitYYMMDD) != 3 || splitYYMMDD[0] == "" || splitYYMMDD[1] == "" || splitYYMMDD[2] == "" {
		err = fmt.Errorf("invalid birthday format")
		pkg.Logging(ctx, err)
		return
	}

	enaNameYY, err := strconv.Atoi(splitYYMMDD[0])
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}
	dateStr := ""
	dateStr += strconv.Itoa(calendar + enaNameYY - 1) // 年
	dateStr += "/" + splitYYMMDD[1]                   // 月
	dateStr += "/" + splitYYMMDD[2]                   // 日

	birthday, err = time.Parse("2006/1/2", dateStr)
	if err != nil {
		pkg.Logging(ctx, err)
		return
	}

	timeNow := time.Now()
	age = timeNow.Year() - birthday.Year()

	// if timeNow.Month() < birthday.Month() || (timeNow.Month() == birthday.Month() && timeNow.Day() < birthday.Day()) {
	if timeNow.Before(birthday.AddDate(age, 0, 0)) {
		age--
	}

	if age < 20 {
		err = fmt.Errorf("age is under 20")
		pkg.Logging(ctx, err)
		return
	}

	return
}
