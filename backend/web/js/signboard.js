let selectedLanguageCode = 1041;
let selectedPriceTotal = 0;
let selectedJANCode;
let isAlreadyCreateTabletScreen = false;
let inventoryMap;
let mapJANCodePrice = new Map();
let placePrice = 0;
let isMiniProductMapOK = false;
let miniProductMap;


// 初期化関数
async function init() {
  miniProductMap = JSON.parse($('#mini_product_map').text());
  signboardProduct = JSON.parse($('#signboard_product').text());
  
  if (version != await getLocalStorage('version')) { 
    await setLocalStorage('version', version);
  }

  loopProduct();

  await setLocalStorage('originURL', window.location.href);
}


function loopProduct() {
  let index0 = 0;
  let index1 = 1;
  let index2 = 2;
  const maxLength = signboardProduct.length;

  setInterval(() => {

    $('#image0') //
      .attr(
        'src',
        '/image/product.svg?item=' +
          miniProductMap[signboardProduct[index0]].jan_code +
          '.jpg?version=' + version
      );
    $('#image1') //
      .attr(
        'src',
        '/image/product.svg?item=' +
          miniProductMap[signboardProduct[index1]].jan_code +
          '.jpg?version=' + version
      );
    $('#image2') //
      .attr(
        'src',
        '/image/product.svg?item=' +
          miniProductMap[signboardProduct[index2]].jan_code +
          '.jpg?version=' + version
      );

    // ---------------------------------------------
    $('#name0') //
      .addClass('h2') //
      .text(
        miniProductMap[signboardProduct[index0]].name_ja
      );
    $('#name1') //
      .addClass('h2') //
      .text(
        miniProductMap[signboardProduct[index1]].name_ja
    );
    $('#name2') //
      .addClass('h2') //
      .text(
        miniProductMap[signboardProduct[index2]].name_ja
      );

    // ---------------------------------------------
    $('#price0') //
      .addClass('h1') //
      .text(
        miniProductMap[signboardProduct[index0]].price + "円"
      );

    $('#price1') //
      .addClass('h1') //
      .text(
        miniProductMap[signboardProduct[index1]].price + "円"
      );

    $('#price2') //
      .addClass('h1') //
      .text(
        miniProductMap[signboardProduct[index2]].price + "円"
      );


    index0 = (index0 + 1) % maxLength;
    index1 = (index1 + 1) % maxLength;
    index2 = (index2 + 1) % maxLength;
  }, 2000);
}



async function buildProductList() {
  $.each(miniProductMap, function (_, product) {

    const divParent = $('<div>') //
      .appendTo('#menu_content') //
      .addClass('col-6 col-md-3 alert products p-1 mb-0') // Updated for responsiveness
      .hide() //
      .attr('tag', product.tag) //
      .attr('id', product.jan_code)
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




// ページロード時の初期化
