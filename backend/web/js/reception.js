
let isCognitionAlcoholMinor = false;
let isCognitionAlcoholCar = false;
let isCognitionChargeCost = false;
let isCognitionRisqueDamage = false;
let isCognitionProblemsBetweenCustomers = false;
let deviceKey;
let isConfirmed = false;

let isInternal = "false" ;
let is_enable_charge = "false" ;
let nowSituate = "" ;

// 初期化関数
async function init() {
  isInternal = getParam('is_internal') === "true"; // true であれば厨房用
  deviceKey = getParam('device_key');
 is_enable_charge = $('#is_enable_charge').attr('name') === "true";
  if (isInternal) { 
    $('#border').hide();
    $('#icon_explain').hide();
    $('#explain_div').hide();
    $('#store_explain').hide();
    $('#icon_status').addClass('col-10 ');
  }


}

async function patchPlace(elem) {
  const id = $(elem).attr('id');
  const lineNumber = $(elem).attr('line_number');
  const customerID = $(elem).attr('customer_id');
  const storeNameOriginal = $('#store_name').attr('name');
  const storeAddress = $('#store_address').attr('name');
  isConfirmed = true //

  try {
    await fetch(
      `${window.location.origin}/v1/place/${id}?device_key=${deviceKey}`,
      {
        method: 'PATCH',
        body: JSON.stringify({
          place_status: placeStatusUsing,
        }),
        headers: {
          'Content-type': 'application/json; charset=UTF-8',
        },
      }
    );

    const date = new Date();
    const Y = date.getFullYear();
    const M = ('00' + (date.getMonth() + 1)).slice(-2);
    const D = ('00' + date.getDate()).slice(-2);

    let newLine = '\n';
    let horizontalLine = '-';
    let storeName = '^^^' + storeNameOriginal;
    let yyyymmdd = Y + '/' + M + '/' + D;
    let placeLineNumber = id + '-' + lineNumber;
    let attention = 'こちらの紙を会計時にご使用願います';
    let qrString = '{code:' + customerID + '; option:qrcode,8,H}';

    let bodyString =
      storeName +
      newLine +
      horizontalLine +
      newLine +
      storeAddress +
      newLine +
      horizontalLine +
      newLine +
      yyyymmdd +
      '-' +
      placeLineNumber +
      newLine +
      horizontalLine +
      newLine +
      newLine +
      '^^^' +
      id +
      newLine +
      newLine +
      qrString +
      newLine +
      newLine +
      horizontalLine +
      newLine +
      newLine +
      attention;



    // QR 印刷
    await fetch(`https://localhost:5495/printer1`, {
      method: 'POST',
      body: bodyString,
      mode: 'no-cors',
    });
  } catch (error) {
    console.error('Error :', error);
  }
}

async function buildOnSse(eventData) {
  if (typeof eventData === 'string') {
    place = JSON.parse(eventData);
  } else {
    place = eventData; // Assuming it's already an object
  }

  if (place != nowSituate) {
    nowSituate = place;
    const audio = new Audio('../../../sound/chime.mp3');
    if (isConfirmed) {
      audio.playbackRate = 0.7;
      audio.play();
    }
  }

  await buildReceptionMonitor(place);
}

