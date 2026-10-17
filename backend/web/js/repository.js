let isOldState = false;
let isNowState = false;
let isReverse = false;
let isConfirmed = false;
let miniProductMap;
  
  
async function repositoryInit() {
  await setLocalStorage('isNoticeable', false);
  isReverse = getParam('is_reverse');
  miniProductMap = JSON.parse($('#mini_product_map').text());
  await createTableDesign();
}

async function createTableDesign() {
  const allTableList = JSON.parse($('#all_table_list').text());
  const bodyDiv = $('#body_div');

  let rowSize = 0;
  let colSize = 0;
  let sideColSize = 0;
  let centerColSize = 0;

  // 水平方向 最大値計算 ---------------------------------
  $.each(allTableList, function (_, row) {
    if (rowSize < row.length) {
      rowSize = row.length;
    }
  });

  // 垂直方向計算 ---------------------------------
  let colList = Object.keys(allTableList);
  if (isReverse) {
    colList.reverse(); // Reverse the colList
  }

  colSize = colList.length;

  if (12 % colSize == 0) {
    colSize = 12 / colSize;
  } else if (colSize == 5) {
    colSize = 5;
    sideColSize = 1;
  } else if (12 % colSize == 1) {
    centerColSize = 1;
  }

  for (let row = 1; row <= rowSize; row++) {
    let number = row;
    if (isReverse) {
      number = rowSize - row + 1; 
    } 

    let rowDiv = $('<div>').appendTo(bodyDiv).addClass('col-12 row');

    $('<div>').appendTo(rowDiv).addClass('divBodySide');

    $.each(colList, function (_, col) {
      let tableDiv = $('<div>')
        .appendTo(rowDiv)
        .addClass('alert alert-info mb-3 col-' + colSize); // ---------------

      let tableHead = $('<div>')
        .appendTo(tableDiv)
        .addClass('text-center font-40px col-12 row ');

      $('<div>')
        .appendTo(tableHead)
        .addClass('text-center col-6')
        .text(col + number);

      $('<div>')
        .appendTo(tableHead)
        .addClass('text-center col-6')
        .append($('<i>').addClass('fas fa-lock-open'))
        .attr('id', col + number + '_status');

      // 注文表示領域
      $('<div>')
        .appendTo(tableDiv)
        .addClass('row')
        .attr('id', col + number);

      });

    $('<div>').appendTo(rowDiv).addClass('divBodySide');
  } // ========================

  
}


async function buildOnSse(eventData) {
  const place = typeof eventData === 'string' ? JSON.parse(eventData) : eventData;
  await buildRepositoryMonitor(place);
}

async function buildRepositoryMonitor(placeList) {
  let isNoticeable = false;
  printDebug('placeList', placeList); // -----------------------------

  const statusMap = {
    0: 'fa-broom',
    1: 'fa-lock-open',
    2: 'fa-lock',
    3: 'fa-money-bill',
    4: 'fa-ban',
  };

  
  await $.each(placeList, function (_, place) {
    $('#' + place.place_name + '_status').empty();
    $('#' + place.place_name).empty();
    let addColor = '';
    if (place.place_status == 0 || place.place_status == 3) {
      addColor = "icon-red";
    }

    $('#' + place.place_name + '_status').append(
      $('<i>').addClass(addColor + ' fas ' + statusMap[place.place_status] )
    );

    const placeJANCountMap = new Map();

    $.each(place.order_list, function (_, order) {
      if (order.order_status === orderStatusPreparing) {
        isNoticeable = true;
        const currentCount = placeJANCountMap.get(order.jan_code) || 0;
        placeJANCountMap.set(order.jan_code, currentCount + order.count);
      }
    });
    placeJANCountMap.forEach(function (count, janCode) {
      const divParent = $('#' + place.place_name);

      const divProduct = $('<div>')
        .appendTo(divParent)
        .addClass(
          'rounded-3 row border border-info border-1 margin-0 col-6 height-min-content background-white'
        );
      // 商品場所
      $('<div>')
        .appendTo(divProduct)
        .addClass('col-12 text-center badge bg-white text-dark font-40px')
        .text(miniProductMap[janCode].place);

      // 商品画像
      const divImage = $('<div>')
        .appendTo(divProduct)
        .addClass('col-7 background-white');

      $('<img>')
        .appendTo(divImage)
        .addClass('img-fluid rounded-start')
        .attr(
          'src',
          '/image/product.svg?item=' +
            janCode +
            '.jpg?version=' + version
        );

      
      // 商品カウント
      $('<div>')
        .appendTo(divProduct)
        .addClass('col-5 text-center badge bg-white text-dark font-70px')
        .text(count);

    });

  });

  const isNoticeableOld = JSON.parse(await getLocalStorage('isNoticeable')); // Convert to boolean
  if (isNoticeableOld == false && isNoticeable == true){
    const audio = new Audio('../../../sound/chime.mp3');
    if (isConfirmed) {
      audio.playbackRate = 0.7;
      audio.play();
    }
  } 
  await setLocalStorage('isNoticeable', isNoticeable);


  
}

function confirm () {
  isConfirmed = true //
  $('#confirm_screen').hide();
  $('#main_screen').show();
} //
