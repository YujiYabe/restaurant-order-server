let allProducts;

const titleIsValid = '有効';
const titlePlace = '場所';
const titleJanCode = 'JANコード';
const titleRecipe = 'レシピ';
const titleCalorie = 'カロリー';

const titleName = '名前';

// const titleInformation = '短名';

const titleStock = '在庫';

const titleCostPrice = '原価';
const titleRetailPrice = '小売価格';
const titleOperationCost = '作業コスト';
const titleEnergyCost = '光熱費';

const titleAllergy = 'アレルギー';
let deviceKey;

async function adminProductInit() {
  deviceKey = getParam('device_key');
  allProducts = JSON.parse($('#all_product_list').text());
  inventoryMap = JSON.parse($('#inventory_map').text());

  await buildAdminMonitor();
  resetMonitorSetting();

}

function resetMonitorSetting () { 
  $(document).on('contextmenu', function (e) {
    return true;
  });

  // マウスカーソル非表示
  $('html').css('cursor', '');

  // テキスト選択不可
  $('body').css('-webkit-user-select', ''); /* Safari */
  $('body').css('-moz-user-select', ''); /* Firefox */
  $('body').css('-ms-user-select', ''); /* IE 10+ */
  $('body').css('user-select', '');

  // ピンチアウト禁止
  document.addEventListener('touchmove', mobileNoScroll, {
    passive: true,
  });
}



