const yyyymmdd = new Intl.DateTimeFormat(
  undefined,
  {
    year:   'numeric',
    month:  '2-digit',
    day:    '2-digit',
    hour:   '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }
)

const placeStatusCleaning = 0; // 掃除中
const placeStatusOpening  = 1; // 受付中 店頭で席を確保する
const placeStatusUsing    = 2; // 利用中 オーダー可能状況
const placeStatusBilling  = 3; // 支払中
const placeStatusStopping = 4; // 停止中

const placeStatusMap = { 
  0 : "placeStatusCleaning", 
  1 : "placeStatusOpening", 
  2 : "placeStatusUsing", 
  3 : "placeStatusBilling", 
  4 : "placeStatusStopping", 
}
const orderStatusCancel    = 0; // 
const orderStatusPreparing = 1; // 
const orderStatusPassed    = 2; // 

const textOrderStatusCancel    = '取り消し';   // 
const textOrderStatusPreparing = '準備中';     // 
const textOrderStatusPassed    = 'お届け済み'; // 

const colorOrderStatusCancel    = 'bg-black    btn-black ';  //
const colorOrderStatusPreparing = 'bg-success  btn-success'; // 
const colorOrderStatusPassed    = 'bg-primary  btn-primary'; // 

let isEnvLocal;

let chargeBaseCapacity = 1;
let chargePrice = 0;


const days = ["日", "月", "火", "水", "木", "金", "土"];

async function commonInit () {
  version = $('#version').attr('name');
  isEnvLocal = $('#environment').attr('name') == 'local';

  if (!isEnvLocal) {
    // 右クリック禁止
    $(document).on('contextmenu', function (e) {
      return false;
    });

    // マウスカーソル非表示
    // $('html').css('cursor', 'none');

    // テキスト選択不可
    $('body').css('-webkit-user-select', 'none'); /* Safari */
    $('body').css('-moz-user-select', 'none'); /* Firefox */
    $('body').css('-ms-user-select', 'none'); /* IE 10+ */
    $('body').css('user-select', 'none');

    // ピンチアウト禁止
    document.addEventListener('touchmove', mobileNoScroll, {
      passive: false,
    });
  }
}


function mobileNoScroll(event) {
  if (event.touches.length >= 2) { // ２本指での操作の場合
    event.preventDefault();  // デフォルトの動作をさせない
  }
}



// SHA-256ハッシュを計算する関数
async function calculateSHA256Hash(place) {
  const encoder = new TextEncoder();
  const dataBuffer = encoder.encode(place);
  const hashBuffer = await crypto.subtle.digest('SHA-256', dataBuffer);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  const hashHex = hashArray
    .map((byte) => byte.toString(16).padStart(2, '0'))
    .join('');
  return hashHex;
}

// https://www.relief.jp/docs/001403.html

function printDebug(title, value) {
  if (isDebug()) {
    console.log(' = ' +yyyymmdd.format( new Date() ) +' = ' + title + ' ======================== ');
    console.log(value);
  }
}

function isDebug() {
  if ('?d' == window.location.search) {
    return true;
  }

  return false;
}

function getParam(name, url) {
  if (!url) url = window.location.href;
  name = name.replace(/[\[\]]/g, '\\$&');
  let regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)'),
    results = regex.exec(url);
  if (!results) return null;
  if (!results[2]) return '';
  return decodeURIComponent(results[2].replace(/\+/g, ' '));
}


