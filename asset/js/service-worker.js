"use strict";
/// <reference lib="webworker" />
// ✅ Type correct pour le scope global du service worker
const CACHE_NAME = 'nopl-cache-v1';
const urlsToCache = [
    '/',
    '/duel',
    '/asset/style.css',
    '/asset/js/duel.js',
    '/asset/js/ui.js',
];
// ✅ Gestion de l'installation du service worker
self.addEventListener('install', (event) => {
    event.waitUntil(caches.open(CACHE_NAME).then(cache => cache.addAll(urlsToCache)));
});
// ✅ Interception des requêtes HTTP pour utiliser le cache
self.addEventListener('fetch', (event) => {
    event.respondWith(caches.match(event.request).then(response => {
        return response || fetch(event.request);
    }));
});
