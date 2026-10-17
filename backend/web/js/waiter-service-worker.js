const CACHE_NAME = 'restaurant-order-server-cache-v5';

const BOOTSTRAP_JS_URL = 'https://cdn.jsdelivr.net/npm/bootstrap@5.3.1/dist/js/bootstrap.bundle.min.js';
const ALPINE_JS_URL = 'https://cdn.jsdelivr.net/npm/alpinejs@3.17.4/dist/cdn.min.js';

const CACHE_FILES = [
  '/js/waiter-manifest.json',
  '/css/common.css',
  '/js/common.js',
  '/js/external/cash.min.js',
  ALPINE_JS_URL,
  '/js/alpine_app.js',
  BOOTSTRAP_JS_URL,
  '/image/app-icon.svg',
];

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.addAll(CACHE_FILES);
    })
  );
});

// 古いキャッシュ削除（そのまま）
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((name) => {
          if (name !== CACHE_NAME) {
            return caches.delete(name);
          }
        })
      );
    })
  );
});

// フェッチ処理
self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.match(event.request).then((response) => {
      return response || fetch(event.request);
    })
  );
});
