/// <reference lib="webworker" />


const CACHE_NAME = 'nopl-cache-v1';

const urlsToCache: string[] = [
  '/',
  '/duel',
  '/asset/style.css',
  '/asset/js/duel.js',
  '/asset/js/ui.js',
];

self.addEventListener('install', (event: ExtendableEvent) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then(cache => cache.addAll(urlsToCache))
  );
});

self.addEventListener('fetch', (event: FetchEvent) => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
