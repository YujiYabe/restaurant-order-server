// グローバル変数
let video = document.getElementById('video');
let canvas = document.getElementById('canvas');



let selectedTagList = {};
let restoreSelectedTagList = {};
let restoreSelectedTagPrimary = 1;

let tagSecondaryList = [];
let selectedLanguageCode = 1041;
let selectedPriceTotal = 0;
let selectedJANCode;
let isAlreadyCreateTabletScreen = false;
let inventoryMap;
let mapJANCodePrice = new Map();
let enterTime;
let placePrice = 0;
let permission;
let isMiniProductMapOK = false;
let miniProductMap;
let chargeBaseMinute;

// 初期化関数
async function init() {
  miniProductMap = JSON.parse($('#mini_product_map').text());
  tagList = JSON.parse($('#tag_list').text());
  
  if (version != await getLocalStorage('version')) { 
    await setLocalStorage('version', version);
    resetConfig();
  }
  chargeBaseMinute = $('#charge_base_minute').attr('name') 

  permission = getParam('permission');

  $('#thanksMessage').hide();
  
  const storedSelectedTagList = await getLocalStorage('selectedTagList');
  if (storedSelectedTagList != "") {
    restoreSelectedTagList = JSON.parse(storedSelectedTagList);
  } else { 
    $.each(
      tagList, async function (_, tagPrimary) {
        let isFirst = true;
        $.each(tagPrimary.tag_secondary_list,
          async function (_, tagSecondary) {
            if (isFirst) {
              isFirst = false;
              restoreSelectedTagList[tagPrimary.id] = Number(tagSecondary.id);
              return 
            }
          }
        );
      }
    );
  }

  const storedSelectedTagPrimary = await getLocalStorage('selectedTagPrimary');
  if (storedSelectedTagPrimary != "") {
    restoreSelectedTagPrimary = storedSelectedTagPrimary;
  } else {
    restoreSelectedTagPrimary = tagList[0].id ;
  }

  await setLocalStorage('originURL', window.location.href);

  await buildProductList();


}


async function resetConfig () {
  localStorage.removeItem("selectedTagPrimary");
  localStorage.removeItem("selectedTagList");
  localStorage.removeItem("originURL");
  location.reload()
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

async function buildProductList() {
  $.each(miniProductMap, function (_, product) {

    const divParent = $('<div>') //
      .appendTo('#menu_content') //
      .addClass(' col-6 alert products p-1 mb-0') //
      .hide() //
      .attr('tag', product.tag) //
      .attr('id', product.jan_code)
      .on('click', function () {
        showConfirmOrder(async () => {
          await postProductPrepare();
        }, product); //
      });
    const rounded = $('<div>')
      .addClass(' rounded-5 alert alert-light background-white m-0') //
      .appendTo(divParent); // //

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
      .text(product.price + '円'); //
  });
  isMiniProductMapOK = true
}

async function buildWaiterMonitor(eventData) {
  if (typeof eventData === 'string') {
    tempData = JSON.parse(eventData);
  } else {
    tempData = eventData; 
  }


  const place = tempData.Place;
  inventoryMap = tempData.InventoryMap;
  printDebug('tempData', tempData); // -----------------------------

  $('#status').attr('name', place.place_status);
  $('#charge_base_capacity').attr('name', place.capacity);
  $('#customer_id').attr('name', place.customer_id);
  attachQRCode();

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
      open('about:blank', '_self').close();  
    break;

    case placeStatusOpening:
      $('#body_opening').show();
    break;

    case placeStatusBilling:
      $('#body_billing').show();
      $('#modal').hide();
    break;

    case placeStatusStopping:
      $('#body_stopping').show();
    break;

    case placeStatusUsing:
      const name = $('#place_name').attr('name');
      $('#body_using').show();
      $('#line_number').text(name + '-' + place.line_number);
      lineNumber = place.line_number;
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
      const diffBase = place.base_count;
      const chargePrice = diffBase * chargeBasePrice * chargeBaseCapacity;
      let chargeUnit = Math.floor(
        placePrice / (chargeBasePrice * chargeBaseCapacity)
      );

      // ---------------------------------
      let extendTime = new Date();
      extendTime.setTime(
        enterTime.getTime() + (chargeUnit +1) * chargeBaseMinute * 60 * 1000
      );
      let [textDateExtended, textTimeExtended] = convertTimeString(extendTime);

      let isPlaceOverCharge = placePrice > chargePrice;
      let isNowOverExtend = nowTime > extendTime;
      diffPrice = chargePrice - placePrice;

      // チャージ料金より購入金額が多い
      // 時間が延長時間未満
      if (isPlaceOverCharge && !isNowOverExtend) {
        $('#diff_price').text(textTimeExtended.slice(0, 5) + 'まで無料');
      } else {
        $('#diff_price').text(diffPrice + '円');
      }
      $('#place_price').text(placePrice + '円');



      if (isAlreadyCreateTabletScreen == false) {
        await createOrderScreen();

        if (restoreSelectedTagList != "") { 
          selectedTagList = restoreSelectedTagList
          await setLocalStorage( 'selectedTagList', JSON.stringify(restoreSelectedTagList));
        }

        if (restoreSelectedTagPrimary != 1) {
          selectedTagPrimary = restoreSelectedTagPrimary
          await setLocalStorage('selectedTagPrimary', restoreSelectedTagPrimary);
        } else { 
          selectedTagPrimary = 1
        }

        $('#tag_primary').val(selectedTagPrimary);

        $('#tag_secondary').val(selectedTagList[Number(restoreSelectedTagPrimary)]);
        await selectTagPrimary();

        filterProductByTag(selectedTagList[Number(restoreSelectedTagPrimary)] );
      }
    break;

    
    default:
      console.error('error status : ', value);
    break;
  }
}

