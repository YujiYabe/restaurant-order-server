// グローバル変数
let selectedTagList = {};
let tagSecondaryList = [];
let selectedLanguageCode = 1041;
let selectedPriceTotal = 0;
let selectedJANCode;
let isAlreadyCreateTabletScreen = false;
let inventoryMap;
let enterTime;
let placePrice = 0;
let diffPrice = 0;
let defaultScreen = 'mobile'
// let defaultScreen = 'tablet'
// let defaultScreen = 'history'
let textDateEnter;
let textTimeEnter;
let order;
let miniProductMap
let isEnableCharge = false;
let isMiniProductMapOK = false;
let deviceKey;
let chargeBaseMinute;

// 初期化関数
async function init() {
  if (isEnvLocal) {
    $('#debug_table').show();
  }
  debug = getParam('debug');

  if (debug) {
    $('#debug_table').show();
  }
  if ( $('#is_enable_charge').attr('name') == "true") {
    isEnableCharge = true;
  }

  chargeBaseMinute = $('#charge_base_minute').attr('name') 

  if (version != await getLocalStorage('version')) { 
    await setLocalStorage('version', version);
    resetConfig();
  }
  
  miniProductMap = JSON.parse($('#mini_product_map').text());

  attachQRCode();
  buildProductList();
  deviceKey = getParam('device_key');


  // setInterval('calculateTime()', 1000);
}

async function buildOnSse(eventData) {
  await new Promise((resolve) => {
    const intervalId = setInterval(() => {
      if (isMiniProductMapOK) {
        clearInterval(intervalId);
        buildWaiterMonitor(eventData);
        resolve();
      }
    }, 100);
  });
}

async function hideAll() {
  $('.screen').hide();
}

async function resetConfig () {
  localStorage.removeItem("selectedTagPrimary");
  localStorage.removeItem("selectedTagList");
  location.reload()
}

async function buildProductList() {

  $.each(miniProductMap, function (_, product) {

    const divParent = $('<div>') //
      .appendTo('#menu_content') //
      .addClass(' col-6 alert products ') //
      .hide() //
      .attr('tag', product.tag) //
      .attr('id', product.jan_code)
      .on('click', function () {
        showConfirmOrder(async () => {
          await postProduct();
        }, product); //
      });
    const rounded = $('<div>')
      .addClass(' rounded-5 alert alert-light background-white ') //
      .appendTo(divParent); //

    $('<img>') //
      .appendTo(rounded) //
      .addClass('img-fluid ') //
      .attr(
        'src',
        '/image/product.svg?item=' +
          product.jan_code +
          '.jpg?version=' + version
      );

      
    $('<div>') //
      .addClass('bg-white h-2') //
      .appendTo(rounded) //
      .text( '(' +  product.price + '円)' + product.name_ja ); //
  });
  isMiniProductMapOK = true
}