async function buildAdminMonitor() {
  try {
    $.each(allProducts, function (_, product) {
      const janCode = product.jan_code;
      printDebug('janCode', janCode); // -----------------------------

      let trRow = $('<tr>').appendTo('#admin_product');
      let content = $('<td>').appendTo(trRow);
      let divRow = $('<div>')
        .appendTo(content)
        .addClass('alert-info col-12 alert border-0 row ');
      let divImg = $('<div>')
        .appendTo(divRow)
        .addClass('col-2 alert alert-light background-white border-0');
      $('<img>')
        .appendTo(divImg)
        .addClass('img-thumbnail')
        .attr(
          'src',
          '/image/product.svg?item=' +
            janCode +
            '.jpg?version=' + version
        );

      let divContent = $('<div>')
        .appendTo(divRow)
        .addClass('col-10 alert alert-light background-white border-0 row ')
        .attr('name', janCode);
      //==================
      // 基本情報
      let divBaseInfo = $('<div>').addClass('col-12 row').appendTo(divContent);

      // -----------------
      // 有効  場所  JANコード レシピ  カロリー
      let divIsValid = $('<div>').addClass('col-2').appendTo(divBaseInfo);
      let divPlace = $('<div>').addClass('col-3').appendTo(divBaseInfo);
      let divJanCode = $('<div>').addClass('col-3').appendTo(divBaseInfo);
      let divRecipe = $('<div>').addClass('col-2').appendTo(divBaseInfo);
      let divCalorie = $('<div>').addClass('col-2').appendTo(divBaseInfo);
      // -----------------
      // 名前
      let divStock = $('<div>').addClass('col-3').appendTo(divBaseInfo);
      let divName = $('<div>').addClass('col-9').appendTo(divBaseInfo);

      // -----------------
      // 商品情報
      // let divInfo = $('<div>').addClass('col-12').appendTo(divBaseInfo);

      //==================
      // 小売価格   原価  作業コスト  光熱費
      let divCost = $('<div>').addClass('col-12 row').appendTo(divContent);

      let divCostPrice = $('<div>').addClass('col-3').appendTo(divCost);
      let divRetailPrice = $('<div>').addClass('col-3').appendTo(divCost);
      let divOperationCost = $('<div>').addClass('col-3').appendTo(divCost);
      let divEnergyCost = $('<div>').addClass('col-3').appendTo(divCost);

      //==================
      // -----------------
      let divOther = $('<div>').addClass('col-12 row').appendTo(divContent);
      let divAllergy = $('<div>').addClass('col-12').appendTo(divOther);

      makeCommonDiv(titleName, product.name_ja, divName);
      makeCommonDiv(titleJanCode, product.jan_code, divJanCode);
      makeCommonDiv(titleRecipe, product.recipe, divRecipe);
      makeCommonDiv(titleCalorie, product.calorie, divCalorie);


      {
        // 有効
        $('<span>')
          .appendTo(divIsValid)
          .text(titleIsValid)
          .addClass('col badge bg-success');

        $('<span>').appendTo(divIsValid).text(' ');

        let divForm = $('<span>').appendTo(divIsValid).addClass('checkbox');
        let checkbox = $('<input>').appendTo(divForm).prop({
          type: 'checkbox',
          checked: product.is_valid,
        });
        checkbox.on('change', function () {
          patchAdminProduct(janCode);
        });

        $('<span>').appendTo(divIsValid).text(' ');

        $('<span>')
          .appendTo(divIsValid)
          .text(product.is_valid.toString())
          .attr('id', 'valid_' + janCode);
      }

      {
        // 場所
        $('<span>')
          .appendTo(divPlace)
          .addClass('col badge bg-success')
          .text(titlePlace);

        $('<span>').appendTo(divPlace).text(' ');

        let textBox = $('<span>')
          .appendTo(divPlace)
          .addClass('text-left col-9');
        $('<input>')
          .appendTo(textBox)
          // .addClass("form-control")
          .attr('name', 'place')
          .val(product.place)
          .prop({
            type: 'text',
          });

        textBox.on('change', function () {
          patchAdminProduct(janCode);
        });
      }

      {
        // 在庫
        $('<span>')
          .appendTo(divStock)
          .addClass('col badge bg-success')
          .text(titleStock);

        $('<span>').appendTo(divStock).text(' ');

        let textBox = $('<span>')
          .appendTo(divStock)
          .addClass('text-left col-9');
        $('<input>')
          .appendTo(textBox)
          .addClass('size-stock')
          .attr('name', 'stock')
          .val(inventoryMap[janCode])
          .prop({
            type: 'number',
          });

        textBox.on('change', function () {
          patchAdminProduct(janCode);
        });
      }

      {
        // 原価 CostPrice
        $('<span>')
          .appendTo(divCostPrice)
          .addClass('col badge bg-success')
          .text(titleCostPrice);

        $('<span>').appendTo(divCostPrice).text(' ');

        let textBox = $('<span>')
          .appendTo(divCostPrice)
          .addClass('text-left col-9');
        $('<input>')
          .appendTo(textBox)
          .addClass('size-stock')
          .attr('name', 'costPrice')
          .val(product.cost_price)
          .prop({
            type: 'number',
          });

        textBox.on('change', function () {
          patchAdminProduct(janCode);
        });
      }

      {
        // 小売価格 RetailPrice
        $('<span>')
          .appendTo(divRetailPrice)
          .addClass('col badge bg-success')
          .text(titleRetailPrice);

        $('<span>').appendTo(divRetailPrice).text(' ');

        let textBox = $('<span>')
          .appendTo(divRetailPrice)
          .addClass('text-left col-9');
        $('<input>')
          .appendTo(textBox)
          .addClass('size-stock')
          .attr('name', 'retailPrice')
          .val(product.retail_price)
          .prop({
            type: 'number',
          });

        textBox.on('change', function () {
          patchAdminProduct(janCode);
        });
      }
    
      {
        // 作業コスト OperationCost
        $('<span>')
          .appendTo(divOperationCost)
          .addClass('col badge bg-success')
          .text(titleOperationCost);

        $('<span>').appendTo(divOperationCost).text(' ');

        let textBox = $('<span>')
          .appendTo(divOperationCost)
          .addClass('text-left col-9');
        $('<input>')
          .appendTo(textBox)
          .addClass('size-stock')
          .attr('name', 'operationCost')
          .val(product.operation_cost)
          .prop({
            type: 'number',
          });

        textBox.on('change', function () {
          patchAdminProduct(janCode);
        });
      }

      {
        // 作業コスト energyCost
        $('<span>')
          .appendTo(divEnergyCost)
          .addClass('col badge bg-success')
          .text(titleEnergyCost);

        $('<span>').appendTo(divEnergyCost).text(' ');

        let textBox = $('<span>')
          .appendTo(divEnergyCost)
          .addClass('text-left col-9');
        $('<input>')
          .appendTo(textBox)
          .addClass('size-stock')
          .attr('name', 'energyCost')
          .val(product.energy_cost)
          .prop({
            type: 'number',
          });

        textBox.on('change', function () {
          patchAdminProduct(janCode);
        });
      }

      {
        // アレルギー
        $('<span>')
          .appendTo(divAllergy)
          .addClass('col-1 badge bg-success')
          .text(titleAllergy);

        let divAllergyBody = $('<div>')
          .appendTo(divAllergy)
          .addClass('col-12 row');
        let cardTop = $('<div>')
          .appendTo(divAllergyBody)
          .addClass('col-12 row ms-0');
        const allergyList = JSON.parse(product.allergy);
        printDebug('allergyList', allergyList); // -----------------------------

        $.each(allergyList, function (key, boolean) {
          if (boolean == false) return;

          let card = $('<div>').appendTo(cardTop).addClass('card col-1 ');
          $('<img>')
            .appendTo(card)
            .addClass('img-thumbnail card-img-top')
            .attr('src', '/image/allergy.svg?name=' + encodeURIComponent(key));
          let cardBody = $('<div>').appendTo(card).addClass('text-center');
          $('<p>')
            .appendTo(cardBody)
            .addClass('text-center')
            .text(langSet[defaultLangCode][key]);
        });
      }
    });

    initializeProductFilter();
  } catch (error) {
    console.error('Error :', error);
  }
}

