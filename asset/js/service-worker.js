"use strict";
/// <reference lib="webworker" />
// @ts-nocheck
const CACHE_NAME = 'nopl-cache-v1';
const urlsToCache = [
    '/',
    '/duel',
    '/asset/style.css',
    '/asset/js/duel.js',
    '/asset/js/ui.js',
];
self.addEventListener('install', (event) => {
    event.waitUntil(caches.open(CACHE_NAME).then(cache => cache.addAll(urlsToCache)));
});
self.addEventListener('fetch', (event) => {
    event.respondWith(caches.match(event.request).then(response => {
        return response || fetch(event.request);
    }));
});