async function buildWaiterMonitor(eventData) {
  if (typeof eventData === 'string') {
    tempData = JSON.parse(eventData);
  } else {
    tempData = eventData; 
  }

  console.log(" ========================= ");
  console.log( tempData );
  console.log(" ========================= ");
  
  
  const place = tempData.Place;

  inventoryMap = tempData.InventoryMap;
  printDebug('tempData', tempData); // -----------------------------

  $('#permission').attr('name', place.permission);
  $('#status').attr('name', place.place_status);
  $('#charge_base_capacity').attr('name', place.capacity);

  
  $.each(inventoryMap, function (key, value) {
    if (value == 0) {
      $('#' + key).remove();
    }
    if (selectedJANCode == key) {
      let count = Number($('#count').text());
      if (value < count) {
        $('#count').text(value);
        $('#product_count').text(value);
        $('#total_product_price').text(value * selectedPriceTotal + '円');
      }
      if (value == 0) {
        $('#confirm_button_product').text('品切れ');
        $('#confirm_button_product').addClass('btn-secondary');
      }
    }
  });

  hideAll();
  
  printDebug('place_status', placeStatusMap[place.place_status]); // -----------------------------
  switch (place.place_status) {
    case placeStatusCleaning:
      $('#body_cleaning').show();
    break;

    case placeStatusOpening:
      $('#body_opening').show();
      resetTagPrimary(1);
      resetTagSecondaryGroup(1);
      // resetTagSecondary(1, 1);
      resetTagSecondary(1, 1);
      setIsRequireReload(1); // set true
    break;

    case placeStatusUsing:
      if (await getIsRequireReload() == 1 ){ // 状態変更 1だけリロードする 
        await setIsRequireReload(0); // set false
        location.reload();
      }

      attachQRCode();
      const placeName = $('#place_name').attr('name');
      $('#body_using').show();
      $('#line_number').text(placeName + '-' + place.line_number);
      enterTime = new Date(place.enter_time);
      [textDateEnter, textTimeEnter] = convertTimeString(enterTime);
      $('#enter_time').text(textTimeEnter.slice( 0, 5 ) + '〜');

      let placePrice = 0;
      $.each(place.order_list, function (_, order) {
        if (order.order_status != orderStatusCancel) {
          placePrice = placePrice + (order.count * miniProductMap[order.jan_code].price);
        }
      });

      // ----------------------------------------
      let nowTime = new Date();
      const chargeBasePrice = $('#charge_base_price').attr('name');
      const chargeBaseCapacity = place.capacity;
      let diffBase;
      if (isEnableCharge) { 
        diffBase  = place.base_count;
      }
      


      const chargePrice = diffBase * chargeBasePrice * chargeBaseCapacity;
      let chargeUnit = Math.floor(
        placePrice / (chargeBasePrice * chargeBaseCapacity)
      );

      // ---------------------------------
      let extendTime = new Date();
      extendTime.setTime(
        // enterTime.getTime() + chargeUnit * chargeBaseMinute * 60 * 1000
        enterTime.getTime() + (chargeUnit +1) * chargeBaseMinute * 60 * 1000

      );
      let [textDateExtended, textTimeExtended] = convertTimeString(extendTime);

      // チャージ料金より購入金額が多い
      // 時間が延長時間未満
      let isPlaceOverCharge = placePrice > chargePrice;
      let isNowOverExtend = nowTime > extendTime;
      diffPrice = chargePrice - placePrice;

      if (isPlaceOverCharge && !isNowOverExtend) {
        $('#diff_price').text(textTimeExtended.slice(0, 5) + 'まで無料');
      } else {
        $('#diff_price').text(diffPrice + '円');
      }
      $('#place_price').text(placePrice + '円');


      if (isAlreadyCreateTabletScreen == false) {
        await createTabletScreen(place);
        resetTagPrimary(1);
        resetTagSecondaryGroup(1);
        // resetTagSecondary(1, 1);
        resetTagSecondary(1, 1);
      }

      await createHistoryScreen(place);

    break;

    case placeStatusBilling:
      $('#body_billing').show();
    break;

    case placeStatusStopping:
      $('#body_stopping').show();
    break;

    default:
      console.error('error status : ', value);
    break;
  }
}

async function createTabletScreen(place) {
  $('#check').on('click', function () {
    showConfirmCheck(async () => {
      await closeOrder(place);
    }); //
  });

  $.each(
    JSON.parse($('#tag_list').text()),
    function (_, tagPrimary) {
    selectedTagList[tagPrimary.id] = 0;

    $('<div>')
      .appendTo('#tag_primary') //
      .addClass(
        'tag_primary alert alert-light background-white bs-primary-bg-subtle rounded-pill'
      ) //
      .attr('id', 'tag_primary_' + tagPrimary.id) //
      .text(tagPrimary.name) //
      .on('click', function () {
        resetTagPrimary(tagPrimary.id);
        resetTagSecondaryGroup(tagPrimary.id);
      }); //

    const tagPrimaryGroupID = 'tag_primary_group_' + tagPrimary.id;
    let divTagSecondary = $('<div>')
      .appendTo('#tag_secondary') //
      .hide() //
      .addClass('tag_primary_group') //
      .attr('id', tagPrimaryGroupID); //
      
    let isFirstSecondaryTag = true
    $.each(tagPrimary.tag_secondary_list, function (_, tagSecondary) {
      tagSecondaryList[tagSecondary.id] = tagSecondary.name;

      let textName = tagSecondary.name
      let replaceName = "全て"
      if (textName.includes(replaceName)) {
        textName = replaceName
      }

      $('<div>')
        .appendTo(divTagSecondary) //
        .addClass(
          'tag_secondary alert alert-light background-white bs-primary-bg-subtle rounded-pill ' +
            tagPrimaryGroupID 
        )
        .attr('id', 'tag_secondary_' + tagSecondary.id) //
        .text(textName) //
        .on('click', function () {
          resetTagSecondary(tagPrimary.id, tagSecondary.id);
        }); //

      if (isFirstSecondaryTag ){
        isFirstSecondaryTag = false
        resetTagSecondary(tagPrimary.id, tagSecondary.id);
      }
    });
  });

  isAlreadyCreateTabletScreen = true;
  changePage(defaultScreen);
}

