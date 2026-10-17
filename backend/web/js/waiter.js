const deviceKey = new URLSearchParams(window.location.search).get('device_key'); // device_keyを取得

document.addEventListener("DOMContentLoaded", function () {
  const allTableList = JSON.parse(document.getElementById("all_table_list").innerText);

  buildReceptionMonitor(allTableList)
});

async function buildReceptionMonitor(placeList) {

  groupList = [];
  nameList = [];

  console.log("Place List:", placeList); // 追加: placeListの内容を表示

  let lengthCol = Object.keys(placeList).length;

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

  await $.each(placeList, function (key) {
    $('<div>')
      .appendTo(divBody)
      .attr('id', 'place_group-' + key)
      .addClass('text-center alert m-0 col py-0');
  });

  $('#body_div').empty();
  $('#body_div').append(divBody);

  
  $.each(placeList, function (key, places) {
    // ソート処理を追加
    places.sort();  // 文字列として昇順ソート

    for (let place of places) {
      let alertColor = 'alert-success border-white '; 
      let icon = '';
      let iconColor = '';

      let placeDiv = $('<div>')
        .appendTo('#place_group-' + key)
        .attr('id', place)
        .addClass('rounded-pill text-center border-1 h3 alert col ' + alertColor)
        .on('click', function () {
          // クリックしたplaceをローカルストレージに保存
          localStorage.setItem('place', place);
          window.location.href = `/v1/monitor/waiter/place/${place}?device_key=${deviceKey}`;
        })
        .text(place + ' ');

      $('<div>')
        .appendTo(placeDiv)
        .addClass('fas ' + icon)
        .addClass(iconColor);
    }
  });

  $('#head_div').empty();
  $('#head_div').append(divHead);
}

// localStorageからplaceを取得し、遷移する
const place = localStorage.getItem('place');

if (place) {
    window.location.href = `/v1/monitor/waiter/place/${place}?device_key=${deviceKey}`;
}