async function createOrderScreen() {

  $.each(
    tagList,
    async function (_, tagPrimary) {

    $('<option>')
      .appendTo('#tag_primary') //
      .attr('value', Number(tagPrimary.id)) //
      .text(tagPrimary.name);
      $.each(tagPrimary.tag_secondary_list,
        async function (_, tagSecondary) {
        tagSecondaryList[tagSecondary.id] = tagSecondary.name;
      });

  });



  $('#tag_secondary').empty();
  let tagPrimaryEntity = tagList[0]; 

  $.each(tagPrimaryEntity.tag_secondary_list, function (_, tagSecondary) {
    let textName = tagSecondary.name
    let replaceName = "全て"
    if (textName.includes(replaceName)) {
      textName = replaceName
    }

    let option = $('<option>')
      .appendTo('#tag_secondary') //
      .attr('value', Number(tagSecondary.id)) //
      .text(textName);

    if (selectedTagList[tagPrimaryEntity.id] == tagSecondary.id) {
      option.prop('selected', true);
    }
  });

  isAlreadyCreateTabletScreen = true;
  // changePage('tablet');
}

async function selectTagPrimary() {

  let tagPrimaryEntity ;
  for (const key in tagList) {
    if (tagList[key].id == $('#tag_primary').val()) { 
      tagPrimaryEntity = tagList[key]
    }
  }  
  
  $('#tag_secondary').empty();


  $.each(tagPrimaryEntity.tag_secondary_list, function (_, tagSecondary) {


    let textName = tagSecondary.name
    let replaceName = "全て"
    if (textName.includes(replaceName)) {
      textName = replaceName
    }

    let option = $('<option>')
      .appendTo('#tag_secondary') //
      .attr('value', Number(tagSecondary.id)) //
      .text(textName); // 

    if (selectedTagList[tagPrimaryEntity.id] == tagSecondary.id) {
      option.prop('selected', true);
    }
  });

  setLocalStorage('selectedTagPrimary', tagPrimaryEntity.id);
  filterProductByTag(selectedTagList[Number($('#tag_primary').val())] );
}

async function selectTagSecondary () {
  await filterProductByTag($('#tag_secondary').val());

  selectedTagList[Number( $('#tag_primary').val())] = Number( $('#tag_secondary').val());
  await setLocalStorage( 'selectedTagList', JSON.stringify(selectedTagList));
}
 