async function createHistoryScreen(place) {
  placePrice = 0;
  let reverseOrderList = [];

  const divOrder = $('<div>').addClass('row');
  const divMenu = $('<div>').addClass('row');

  await $.each(place.order_list, function (_, order) {
    reverseOrderList.push(order);
  });
  reverseOrderList.sort().reverse();

  const divPreparing = $('<div>').appendTo(divOrder).addClass('row');
  $('<hr>').appendTo(divOrder);
  const divPassed = $('<div>').appendTo(divOrder).addClass('row');

  $.each(reverseOrderList, function (_, order) {
    if (order.order_status != orderStatusCancel){
      placePrice = placePrice + (order.count * miniProductMap[order.jan_code].price) ;
    }

    const divCard = $('<div>').addClass('col-6 alert mb-3 ');
    const divRow = $('<div>')
      .appendTo(divCard)
      .addClass('row rounded-5 alert alert-light background-white ');
    const divCardLeft = $('<div>').appendTo(divRow).addClass('col-md-6');

    $('<img>')
      .appendTo(divCardLeft)
      .addClass('img-fluid rounded-start ')
      .attr(
        'src',
        '/image/product.svg?item=' +
          order.jan_code +
          '.jpg?version=' + version
      );

    let rightSpace = $('<div>')
      .appendTo(divRow)
      .addClass('col-md-6 text-center text-dark font-select ')

    $('<div>')
      .appendTo(rightSpace)
      .addClass('text-center text-dark font-select text-count')
      .text(order.count);

    if (miniProductMap[order.jan_code].recipe != "") { 
      $('<hr>').appendTo(rightSpace);
      $('<div>')
        .appendTo(rightSpace)
        .addClass('text-center text-dark font-select text-recipe')
        .text(miniProductMap[order.jan_code].recipe);
    }

    let status = '';
    let statusColor = '';

    switch (order.order_status) {
      case orderStatusCancel:
        status = textOrderStatusCancel;
        statusColor = 'bg-secondary';
        break;
      case orderStatusPreparing:
        status = textOrderStatusPreparing;
        statusColor = colorOrderStatusPreparing;
        $('<div>')
          .appendTo(divCard)
          .addClass(
            ' border-0 badge align-items-center font-20px position-absolute top-0 start-0 ' +
              statusColor
          )
          .text(status);
        // divCard.appendTo(divPreparing);
        const divCardHistory = divCard.clone();
        divCardHistory.addClass('col-6 ').appendTo(divPreparing);

        const divCardMenu = divCard.clone();
        divCardMenu.addClass('col-12 ').appendTo(divMenu);

        break;
      case orderStatusPassed:
        status = textOrderStatusPassed;
        statusColor = colorOrderStatusPassed;
        $('<div>')
          .appendTo(divCard)
          .addClass(
            ' border-0 badge align-items-center font-20px position-absolute top-0 start-0 ' +
              statusColor
          )
          .text(status);
        divCard.appendTo(divPassed);

        break;
      default:
        console.error('error status : ', value);
        break;
    }
  });

  //

  await $('#order_history').empty();
  divOrder.appendTo('#order_history');
  $('#order_history').find('.text-count').addClass('font-50px');

  await $('#menu_util').empty();
  divMenu.appendTo('#menu_util');
  $('#menu_util').find('.text-count').addClass('font-30px');
}





async function getIsRequireReload() {
  return localStorage.getItem('is_require_reload');
}

