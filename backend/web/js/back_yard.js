// グローバル変数
let placeStatus = placeStatusCleaning;
let selectedOrderNumber = 0;
let miniProductMap;
let deviceKey;

async function init () {
  miniProductMap = JSON.parse($('#mini_product_map').text());

}
async function buildOnSse(eventData) {
  await new Promise((resolve) => {
    const intervalId = setInterval(() => {
      clearInterval(intervalId);
      buildBackYardMonitor(eventData);
      resolve();
    }, 100);
  });
}

async function buildBackYardMonitor(eventData) {
  if (typeof eventData === 'string') {
    place = JSON.parse(eventData);
  } else {
    place = eventData; 
  }

  deviceKey = getParam('device_key');

  placeStatus = place.place_status;
  printDebug('place', place); // -----------------------------

  // ステータス変更
  $('.page-title').removeClass('border-danger border-2');
  $('#' + placeStatus).addClass('border-danger border-2');

  const placeName = $('#place_name').attr('name');
  $('#line_number').text(placeName);

  if (placeStatus == placeStatusCleaning) {
    await $('#order_history').empty();
  }

  if (placeStatus == placeStatusUsing || placeStatus == placeStatusBilling) {
    await createHistoryScreen(place);
    $('#line_number').text(placeName + '-' + place.line_number);
  }
}

async function createHistoryScreen(place) {
  let reverseOrderList = [];

  const divOrder = $('<div>').addClass('row');

  await $.each(place.order_list, function (_, order) {
    reverseOrderList.push(order);
  });
  // reverseOrderList.sort().reverse();

  const divPreparing = $('<div>').appendTo(divOrder).addClass('row');
  $('<hr>').appendTo(divOrder);
  const divPassed = $('<div>').appendTo(divOrder).addClass('row');

  reverseOrderList.forEach((order) => {
    const divCard = $('<div>').addClass('col-6 alert mb-3 ');
    const divRow = $('<div>')
      .appendTo(divCard)
      .addClass('row rounded-5 alert alert-light background-white ');

    // 商品場所
    $('<div>')
      .appendTo(divRow)
      .addClass('col-12 text-center badge bg-white text-dark font-40px')
      .text(miniProductMap[order.jan_code].place);

    // 商品画像
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


    // 商品カウント
    $('<div>')
      .appendTo(divRow)
      .addClass('col-md-6 text-center text-dark font-150px')
      .text(order.count);

    let orderStatus = orderStatusCancel
    let textOrderStatus = '';
    let colorOrderStatus = '';

    switch (order.order_status) {
      case orderStatusCancel:
        orderStatus = orderStatusCancel;
        textOrderStatus = textOrderStatusCancel;
        colorOrderStatus = colorOrderStatusCancel ;
        $('<div>')
          .appendTo(divCard)
          .addClass(
            'border-0 badge align-items-center font-20px position-absolute top-0 start-0 ' +
              colorOrderStatus
          )
          .text(textOrderStatus);
        divCard.appendTo(divPassed);
        divCard.on('click', async () => {
          await showConfirmModal(order);
        });
        break;

      case orderStatusPreparing:
        orderStatus = orderStatusPreparing;
        textOrderStatus = textOrderStatusPreparing;
        colorOrderStatus = colorOrderStatusPreparing;

        $('<div>')
          .appendTo(divCard)
          .addClass(
            'border-0 badge align-items-center font-20px position-absolute top-0 start-0 ' +
              colorOrderStatus
          )
          .text(textOrderStatus);
        divCard.appendTo(divPreparing);
        divCard.on('click', async () => {
          await showConfirmModal(order);
        });

        break;
      case orderStatusPassed:
        orderStatus = orderStatusPassed;
        textOrderStatus = textOrderStatusPassed;
        colorOrderStatus = colorOrderStatusPassed;
        $('<div>')
          .appendTo(divCard)
          .addClass(
            'border-0 badge align-items-center font-20px position-absolute top-0 start-0 ' +
              colorOrderStatus
          )
          .text(textOrderStatus);
        divCard.appendTo(divPassed);
        divCard.on('click', async () => {
          await showConfirmModal(order);
        });
        break;
      default:
        console.error('error status : ', order);
        break;
    }
  });

  await $('#order_history').empty();
  divOrder.appendTo('#order_history');
}