async function showConfirmOrder(callback, product) {

  selectedJANCode = product.jan_code;
  selectedPriceTotal = product.price;

  const modal = new bootstrap.Modal(document.getElementById('modal'));
  $('#confirm_button_product').text('注文する');
  $('#confirm_button_product').removeClass('btn-secondary');

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

  $('<hr>')
    .appendTo(divCardBody)

  if (product.recipe != "") { 
    // レシピ
    $('<div>')
      .appendTo(divCardBody)
      .addClass('h4 text-center')
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
        .addClass('card col-3 width_18rem');
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
  // アルコールであればアラート表示;
  if (product.tag.includes('alcohol')) {
    let uuidVerifyTwenty = await getLocalStorage('uuidVerifyTwenty');
    let isIncludeList = false
    // uuidがあれば そのuuidが既に注文済みかを get処理
    if (uuidVerifyTwenty != "") {

      $('#verify_text').text("顔認証");

      const res = await fetch(
        `${window.location.origin}/v1/alcohol/include_list?uuidVerifyTwenty=${uuidVerifyTwenty}`, // url
        {
          method: 'GET',
          headers: {
            'Content-type': 'application/json; charset=UTF-8',
          },
        }
      );
      if(res.status == 200) {
        isIncludeList = true;
      }
    }


    if (!isIncludeList) {

      // 未登録
      $('.modal_footer').hide();
      $('#modal_footer_verify').show();
      const alertBody = $('<div>')
        .appendTo(divCardBody)
        .addClass('alert alert-danger')
        .attr('id', 'alert_alcohol');

      $('<div>')
        .appendTo(alertBody)
        .text(
          '当店でアルコールを購入するには 以下のご対応が必要となります。'
        );
      $('<div>')
        .appendTo(alertBody)
        .addClass(' text_decoration') //

        .text(
          '①身分証明書による年齢確認'
        );
      $('<div>')
        .appendTo(alertBody)
        .addClass('text_decoration')
        .text(
          '②本人確認'
        );
      $('<div>')
        .appendTo(alertBody)
        .text(
          'ご登録頂くと、このスマートフォンから次回以降 本人確認のみでご購入可能となります。'
        );

    } else { 

      // 登録済 
      const alertBody = $('<div>')
        .appendTo(divCardBody)
        .addClass('alert alert-danger')
        .attr('id', 'alert_alcohol');

      $('<div>')
        .appendTo(alertBody)
        .text(
          '20歳未満の飲酒は法律で禁じられています。'
        );

      $('<div>')
        .appendTo(alertBody)
        .text(
          '飲酒運転は法律で禁じられています。'
        );

      $('<br>')
        .appendTo(alertBody)

      $('<div>')
        .appendTo(alertBody)
        .addClass('text_decoration')
        .text(
          'アルコールは店内のみの提供(持ち帰り不可)となり、店員が開封してのお渡しとなります。'
        );

      $('<br>')
        .appendTo(alertBody)

      $('<div>')
        .appendTo(alertBody)
        .text('同意いただける方のみ注文ください');

      $('.modal_footer').hide();
      $('#modal_footer_product').show();
    }
  } else {
    $('.modal_footer').hide();
    $('#modal_footer_product').show();
  }

  $('#confirm_button_product').off('click');
  $('#confirm_button_product').on('click', async () => await callback());

  await $('#modal_body').empty();
  divParent.appendTo('#modal_body');

  modal.show();
}




async function jumpToVerifyPage() {
  window.open("/v1/alcohol/verify_twenty");
}


function showConfirmButtonGroup() {
  $('.modal_footer').hide();
  $('#alert_alcohol').hide();
  $('#modal_footer_product').show();
}

async function postProductPrepare () {
  // 商品がアルコールかどうかを判断。 
  // ├ アルコールでなければ そのまま注文
  // └ アルコールであれば、以下の処理を実施
  //    └uuidがあれば そのuuidが既に注文済みかを get処理
  //      ├ その日2回目以降の注文であればそのまま注文
  //      └ その日初めての注文であれば、カメラを起動

  const placeName = $('#place_name').attr('name');

  let countParam = Number($('#count').text());
  if (countParam == 0) {
    return;
  }

  let isAlcohol = false
  await $.each(miniProductMap, function (_, product) {
    if (product.jan_code == selectedJANCode) { 
      if (product.tag.includes('alcohol')) {
        isAlcohol = true
      }
    }
  });

  // アルコールでなければ、そのまま注文
  if (!isAlcohol) {
    await postProduct(
     `${window.location.origin}/v1/internet/${placeName}/order?permission=${permission}`, // url
      JSON.stringify({                                                                    // data
        jan_code: selectedJANCode,
        count: countParam,
      }),
    );

    return;
  }
  
  // アルコールであれば、以下の処理を実施
  let uuidVerifyTwenty = await getLocalStorage('uuidVerifyTwenty');

  // その日2回目以降の注文であればそのまま注文
  let data = JSON.stringify({
    jan_code: selectedJANCode,
    count: countParam,
    uuid: uuidVerifyTwenty,
  })

  await postProduct(
    `${window.location.origin}/v1/internet/${placeName}/order?permission=${permission}`, 
     data,
  );

  return;
}

async function postProduct(url, data) {
  try {
    if (url == ''){ 
      return 
    }

    const res = await fetch(
      url,
      {
        method: 'POST',
        body: data,
        headers: {
          'Content-type': 'application/json; charset=UTF-8',
        },
      }
    );
    if (!res.ok) {
    if(res.status == 401) {
        alert('注文できませんでした。');
    }
      throw new Error(res);
    }else{
      animateThanksMessage();

      localStorage.setItem('base64Image', '');

    }
  } catch (error) {
    localStorage.setItem('base64Image', '');
    console.error('Error :', error);
  }
}

async function animateThanksMessage() {
  const message = document.getElementById('thanksMessage');
  message.style.display = '';
  await message.animate([{ opacity: 0 }, { opacity: 1 }], 500).finished;
  await new Promise((resolve) => window.setTimeout(resolve, 1000));
  await message.animate(
    [
      { opacity: 1, maxHeight: `${message.scrollHeight}px` },
      { opacity: 0, maxHeight: 0 },
    ],
    1000,
  ).finished;
  message.style.display = 'none';
}

function attachQRCode() {
  const customerID = $('#customer_id').attr('name');

  $('#image_qr').empty();

  $('#image_qr')
    .html('')
    .qrcode({
      // width: 400,
      // height: 400,
      text: customerID,
    });
}





// ページロード時の初期化