async function setIsRequireReload(value) {
  await localStorage.setItem('is_require_reload', value);
}


function resetTagPrimary(tagPrimaryID) {
  $('.tag_primary').removeClass(' border-danger border-2');
  $('#tag_primary_' + tagPrimaryID).addClass(' border-danger border-2');
  $('.products').hide();

  if (selectedTagList[tagPrimaryID] != 0) {
    filterProductByTag(selectedTagList[tagPrimaryID]);
  }
}



function resetTagSecondaryGroup(tagPrimaryID) {
  $('.tag_primary_group').hide();
  $('#tag_primary_group_' + tagPrimaryID).show();
}

function resetTagSecondary(tagPrimaryID, tagSecondaryID) {
  const tagPrimaryGroupID = 'tag_primary_group_' + tagPrimaryID;
  $('.' + tagPrimaryGroupID).removeClass(' border-danger border-2');
  $('#tag_secondary_' + tagSecondaryID).addClass(' border-danger border-2');
  selectedTagList[tagPrimaryID] = tagSecondaryID; //
  selectedTagSecondary = tagSecondaryID;
  filterProductByTag(tagSecondaryID);
}


async function showConfirmOrder(callback, product) {
  selectedJANCode = product.jan_code;
  selectedPriceTotal = product.price;

  const modal = new bootstrap.Modal(document.getElementById('modal'));
  $('#confirm_button_product').text('注文する');
  $('#confirm_button_product').removeClass('btn-secondary');

  $('#modal-size') //
    .addClass('modal-xl'); //

  $('#count').text(1);
  $('#product_count').text(1);

  $('#product_price').text(selectedPriceTotal + '円');
  $('#total_product_price').text(selectedPriceTotal + '円');

  const divParent = $('<div>');
  const divCard = $('<div>').appendTo(divParent);
  const divRow = $('<div>').appendTo(divCard).addClass('row');
  const divCardLeft = $('<div>').appendTo(divRow).addClass('col-md-3');
  $('<img>')
    .appendTo(divCardLeft)
    .addClass('img-fluid rounded-start ')
    .attr(
      'src',
      '/image/product.svg?item=' +
        product.jan_code +
        '.jpg?version=' + version
    );

      
  const divCardRight = $('<div>').appendTo(divRow).addClass('col-md-9');
  const divCardBody = $('<div>').appendTo(divCardRight).addClass('card-body');

  // 名前
  $('#modal_header_text')
    .addClass('text-center') //
    .text(product.name_ja); //

  if (product.recipe != "") { 
    // レシピ
    $('<div>')
      .appendTo(divCardBody)
      .addClass('h4 card-text')
      .text(product.recipe); //
    $('<hr>')
      .appendTo(divCardBody)
  }

  // 説明
  $('<div>')
    .appendTo(divCardBody)
    .addClass('h6 card-text')
    .text(product.information_ja);

  $('<div>')
    .appendTo(divCardBody)
    .addClass('h6 card-text')
    .attr('id', 'calorie_' + product.jan_code);

  const li_allergy = $('<div>').appendTo(divCardBody).addClass('h6 card-text');
  $('<div>')
    .appendTo(li_allergy)
    .addClass('badge bg-danger my-2')
    // .text(allergyName);
    .text('アレルギー情報');

  const cardTop = $('<div>')
    .appendTo(li_allergy)
    .addClass('row mx-auto col-12 ');
  const allergyList = JSON.parse(product.allergy);

  for (const key in allergyList) {
    if (allergyList[key]) {
      const card = $('<div>')
        .appendTo(cardTop)
        .addClass('card col-1 width_18rem');
      $('<img>')
        .appendTo(card)
        .addClass('img-thumbnail card-img-top')
        .attr('src', '/image/allergy.svg?name=' + encodeURIComponent(key));

      const cardBody = $('<div>').appendTo(card).addClass('text-center');

      $('<p>')
        .appendTo(cardBody)
        .addClass('text-center')
        .text(langSet[selectedLanguageCode][key]);
    }
  }


  await $('#modal_body').empty();
  divParent.appendTo('#modal_body');

  $('#modal_footer_normal').hide();
  $('#modal_footer_product').show();

  modal.show();
  $('#confirm_button_product').off('click');
  $('#confirm_button_product').on('click', async () => await callback());
}


