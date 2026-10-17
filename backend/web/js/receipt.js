

async function createReceiptContent (place, diffPrice) {
	let bodyString = "";
    const newLine = '\n';

	// タイトル  ---------------------
	const title = "{w:*;b:line}" + newLine +
					"^^^取  引  明  細" + newLine +
					"{w:auto;b:space}"
	bodyString = bodyString + title + newLine

	// 店名 ---------------------
	const storeName= $('#store_name').attr('name');
	bodyString = bodyString + storeName + newLine

	// 店住所 ---------------------
	const storeAddress= $('#store_address').attr('name');
	bodyString = bodyString + storeAddress + newLine
	
	
	// 利用時間 ---------------------
	const enterTime = new Date(place.enter_time)
	
	const [textDateEnter, textTimeEnter] = convertTimeString(enterTime);
	const [textDateExit, textTimeExit] = convertTimeString(new Date());
	if (textDateEnter != textDateExit) { 
		// 日付が違う場合はtextDateExitを24+Nの値にする
	}
	const usedSpan = "ご利用時間 " + textDateEnter + "("+ days[enterTime.getDay()] + ") " + textTimeEnter.slice(0, 5) + "〜" + textTimeExit.slice(0, 5) ;
	bodyString = bodyString + usedSpan + newLine

	// 伝票コード ---------------------
    const Y = enterTime.getFullYear();
    const M = ('00' + (enterTime.getMonth() + 1)).slice(-2);
    const D = ('00' + enterTime.getDate()).slice(-2);

	const StoreID= $('#store_id').attr('name');
	const receiptCode = "明細コード " + Y+M+D + "-" + StoreID + "-" + place.place_name + "-" +place.line_number
	bodyString = bodyString + receiptCode + newLine

	// 購入一覧 ---------------------
	let sumPrice = 0 ;
	let sumCount = 0 ;
	let mapJANCodeCount = new Map(); 
	await $.each(place.order_list, function (_, order) {
		mapJANCodeCount[order.jan_code] = (mapJANCodeCount[order.jan_code] || 0) + order.count;
	}); 

	let consumeList = "{w:*,4,8;b:line}" + newLine;
	consumeList = consumeList + " 品名 | 数量 | 金額" + newLine + "-" + newLine;
	
	$.each(mapJANCodeCount, function (janCode, count) {
		let price = mapJANCodePrice[janCode] *count;

		sumPrice = sumPrice + price;
		sumCount = sumCount + count;
		const shortName = ""
		consumeList = consumeList + shortName + " | " + count + "| ¥" + price + newLine;
		
	});
	// チャージ料金
	consumeList = consumeList + "チャージ料金" + " | "+ "1" + "| ¥" + diffPrice + newLine;
	sumPrice = sumPrice + diffPrice;
	sumCount= sumCount + 1;

	// 計 ------
	consumeList = consumeList + "-" + newLine + "　" + "| "+ sumCount+ "| ¥" + sumPrice + newLine;
	consumeList = consumeList + "{w:auto;b:space}";
	bodyString = bodyString + consumeList + newLine;


	// 集計 ---------------------
	bodyString = bodyString + "-" + newLine;

	let tax = Math.round((sumPrice /1.1*0.1));
	let summary = "{w:*,*}" + newLine + "|^^^ 合計 | ^^^¥" + sumPrice + newLine ;
	let subSummary = "|　 内消費税率等 10% | ¥" + tax + newLine + newLine + "{w:auto;b:space}";

	bodyString = bodyString + summary + subSummary + newLine ;

	// 注意事項 ---------------------
	bodyString = bodyString +"-" + newLine + "以下のバーコードにてセルフレジでお会計願います" + newLine + "-" + newLine ;



	// nonPLU バーコード ---------------------
	// 0123456789012 | 13桁
	// 0212345309803 | サンプル
	// 02000XYYYYYYZ | 最後の1桁は強制

	let nonPLUFlag = "02" ; // 固定
	let nonPLUCode = "001" ; // 店id
	let priceCheckDigit = "" ; // X
	let nonPLUPrice = ('0000000' + sumPrice).slice(-7) ;// YYYYYY

	let barCode = "{code:" + nonPLUFlag + nonPLUCode + priceCheckDigit + nonPLUPrice +"; option:jan,2,72,hri}";
	bodyString = bodyString + barCode + newLine + "-" + newLine + newLine;

	return bodyString;
}
