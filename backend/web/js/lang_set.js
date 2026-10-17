const defaultLangCode = 1041;

const langSet = {
  1041: {
    // name: '日本語',
    // translate: 'jp',
    // speech: 'ja-JP',
    // prepare_txt: ' が準備できました。',

    // notice_header: '振動と音声でお知らせします',
    // notice_body: '(機能がONになっている必要があります)',

    // wait_number: 'お待ち番号',
    // call_number: 'お呼び出し中の番号',

    // price_discount: '現在の値引き額',
    // price_hour_text:
    //   '当店の商品は昼食時の価格を基準に、ダイナミックプライシングを導入しております。',
    // price_multi_text: '2点目以降、全商品100円引きとなります。',

    // step1_text: '以下のカードを選んでお盆に載せてください。',
    // step1_card_lang: '言語カード（日本語の場合は必要ありません）',
    // step1_card_pay: '支払い方法カード',
    // step1_card_prd: '商品カード',

    // calorie_name: 'カロリー',
    // allergy_name: 'アレルギー情報',
    // tax_name: ' 円(税込み)',
    // step2_price: '合計金額',

    // step2_text:
    //   '列にお待ちの間 モニターで以下をご確認いただき、会計の準備をお願いします。',
    // step3_text:
    //   'レジにてお会計後、セルフサービスでお受け取り願います。お呼び出しが前後する可能性があります。',
    // step4_text: 'お食事後の食器はセルフサービスでお片付け願います。',

    // inactive_hour: '営業時間外',

    // 乳卵
    egg: '卵',
    milk: '乳',

    // 魚介
    shrimp: 'えび',
    crab: 'かに',
    squid: 'いか',
    salmon: 'さけ',
    mackerel: 'さば',
    salmon_roe: 'いくら',
    abalone: 'あわび',

    // 種子
    peanut: 'ピーナッツ',
    walnut: 'くるみ',
    soybean: '大豆',
    almond: 'アーモンド',
    cashew_nut: 'カシューナッツ',
    sesame: 'ごま',

    // 穀物
    wheat: '小麦',
    buckwheat: 'そば',

    // フルーツ
    apple: 'りんご',
    orange: 'オレンジ',
    banana: 'バナナ',
    peach: 'もも',
    kiwi_fruit: 'キウイフルーツ',

    // 肉
    pork: '豚肉',
    chicken: '鶏肉',
    beef: '牛肉',

    // その他
    wild_yam: 'やまいも',
    matsutake_mushroom: 'まつたけ',
    gelatin: 'ゼラチン',
  },
  1033: {
    // name: 'english',
    // translate: 'en',
    // speech: 'en-US',
    // prepare_txt: ' is ready.',

    // notice_header: 'Vibration and voice notification',
    // notice_body: '(function must be turned on)',

    // wait_number: 'Waiting number',
    // call_number: 'Calling number',

    // price_discount: 'Current discounts',
    // price_hour_text:
    //   'Our products are priced based on lunch prices and we have introduced dynamic pricing.',
    // price_multi_text:
    //   'From the second item onward, all products will be discounted by 100 yen.',

    // step1_text: 'Please select a card below and place it on the tray.',
    // step1_card_lang: 'Language cards',
    // step1_card_pay: 'Payment method cards',
    // step1_card_prd: 'Product cards',

    // calorie_name: 'calorie',
    // allergy_name: 'allergy info',
    // tax_name: ' yen(tax included)',

    // step2_price: 'total amount',

    // step2_text:
    //   'While waiting in line, please check the following on the monitor and prepare your bill.',
    // step3_text:
    //   'Please pick up your order by self-service after paying at the cashier. Please note that the call may be delayed.',
    // step4_text: 'Please clean up your dishes after meals by self-service.',

    // inactive_hour: 'in preparation',

    // 乳卵
    egg: 'Egg',
    milk: 'Milk',

    // 魚介
    shrimp: 'Shrimp',
    crab: 'Crab',
    squid: 'Squid',
    salmon: 'Salmon',
    mackerel: 'Mackerel',
    salmon_roe: 'Salmon roe',
    abalone: 'Abalone',

    // 種子
    peanut: 'Peanut',
    walnut: 'Walnut',
    soybean: 'Soybean',
    almond: 'Almond',
    cashew_nut: 'Cashew nut',
    sesame: 'Sesame',

    // 穀物
    wheat: 'Wheat',
    buckwheat: 'Buckwheat',

    // フルーツ
    apple: 'Apple',
    orange: 'Orange',
    banana: 'Banana',
    peach: 'Peach',
    kiwi_fruit: 'Kiwi fruit',

    // 肉
    pork: 'Pork',
    chicken: 'Chicken',
    beef: 'Beef',

    // その他
    wild_yam: 'Wild yam',
    matsutake_mushroom: 'Matsutake mushroom',
    gelatin: 'Gelatin',
  },
  1034: {
    // name: 'español',
    // translate: 'es',
    // speech: 'es-ES',
    // prepare_txt: ' está listo.',

    // notice_header: 'Notificación por vibración y audio.',
    // notice_body: '(la función debe estar activada).',

    // wait_number: 'Número de espera',
    // call_number: 'Número al que llama',

    // price_discount: 'Descuentos actuales',
    // price_hour_text:
    //   'Los precios de nuestros productos se basan en los precios de los almuerzos y hemos introducido precios dinámicos.',
    // price_multi_text:
    //   'A partir del segundo artículo, todos los productos tendrán un descuento de 100 yenes.',

    // step1_text: 'Seleccione una tarjeta y colóquela en la bandeja.',
    // step1_card_lang: 'Tarjetas de idiomas',
    // step1_card_pay: 'Tarjetas forma de pago',
    // step1_card_prd: 'Tarjetas de producto',

    // calorie_name: 'caloría',
    // allergy_name: 'Información sobre alergias',
    // tax_name: ' yen(impuesto incluido)',
    // step2_price: 'importe total',

    // step2_text:
    //   'Mientras espera en la fila, compruebe lo siguiente en el monitor y prepare su factura',
    // step3_text:
    //   'Por favor, recoja su pedido en la caja y en el autoservicio. Las llamadas pueden sufrir retrasos.',
    // step4_text: 'Por favor, limpie después de su comida por autoservicio.',

    // inactive_hour: 'en preparación',

    // 乳卵
    egg: 'Huevo',
    milk: 'Leche',

    // 魚介
    shrimp: 'Camarón',
    crab: 'Cangrejo',
    squid: 'Calamar',
    salmon: 'Salmón',
    mackerel: 'Caballa',
    salmon_roe: 'Huevas de salmón',
    abalone: 'Abulón',

    // 種子
    peanut: 'Maní',
    walnut: 'Nuez',
    soybean: 'Soja',
    almond: 'Almendra',
    cashew_nut: 'Anacardo',
    sesame: 'Sésamo',

    // 穀物
    wheat: 'Trigo',
    buckwheat: 'Trigo sarraceno',

    // フルーツ
    apple: 'Manzana',
    orange: 'Naranja',
    banana: 'Plátano',
    peach: 'Durazno',
    kiwi_fruit: 'Kiwi',

    // 肉
    pork: 'Cerdo',
    chicken: 'Pollo',
    beef: 'Carne de res',

    // その他
    wild_yam: 'Ñame silvestre',
    matsutake_mushroom: 'Seta matsutake',
    gelatin: 'Gelatina',
  },
  2052: {
    // name: '中文-简体',
    // translate: 'zh',
    // speech: 'zh-CN',
    // prepare_txt: ' 准备好了',

    // notice_header: '振动和声音通知',
    // notice_body: '（功能必须打开）',

    // wait_number: '等候号码',
    // call_number: '您要拨打的号码',

    // price_discount: '当前折扣',
    // price_hour_text: '我们的产品以午餐价格为基础定价，并引入了动态定价',
    // price_multi_text: '从第2件开始，所有商品均减价100日元',

    // step1_text: '请在下面选择一张卡片并将其放在托盘上',
    // step1_card_lang: '语言卡',
    // step1_card_pay: '付款方式 卡',
    // step1_card_prd: '产品卡',

    // calorie_name: '卡路里',
    // allergy_name: '过敏信息',
    // tax_name: ' yen(含税)',
    // step2_price: '总金额',

    // step2_text: '排队等候时，请检查显示器上的以下内容，并准备好您的账单',
    // step3_text: '请在收银台自助取餐。 叫号可能会延迟。',
    // step4_text: '餐后请自助清理餐具。',

    // inactive_hour: '准备中',

    // 乳卵
    egg: '鸡蛋',
    milk: '牛奶',

    // 魚介
    shrimp: '虾',
    crab: '蟹',
    squid: '鱿鱼',
    salmon: '鲑鱼',
    mackerel: '鲭鱼',
    salmon_roe: '鲑鱼籽',
    abalone: '鲍鱼',

    // 種子
    peanut: '花生',
    walnut: '核桃',
    soybean: '大豆',
    almond: '杏仁',
    cashew_nut: '腰果',
    sesame: '芝麻',

    // 穀物
    wheat: '小麦',
    buckwheat: '荞麦',

    // フルーツ
    apple: '苹果',
    orange: '橙子',
    banana: '香蕉',
    peach: '桃子',
    kiwi_fruit: '猕猴桃',

    // 肉
    pork: '猪肉',
    chicken: '鸡肉',
    beef: '牛肉',

    // その他
    wild_yam: '山药',
    matsutake_mushroom: '松茸',
    gelatin: '明胶',
  },

  14337: {
    // 乳卵
    egg: 'بيض',
    milk: 'حليب',

    // 魚介
    shrimp: 'روبيان',
    crab: 'سرطان البحر',
    squid: 'حبار',
    salmon: 'سمك السلمون',
    mackerel: 'سمك الأسماك',
    salmon_roe: 'بيض السلمون',
    abalone: 'أبلون',

    // 種子
    peanut: 'فول السوداني',
    walnut: 'جوز',
    soybean: 'فول الصويا',
    almond: 'لوز',
    cashew_nut: 'كاجو',
    sesame: 'سمسم',

    // 穀物
    wheat: 'قمح',
    buckwheat: 'الغرينية',

    // フルーツ
    apple: 'تفاحة',
    orange: 'البرتقال',
    banana: 'موز',
    peach: 'خوخ',
    kiwi_fruit: 'كيوي',

    // 肉
    pork: 'لحم الخنزير',
    chicken: 'دجاج',
    beef: 'لحم بقر',

    // その他
    wild_yam: 'اليام البري',
    matsutake_mushroom: 'فطر ماتسوتاكي',
    gelatin: 'جيلاتين',
  },
};