async function showConfirmCheck (callback) {


  const modal = new bootstrap.Modal(document.getElementById('modal'));

  $('#modal-size') //
    .addClass('modal-lg'); //

  $('#modal_header_text')
    .addClass('text-center') //
    .text('お会計でよろしいでしょうか？'); //

  $('#modal_body').empty();
  if (placePrice == 0
    && chargePrice == 0
  ) {
    $('<div>')
     .addClass('text-center h4') //
      .text('お支払い料金はありません。このままご退店できます。') //
      .appendTo('#modal_body');

  } else { 
    $('<div>')
     .addClass('text-center h4') //
      .text('座席票を持ってレジにて会計をお願いします。') //
      .appendTo('#modal_body');
  }

    if (diffPrice > 0 ){
    $('<br>')
      .appendTo('#modal_body');
      $('<div>')
      .addClass('text-center h4 text_decoration text-red') //
      .text(diffPrice + '円のチャージ料金が発生しますが、よろしいでしょうか？') //
      .appendTo('#modal_body');

    $('<div>')
      .addClass('text-center h4') //
      .text('チャージ料金は注文することで 0円 にすることができます。') //
      .appendTo('#modal_body');
    $('<div>')
      .addClass('text-center h4') //
      .text('また商品の持ち帰りも可能です。') //
      .appendTo('#modal_body');
    }
    
  $('#modal_footer_normal').show();
  $('#modal_footer_product').hide();

  modal.show();

  $('#confirm_button').off('click');
  $('#confirm_button').on('click', async () => await callback());
}

async function postProduct() {
  const placeName = $('#place_name').attr('name');

  let countParam = Number($('#count').text());
  if (countParam == 0) {
    return;
  }

  try {
    const res = await fetch(
      `${window.location.origin}/v1/place/${placeName}/order?device_key=${deviceKey}`,
      {
        method: 'POST',
        body: JSON.stringify({
          jan_code: selectedJANCode,
          count: countParam,
        }),
        headers: {
          'Content-type': 'application/json; charset=UTF-8',
        },
      }
    );
    if (!res.ok) {
      throw new Error(res.status);
    }
  } catch (error) {
    showErrorModal(error);
    console.error('Error :', error);
  }
}

async function closeOrder(place) {
  const placeName = $('#place_name').attr('name');
  let placeStatus = placeStatusBilling;
  
  if (placePrice == 0 && chargePrice == 0
  ) {
    placeStatus = placeStatusCleaning;
  }

  try {
    // 座席のステータスを変更
    const resStatus = await fetch(
      `${window.location.origin}/v1/place/${placeName}?device_key=${deviceKey}`,
      {
        method: 'PATCH',
        body: JSON.stringify({
          place_status: placeStatus,
        }),
        headers: {
          'Content-type': 'application/json; charset=UTF-8',
        },
      }
    );
    if (!resStatus.ok) {
      throw new Error(resStatus.status);
    }
  } catch (error) {
    console.error('Error :', error);
  }

    // bodyString = await createReceiptContent(place, diffPrice);
    
    // 取引明細レシートを送信
    // const res = await fetch(`${window.location.origin}/v1/place/${placeName}`, {
    //   method: 'PATCH',
    //   body: JSON.stringify({
    //     place_status: placeStatusBilling,
    //   }),
    //   headers: {
    //     'Content-type': 'application/json; charset=UTF-8',
    //   },
    // });
    // if (!res.ok) {
    //   throw new Error(res.status);
    // }


}

function attachQRCode() {

  let locationHost = location.host
  // if (locationHost == 'localhost:6789'){
  //   locationHost = '192.168.10.10:6789'
  // }

  const baseHost = location.protocol + '//' + locationHost ; 

  
  const placeName = $('#place_name').attr('name');

  basePath = '/v1/monitor/internet/' + placeName; //
  queryParam = '?permission=' + $('#permission').attr('name');

  $('.image_qr')
    .html('')
    .qrcode({
      width: 500,
      height: 500,
      text: baseHost + basePath + queryParam,
    });
}

function showQR() {
  $('#normal_billing').hide();
  $('#extra_billing').show();
}



// ページロード時の初期化