async function changeStatus(elem) {
  const status = Number(elem.id);
  const name = $('#place_name').attr('name');

  // 同じステータスの場合はなにもしない
  if (placeStatus == status) return;

  // 本番稼働ではbackyardのステータス変更は制限する。
  // placeStatusUsing/placeStatusBilling には変更不可
  if (!isEnvLocal) {
    if (status == placeStatusUsing || status == placeStatusBilling) return;

    // placeStatusBilling → placeStatusCleaningはシステム上でしか許容しない
    // if (placeStatus == placeStatusBilling && status == placeStatusCleaning) return;
  }

  const confirmChange = confirm('ステータスを変更しますか?');
  if (!confirmChange) return; // ユーザーがキャンセルした場合、処理を中止


  try {
    await fetch(
      `${window.location.origin}/v1/place/${name}?device_key=${deviceKey}`,
      {
        method: 'PATCH',
        body: JSON.stringify({
          place_status: status,
        }),
        headers: {
          'Content-type': 'application/json; charset=UTF-8',
        },
      }
    );
  } catch (error) {
    console.error('Error :', error);
  }
}
async function showConfirmModal(order) {
  selectedOrderNumber = order.order_number;

  let targetStatuses;
  if (order.order_status === orderStatusCancel) {
    targetStatuses = [orderStatusPreparing, orderStatusPassed];
  } else if (order.order_status === orderStatusPreparing) {
    targetStatuses = [orderStatusCancel, orderStatusPassed];
  } else {
    targetStatuses = [orderStatusCancel, orderStatusPreparing];
  }

  const statusColors = {
    [orderStatusCancel]: colorOrderStatusCancel,
    [orderStatusPreparing]: colorOrderStatusPreparing,
    [orderStatusPassed]: colorOrderStatusPassed,
  };
  const statusTexts = {
    [orderStatusCancel]: textOrderStatusCancel,
    [orderStatusPreparing]: textOrderStatusPreparing,
    [orderStatusPassed]: textOrderStatusPassed,
  };

  const modal = new bootstrap.Modal(document.getElementById('modal'));
  $('#modal-size') //
    .addClass('modal-lg'); //

  const divParent = $('<div>');
  const divRow = $('<div>').appendTo(divParent).addClass('row');
  const divCardLeft = $('<div>')
    .appendTo(divRow)
    .addClass('col-md-6 text-center');
  $('<img>')
    .appendTo(divCardLeft)
    .addClass('img-fluid rounded-start img_size300px')
    .attr(
      'src',
      '/image/product.svg?item=' +
        order.jan_code +
        '.jpg?version=' +
        version
    );

  const divCardRight = $('<div>')
    .appendTo(divRow)
    .addClass('col-md-6 text-center');
  $('<div>')
    .appendTo(divCardRight)
    .addClass('col-md-6 text-center  text-dark font-150px')
    .text(order.count);

  $('#left_button').remove();
  $('#cancel_button').hide();

  // 左側のボタン
  $('<button>')
    .prependTo('#modal_footer_normal')
    .addClass('btn text-white btn-lg py-3 col-6 mx-0 ' + statusColors[targetStatuses[0]])
    .attr('id','left_button')
    .prop('type','button')
    .text(statusTexts[targetStatuses[0]])
    .on('click', async () => await updateOrderStatus(modal, targetStatuses[0]));

  $('#confirm_button').removeClass(
      ' ' + colorOrderStatusCancel ,
    )

  $('#confirm_button').removeClass(
      ' ' + colorOrderStatusPreparing ,
    )

  $('#confirm_button').removeClass(
      ' ' + colorOrderStatusPassed ,
    )

  // 右側のボタン
  $('#confirm_button')
    .removeAttr('data-bs-dismiss')
    .prop('disabled', false)
    .addClass(' col-6 mx-0 ' + statusColors[targetStatuses[1]])
    .text(statusTexts[targetStatuses[1]]);
  await $('#modal_body').empty();
  divParent.appendTo('#modal_body');

  $('#modal_footer_normal').show();
  $('#modal_footer_product').hide();
  modal.show();


  $('#confirm_button').off('click');
  $('#confirm_button').on('click', async () => await updateOrderStatus(modal, targetStatuses[1]));
}

async function updateOrderStatus(modal, orderStatus) {
  $('#left_button, #confirm_button').prop('disabled', true);
  const updated = await patchOrder(orderStatus);
  if (!updated) {
    $('#left_button, #confirm_button').prop('disabled', false);
    return;
  }

  if (document.activeElement instanceof HTMLElement) {
    HTMLElement.prototype.blur.call(document.activeElement);
  }
  modal.hide();
}

async function patchOrder(orderStatus) {
  const placeName = $('#place_name').attr('name');
  const orderNumber = selectedOrderNumber;

  try {
    const res = await fetch(
      `${window.location.origin}/v1/place/${placeName}/order/${orderNumber}?device_key=${deviceKey}`,
      {
        method: 'PATCH',
        body: JSON.stringify({
          order_status: orderStatus,
        }),
        headers: {
          'Content-type': 'application/json; charset=UTF-8',
        },
      }
    );

    if (!res.ok) {
      throw new Error(res.status);
    }
    return true;
  } catch (error) {
    console.error('Error :', error);
    return false;
  }
}