async function buildReceptionMonitor(placeList) {
  groupList = [];
  capacityList = [];
  nameList = [];
  printDebug('placeList', placeList); // -----------------------------

  await $.each(placeList, function (_, place) {
    capacityList.push(place.capacity);
    groupList.push(place.place_group);
    nameList.push(place.capacity + ':' + place.place_group);
  });

  let nameSort = Array.from(new Set(nameList)).sort();
  let nameUnique = new Set(nameSort);

  let capacityMapCount = new Map();
  let lengthCol = 0;

  await $.each(nameSort, function (_, name) {
    const words = name.split(':');
    const capacity = words[0];
    capacityMapCount.set(capacity, (capacityMapCount.get(capacity) || 0) + 1);
    lengthCol++;
  });

  const divHead = $('<div>').addClass('row mx-0');
  const divBody = $('<div>').addClass('row');

  let colSize = 0;
  let sideColSize = 0;

  if (12 % lengthCol == 0) {
    colSize = 12 / lengthCol;
  } else if (lengthCol == 5) {
    colSize = 2;
    sideColSize = 1;
  }

  if (sideColSize != 0) {
    $('<div>')
      .appendTo(divHead)
      .addClass('text-center alert alert-light col-' + sideColSize);
  }

  $('#explain_charge').empty();
  for await (let key of capacityMapCount) {
    let capacity = key[0];
    let quantity = key[1];
    let colX = 'col-' + colSize * quantity;

    // 一番上 人数の行
    $('<div>')
      .appendTo(divHead)
      // .attr('id', 'capacity-' + capacity)
      .addClass('text-center border-white h3 alert alert-info ' + colX)
      .text(+capacity + '人席');

    let chargeBasePrice = $('#charge_base_price').attr('name');
    
    const contentText = '・'+ capacity + '人席' + '：' +capacity*chargeBasePrice  + '円'
    let candidate = $('<h6>')
      .text(contentText);

    $('#explain_charge').append(candidate);
  }

  $('#body_div').empty();
  $('#body_div').append(divBody);

  for await (let key of nameUnique) {
    const word = key.split(':');
    // let capacity = word[0];
    let place = word[1];

    $('<div>')
      .appendTo(divBody)
      .attr('id', 'place_group-' + place)
      .addClass('text-center alert m-0 col py-0');
  }

  await $.each(placeList, function (_, place) {
    let alertColor = 'alert-secondary border-white ';
    let icon = '';
    let iconColor = '';

    switch (place.place_status) {
      case placeStatusCleaning:
        icon = ' fa-broom ';
        if (isInternal ) {
          iconColor = ' icon-red ';
        }
        break;

      case placeStatusOpening:
        icon = ' fa-lock-open ';
        alertColor = 'alert-light border-info ';
        break;

      case placeStatusUsing:
        icon = ' fa-lock ';
        break;

      case placeStatusBilling:
        icon = ' fa-money-bill ';
        if (isInternal) {
          iconColor = ' icon-red ';
        }
        break;

      case placeStatusStopping:
        icon = ' fa-ban ';
        break;

      default:
        console.error('error status : ', order);
        break;
    }

    let placeDiv = $('<div>')
      .appendTo('#place_group-' + place.place_group)
      .attr('id', place.place_name)
      .attr('line_number', place.line_number)
      .attr('status', place.place_status)
      .attr('permission', place.permission)
      .attr('customer_id', place.customer_id)
      .addClass('rounded-pill text-center border-1 h3 alert col ' + alertColor)
      .on('click', function () {
        showConfirmModal(
          async () => {
            await patchPlace(this);
          }, //
          place.place_status, //
          place.capacity, //
          place.place_name //
        );
      })
      .text(place.place_name + ' ');

    $('<div>')
      .appendTo(placeDiv)
      .addClass('fas' + icon)
      .addClass(iconColor);
  });

  $('#head_div').empty();
  $('#head_div').append(divHead);
  
}