function calculateTime() {
  if ($('#status').attr('name') != placeStatusUsing) {
    return;
  }
  const chargeBasePrice = $('#charge_base_price').attr('name');
  const chargeBaseCapacity = $('#charge_base_capacity').attr('name');
  const chargeBaseMinute = $('#charge_base_minute').attr('name');

  // ----------------------------------------
  let nowTime = new Date();

  // ----------------------------------------
  let diffTime = nowTime.getTime() - enterTime.getTime();
  let diffHour = diffTime / (1000 * 60 * 60);
  let diffMinute = (diffHour - Math.floor(diffHour)) * 60;
  let diffSecond = (diffMinute - Math.floor(diffMinute)) * 60;


  let diffBaseHour = Math.floor(diffHour) * 6
  // let diffBaseMinute = 1 + Math.floor(Math.floor(diffMinute) / chargeBaseMinute);
  let diffBaseMinute = Math.floor(Math.floor(diffMinute) / chargeBaseMinute);
  let diffBase = diffBaseHour + diffBaseMinute;
  if (!isEnableCharge) { 
    diffBase  = 0;
  }


  chargePrice = diffBase * chargeBasePrice * chargeBaseCapacity;
  let chargeUnit = Math.floor(
    placePrice / (chargeBasePrice * chargeBaseCapacity)
  );

/*
  現在時間   開始時間   チャージ金額    会計金額   ステータス  アクション
  20:01      20:01〜    400                    0    初期状態    未購入 
  20:02      20:01〜    200                  200    400円未満   200購入
  20:03      20:01〜    20:11までは無料      500    400円以上   300購入
  20:12      20:01〜    300                  500    800円未満   時間経過
  20:13      20:01〜    20:31までは無料      900    800円以上   400購入
  20:14      20:01〜    20:41までは無料     1400   1200円以上   500購入
  20:42      20:01〜    200                 1400   1600円未満   時間経過

  チャージ金額を計算
  会計金額を計算

  会計金額＞チャージ金額 ＆＆ 現在時刻 ＞ チャージ時刻 
  チャージ金額分の時間まで無料にする

  会計金額＞チャージ金額 ＆＆ 現在時刻 ＜ チャージ時刻 
  残りチャージ料金を表示

  会計金額＜チャージ金額
  残りチャージ料金を表示
*/

  const textFree = 'まで無料';
  $('#debug_table').empty();

  let debugTable = $('<table>')
    .addClass('table table-striped table-bordered')
    .appendTo($('#debug_table'));
  let debugTr1 = $('<tr>').appendTo(debugTable);
  let debugTr2 = $('<tr>').appendTo(debugTable);

  $('<td>').appendTo(debugTr1).text('enterTime');
  $('<td>').appendTo(debugTr1).text('nowTime');
  $('<td>').appendTo(debugTr1).text('diffTime');
  $('<td>').appendTo(debugTr1).text('chargePrice');
  $('<td>').appendTo(debugTr1).text('placePrice');
  $('<td>').appendTo(debugTr1).text('diffBase');
  $('<td>').appendTo(debugTr1).text('extendTime');

  // ---------------------------------
  let textEnter = $('<td>').appendTo(debugTr2);
  $('<div>').appendTo(textEnter).text(textDateEnter);
  $('<div>').appendTo(textEnter).text(textTimeEnter);

  
  // ---------------------------------
  let [textDateNow, textTimeNow] = convertTimeString(nowTime);
  let textNow = $('<td>').appendTo(debugTr2);
  $('<div>').appendTo(textNow).text(textDateNow);
  $('<div>').appendTo(textNow).text(textTimeNow);

  // ---------------------------------
  let textDiff = $('<td>').appendTo(debugTr2);
  const diffDays = Math.floor(diffHour / 24);
  const diffHH = ('00' + Math.floor(diffHour % 24)).slice(-2);
  const diffMM = ('00' + Math.floor(diffMinute)).slice(-2);
  const diffSS = ('00' + Math.round(diffSecond)).slice(-2);
  $('<div>').appendTo(textDiff).text(diffDays + 'd ' + diffHH + ':' + diffMM + ':' + diffSS);

  $('<td>').appendTo(debugTr2).text(chargePrice);
  $('<td>').appendTo(debugTr2).text(placePrice);
  $('<td>').appendTo(debugTr2).text(diffBase);

  // ---------------------------------
  let extendTime = new Date();
  extendTime.setTime(
    enterTime.getTime() + chargeUnit * chargeBaseMinute * 60 * 1000
  );

  let [textDateExtended, textTimeExtended] = convertTimeString(extendTime);
  let textExtended = $('<td>').appendTo(debugTr2);
  $('<div>').appendTo(textExtended).text(textDateExtended);
  $('<div>').appendTo(textExtended).text(textTimeExtended);

}

function convertTimeString (targetTime) {
  const YYYYMMDD = targetTime.toLocaleDateString('sv-SE');

  const HH = ('00' + targetTime.getHours()).slice(-2);
  const MM = ('00' + targetTime.getMinutes()).slice(-2);
  const SS = ('00' + targetTime.getSeconds()).slice(-2);

  return [YYYYMMDD , HH + ':' + MM + ':' + SS];
}

async function getLocalStorage(key) {
  // local storage にあれば使う
  if (
    localStorage.getItem(key) == undefined ||
    localStorage.getItem(key) == null ||
    localStorage.getItem(key) == ""
  ) { 
    return "";
  }

  if (localStorage.hasOwnProperty(key)) {
    return localStorage.getItem(key);
  }

  return "";
}

async function setLocalStorage(key, value) {
  localStorage.setItem(key, value);
}

function countMinus() {
  let count = Number($('#count').text());
  if (count > 1) {
    count--;
    $('#count').text(count);
    $('#product_count').text(count);

    $('#total_product_price').text(count * selectedPriceTotal + '円');
  }
}

function countPlus() {
  let count = Number($('#count').text());
  count++;

  if (inventoryMap[selectedJANCode] < count) {
    return;
  }

  $('#count').text(count);
  $('#product_count').text(count);
  $('#total_product_price').text(count * selectedPriceTotal + '円');
}

function changePage(pageName) {
  $('.page-title').removeClass(' border-danger border-2');
  $('.page-body').hide();

  $('#' + pageName).addClass(' border-danger border-2');
  $('#body_' + pageName).show();
}

function filterProductByTag(tagSecondaryID) {
  $('#stock_none').hide();
  $('.products').hide();
  const keyWord = tagSecondaryList[tagSecondaryID];
  let isExist = false
  $.each(miniProductMap, function (_, product) {
    if (product.tag.includes(keyWord)) {
      if (inventoryMap[product.jan_code] != 0) {
        $('#' + product.jan_code).show();
        isExist = true
      }
    }
  });
  if (!isExist) {
    $('#stock_none').show();
  }
}
