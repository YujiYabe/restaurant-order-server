let video = document.getElementById('video');
let canvas = document.getElementById('canvas');
let videoFacingMode = 'environment';
let videoZoom = 2;
let defaultCertificateType = 'license_driver';
let certificateType = defaultCertificateType;
let isFace = false;
let temporaryUUID = '';
let uuidVerifyTwenty = '';
let isIncludeList = false;
async function init() {
  isEnvLocal = $('#environment').attr('name') == 'local';

  loadModule();

  uuidVerifyTwenty = await getLocalStorage('uuidVerifyTwenty');

  // 登録済み
  if (uuidVerifyTwenty != '') {
    changeScreen('screen_main');

    isFace = true;
    toggleCamera(isFace);
    
    
    isIncludeList = await getIsIncludeList()
    // 登録済みであるが本日まだ認証できていないケース
    if (isIncludeList) {
      await changeScreen('screen_already_exist');
    } else { 
      changeMode('mode_register');
    }

  } else { 
    await changeScreen('screen_wait');
  }
}
async function loadModule() {
  await setupCamera();
  $('#shutter_block').show();
  $('#explain_block').hide();
}

async function getIsIncludeList() {
  const res = await fetch(
    `${window.location.origin}/v1/alcohol/include_list?uuidVerifyTwenty=${uuidVerifyTwenty}`, // url
    {
      method: 'GET',
      headers: {
        'Content-type': 'application/json; charset=UTF-8',
      },
    }
  );

  if (res.status == 200) {
    return true;
  }

  return false;
}


async function toggleCamera(isFaceBool) {
  isFace = isFaceBool;
  if (isFace) {
    videoFacingMode = 'user'; // フロントカメラ
    videoZoom = 1;
    $('#explain_certificate_block').hide();
    $('#explain_latest_face_block').show();
  } else {
    videoFacingMode = 'environment'; // バックカメラ
    videoZoom = 2;
    $('#explain_certificate_block').show();
    $('#explain_latest_face_block').hide();
  }
  
  await setupCamera();
}


async function setupCamera() {
  const stream = await navigator.mediaDevices.getUserMedia({
    video: {
      facingMode: videoFacingMode,
      zoom: videoZoom,
      width: { min: 640, ideal: 1920, max: 7680 }, // 最大解像度8K
      height: { min: 480, ideal: 1080, max: 4320 }, // 4K以上の解像度を指定
      aspectRatio: { ideal: 1.7777777778 }, // 16:9のアスペクト比
      frameRate: { min: 15, ideal: 30, max: 60 }, // フレームレートの設定
      focusMode: 'continuous',
      whiteBalanceMode: 'continuous',
      exposureMode: 'continuous',
      // torch: true,
    },
  });
  video.srcObject = stream;
}

async function takePhoto() {
  canvas.width = video.videoWidth;
  canvas.height = video.videoHeight;
  const context = canvas.getContext('2d');
  context.drawImage(video, 0, 0, canvas.width, canvas.height);

  const base64Image = canvas.toDataURL('image/jpeg', 0.9);

  if (isFace) {
    if (uuidVerifyTwenty == '') {
      await registerFace(base64Image);
    } else { 
      await verifyFace(base64Image);
    }
    
  } else {
    await postCertificate(base64Image);
  }
}

async function postCertificate(base64Image) {

  const data = JSON.stringify({
    certificate_type: certificateType,
    base64Source: base64Image,
  });

  try {
    changeScreen('screen_auth_attempt');
    const response = await fetch(
      '/v1/alcohol/verify_twenty',
      {
        method: 'POST',
        headers: {
          'Content-type': 'application/json; charset=UTF-8',
        },
        body: data,
      }
    );

    const result = await response.json(); // Assuming the server responds with JSON
    if (response.status === 200) {

      $('#age').text(result.age);
      $('#birthday').text(result.birthday);
      temporaryUUID = result.uuid

      changeScreen('screen_auth_success');

      // URLとカメラモードを変更
      await toggleCamera(true);
      changeMode('mode_register');
      return;
    }
    
    changeScreen('screen_auth_fail');
  } catch (error) {
    changeScreen('screen_auth_fail');
    console.error('Error:', error);
  }
}

async function registerFace(base64Image) {
  const data = JSON.stringify({
    base64Source: base64Image,
    uuid: temporaryUUID,
    certificate_type: '',
  });

  try {
    changeScreen('screen_auth_attempt');
    const response = await fetch(
      '/v1/alcohol/register_face',
      {
        method: 'POST',
        headers: {
          'Content-type': 'application/json; charset=UTF-8',
        },
        body: data,
      }
    );

    const result = await response.json(); // Assuming the server responds with JSON
    if (response.status === 200) {

      await localStorage.setItem('uuidVerifyTwenty', result);

      changeScreen('screen_register_face_finished');

      return;
    }
    
    changeScreen('screen_latest_face_fail');
  } catch (error) {
  changeScreen('screen_latest_face_fail');
    console.error('Error:', error);
  }
}

async function verifyFace(base64Image) {


  const data = JSON.stringify({
    base64Source: base64Image,
    uuid: uuidVerifyTwenty,
    certificate_type: '',
  });

  try {
    changeScreen('screen_auth_attempt');
    const response = await fetch(
      '/v1/alcohol/verify_face',
      {
        method: 'POST',
        headers: {
          'Content-type': 'application/json; charset=UTF-8',
        },
        body: data,
      }
    );

    const result = await response.json(); // Assuming the server responds with JSON
    if (response.status === 200) {

      changeScreen('screen_verify_face_finished');

      return;
    }
    
    changeScreen('screen_latest_face_fail');
  } catch (error) {
  changeScreen('screen_latest_face_fail');
    console.error('Error:', error);
  }
}

function returnOrderPage () {
  window.close()
}

function changeScreen(screenName) {
  $('.screen').hide();
  $('#' + screenName).show();
}

function changeMode(modeName) {
  $('.mode').hide();
  $('#' + modeName).show();
}

function selectCertificate(selectedCertificateType) {
  certificateType = selectedCertificateType;

  $('.certificate_button').removeClass('alert-primary');
  $('.certificate_button').removeClass('alert-secondary');

  $('.certificate_button').addClass('alert-secondary');

  $('#' + selectedCertificateType).removeClass('alert-secondary');
  $('#' + selectedCertificateType).addClass('alert-primary');
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