async function showConfirmModal (
  callback,
  place_status,
  place_capacity,
  place_name
) {
  const modal = new bootstrap.Modal(document.getElementById('modal'));
  if (place_status != placeStatusOpening) return;

  $('#modal-size') //
    .addClass('modal-lg'); //

  $('#modal_header_text') //
    .addClass('text-left') //
    .text(
      place_capacity + '人席:' + place_name + ' をご利用でよろしいですか？'
    ); //

  $('#modal_body').empty();
  if (place_status == placeStatusOpening) {
    $('#modal_footer_normal').show();
  } else { 
    $('#modal_footer_normal').hide();
  }

  $('<br>')
    .appendTo('#modal_body')

  // $('<div>')
  //   .appendTo('#modal_body')
  //   .addClass('text-left h3') //
  //   .text('ご利用には以下の同意をお願いいたします。') //
  // ;

  // isCognitionAlcoholMinor = false;
  // isCognitionAlcoholCar = false;
  isCognitionChargeCost = false;
  isCognitionRisqueDamage = false;
  isCognitionProblemsBetweenCustomers = false;
  isEnvLocal = getParam('env') == 'local';

  if (isEnvLocal) { 
    $('#modal_footer_normal').show();
    $('#notice_number_table').show();
  }


  // let checkbox1 = await createAgree(
  //   'isCognitionAlcoholMinor',
  //   ' 20歳未満の飲酒は法律で罰せられることを知っています。'
  // );
  // checkbox1.appendTo('#modal_body');

  // let checkbox2 = await createAgree(
  //   'isCognitionAlcoholCar',
  //   ' 飲酒運転は法律で罰せられることを知っています。'
  // );
  // checkbox2.appendTo('#modal_body');


  // let checkbox3 = await createAgree(
  //   'isCognitionChargeCost',
  //   ' 一定時間でチャージ料金が発生することを同意しました。'
  // );
  // checkbox3.appendTo('#modal_body');


  // let checkbox4 = await createAgree(
  //   'isCognitionRisqueDamage',
  //   ' 正規な利用以外による器物損害について弁償請求する可能性があることを同意しました。'
  // );
  // checkbox4.appendTo('#modal_body');


  // let checkbox5 = await createAgree(
  //   'isCognitionProblemsBetweenCustomers',
  //   ' お客様同士のトラブルは当事者で解決することを同意しました。'
  // );
  // checkbox5.appendTo('#modal_body');

    $('<br>')
    .appendTo('#modal_body')

  $('<div>')
    .appendTo('#modal_body')
    .addClass('text-left h4 ') //
    .attr('id', 'notice_number_table')
    // .hide()
    .text('座席票が発行されますので、お会計時まで保管をお願いいたします。') //
  ;

  $('<br>')
    .appendTo('#modal_body');

  modal.show();

  $('#confirm_button').off('click');
  $('#confirm_button').on('click', async () => await callback());
}



async function changeAgree( elem ) {
  targetID = $(elem).attr('targetAgree') 

  const uncheck = 'fa-square'
  const checked = 'fa-check-square'
  
  $('#' + targetID)
    .removeClass(uncheck)
    .removeClass(checked);

  let className = uncheck
  switch (targetID) {
    // case 'isCognitionAlcoholMinor':
    //   isCognitionAlcoholMinor = !isCognitionAlcoholMinor
    //   if (isCognitionAlcoholMinor){
    //     className = checked
    //   }
    // break;

    // case 'isCognitionAlcoholCar':
    //   isCognitionAlcoholCar = !isCognitionAlcoholCar
    //   if (isCognitionAlcoholCar){
    //     className = checked
    //   }
    // break;

    case 'isCognitionChargeCost':
      isCognitionChargeCost = !isCognitionChargeCost
      if (isCognitionChargeCost){
        className = checked
      }
    break;

    case 'isCognitionRisqueDamage':
      isCognitionRisqueDamage = !isCognitionRisqueDamage
      if (isCognitionRisqueDamage){
        className = checked
      }
    break;

    case 'isCognitionProblemsBetweenCustomers':
      isCognitionProblemsBetweenCustomers = !isCognitionProblemsBetweenCustomers
      if (isCognitionProblemsBetweenCustomers){
        className = checked
      }
    break;

    
    default:
      console.error('error status : ', order);
    break;
  }


  $('#' + targetID)
    .addClass(className)

  // isAgree =
  //   // isCognitionAlcoholMinor &&
  //   // isCognitionAlcoholCar &&
  //   isCognitionChargeCost &&
  //   isCognitionRisqueDamage &&
  //   isCognitionProblemsBetweenCustomers 

  // $('#modal_footer_normal').hide();
  // $('#notice_number_table').hide();
isAgree =true
  if (isAgree) {
    $('#modal_footer_normal').show();
    $('#notice_number_table').show();
  }

}
  

async function createAgree( idName, agreeText ) {
  let checkbox = $('<div>')
  .addClass('text-left ms-3 ');

  $('<span>')
  .appendTo(checkbox)
  .addClass('h4 far fa-square')
  .attr('id', idName)
  .attr('targetAgree', idName)
  .on('click', function () {
    changeAgree(this);
  });


  $('<span>')
  .appendTo(checkbox)
  .addClass(' h4')
  .attr('targetAgree', idName)
  .text(agreeText)
  .on('click', function () {
    changeAgree(this);
  });

  return checkbox
}


// ページロード時の初期化