function initializeProductFilter() {
  const filter = document.getElementById('my_filter');
  const rows = document.querySelectorAll('#admin_product tr');
  filter.addEventListener('input', () => {
    const keyword = filter.value.trim().toLocaleLowerCase();
    rows.forEach((row) => {
      row.hidden = keyword !== '' && !row.textContent.toLocaleLowerCase().includes(keyword);
    });
  });
}

function makeCommonDiv(name, value, addTarget) {
  let div = $('<div>').appendTo(addTarget).addClass('col-12 list-group-item');
  $('<span>').appendTo(div).addClass('col badge bg-success').text(name);

  $('<span>')
    .appendTo(div)
    .addClass('col')
    .text(' ' + value);
}

async function patchAdminProduct(janCode) {
  try {
    const product = document.getElementsByName(String(janCode))[0];
    const isValid = product.querySelector('input[type="checkbox"]').checked;
    const place = product.querySelector('input[name="place"]').value;
    const stock = product.querySelector('input[name="stock"]').value;
    const costPrice = product.querySelector('input[name="costPrice"]').value;
    const retailPrice = product.querySelector('input[name="retailPrice"]').value;
    const operationCost = product.querySelector('input[name="operationCost"]').value;
    const energyCost = product.querySelector('input[name="energyCost"]').value;

    $(`span[id=valid_${janCode}]`).text(isValid.toString());
    const body = JSON.stringify(
      {
        is_valid: isValid,
        place: place,
        stock: Number(stock),
        cost_price: Number(costPrice),
        retail_price: Number(retailPrice),
        operation_cost: Number(operationCost),
        energy_cost: Number(energyCost),
      }
    )


    await fetch(
      `${window.location.origin}/v1/admin/product/${janCode}?device_key=${deviceKey}`,
      {
        method: 'PATCH',
        body: body,
        headers: {
          'Content-type': 'application/json; charset=UTF-8',
        },
      }
    );

  } catch (error) {
    console.error('Error :', error);
  }
}

// ページロード時の初期化
