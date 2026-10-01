var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
function selectSong(songId) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            const response = yield fetch(`/api/lyrics/${songId}`);
            const data = yield response.json();
            const lyricsContainer = document.getElementById('lyrics-text');
            const sectionParoles = document.getElementById('lyrics-container');
            if (lyricsContainer && sectionParoles) {
                lyricsContainer.innerHTML = '';
                Object.entries(data.parole).forEach(([sectionName, lines]) => {
                    const sectionDiv = document.createElement('div');
                    sectionDiv.className = 'lyric-section';
                    sectionDiv.innerHTML = `<h5>${sectionName}</h5>`;
                    lines.forEach(line => {
                        const p = document.createElement('p');
                        p.textContent = line;
                        sectionDiv.appendChild(p);
                    });
                    lyricsContainer.appendChild(sectionDiv);
                });
                sectionParoles.style.display = 'block';
            }
        }
        catch (error) {
            console.error("Erreur lors du chargement des paroles:", error);
        }
    });
}
// ---------------------------------------------------------------------------
// Utilitaires
// ---------------------------------------------------------------------------
function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}
/** "0:1.30" -> 1.3 | "1:03.24" -> 63.24 | 12.5 -> 12.5 */
function parseTime(raw) {
    if (typeof raw === 'number')
        return raw;
    return raw.trim().replace(',', '.').split(':').reduce((acc, p) => acc * 60 + parseFloat(p), 0);
}
const WORD_RE = /[a-zA-ZÀ-ÿ]/;
function normalize(s) {
    return s.toLowerCase().replace(/œ/g, 'oe').normalize('NFD').replace(/[\u0300-\u036f]/g, '')
        .replace(/[^a-z0-9\s]/g, ' ').replace(/\s+/g, ' ').trim();
}
function buildMaskedLine(line, count) {
    const tokens = line.split(/(\s+)/);
    const wordIdx = tokens.flatMap((t, i) => (WORD_RE.test(t) ? [i] : []));
    if (wordIdx.length === 0)
        return { parts: [{ text: line, blank: false }], answer: '' };
    const n = Math.min(count, wordIdx.length);
    const start = getRandomInt(0, wordIdx.length - n);
    const first = wordIdx[start];
    const last = wordIdx[start + n - 1];
    const parts = [];
    if (first > 0)
        parts.push({ text: tokens.slice(0, first).join(''), blank: false });
    parts.push({ text: tokens.slice(first, last + 1).join(''), blank: true });
    if (last < tokens.length - 1)
        parts.push({ text: tokens.slice(last + 1).join(''), blank: false });
    return { parts, answer: tokens.slice(first, last + 1).join('') };
}
/** Choisit une ligne qui contient assez de mots (évite de masquer "Alors" avec 8 mots à cacher). */
function pickLineToMask(lines, wanted) {
    const counts = lines.map(l => (l.text.match(/[a-zA-ZÀ-ÿ]\S*/g) || []).length);
    const ok = counts.flatMap((c, i) => (c >= wanted ? [i] : []));
    if (ok.length)
        return ok[getRandomInt(0, ok.length - 1)];
    return counts.indexOf(Math.max(...counts));
}
const CUT_LEAD_SECONDS = 0; // > 0 : la musique s'arrête un peu avant la ligne à compléter
let currentRound = null; // destroy() du round en cours
function renderLine(rl) {
    const p = rl.el;
    if (rl.kind === 'locked') {
        p.textContent = '🔒 [Paroles masquées]';
        p.classList.add('same-song-hidden-line');
    }
    else if (rl.kind === 'masked' && rl.masked) {
        p.classList.add('masked');
        p.textContent = '';
        rl.blanks = [];
        rl.masked.parts.forEach(part => {
            const span = document.createElement('span');
            if (part.blank) {
                span.className = 'blank';
                span.textContent = part.text.replace(/\S+/g, '_');
                rl.blanks.push(span);
            }
            else {
                span.textContent = part.text;
            }
            p.appendChild(span);
        });
    }
    else {
        p.textContent = rl.text;
    }
}
function startSync(opts) {
    const { audio, lines, box, mode, points } = opts;
    const cutIndex = lines.findIndex(l => l.kind === 'masked');
    let cutDone = cutIndex < 0;
    let activeIndex = -2;
    let raf = 0;
    let panel = null;
    // Fenêtre de 3 lignes : précédente, active, suivante. Tout le reste est masqué (display: none).
    const applyStates = (active) => {
        lines.forEach((l, i) => {
            l.el.classList.toggle('is-active', i === active);
            l.el.classList.toggle('is-past', i < active);
            l.el.classList.toggle('is-upcoming', i > active);
            l.el.classList.toggle('is-hidden', i < active - 1 || i > active + 1);
        });
        // une section sans ligne visible disparaît aussi (marges comprises)
        new Set(lines.map(l => l.el.parentElement)).forEach(sec => {
            sec === null || sec === void 0 ? void 0 : sec.classList.toggle('is-hidden', !sec.querySelector('.lyric-line:not(.is-hidden)'));
        });
    };
    box.classList.remove('is-finished');
    applyStates(-1); // avant le 1er timecode : seule la 1re ligne (la "suivante") est visible
    activeIndex = -1;
    const tick = () => {
        const t = audio.currentTime;
        let active = -1;
        for (let i = 0; i < lines.length; i++) {
            if (lines[i].start <= t)
                active = i;
            else
                break;
        }
        if (active !== activeIndex) {
            activeIndex = active;
            applyStates(active);
        }
        if (!cutDone && t >= lines[cutIndex].start - CUT_LEAD_SECONDS) {
            cutDone = true;
            audio.pause();
            openAnswerPanel();
            return;
        }
        raf = requestAnimationFrame(tick);
    };
    const onPlay = () => { cancelAnimationFrame(raf); raf = requestAnimationFrame(tick); };
    const onPause = () => cancelAnimationFrame(raf);
    const onSeeking = () => {
        // anti-triche : impossible de sauter au-delà du point de coupure
        if (!cutDone && audio.currentTime > lines[cutIndex].start) {
            audio.currentTime = Math.max(0, lines[cutIndex].start - 0.05);
        }
    };
    const onEnded = () => {
        cancelAnimationFrame(raf);
        // fin du morceau : on ré-affiche toute la chanson, défilable
        box.classList.add('is-finished');
        document.querySelectorAll('.is-hidden').forEach(e => { if (box.contains(e))
            e.classList.remove('is-hidden'); });
        lines.forEach(l => {
            l.el.classList.remove('is-upcoming', 'is-active');
            l.el.classList.add('is-past');
            if (l.kind === 'locked') {
                l.el.textContent = l.text;
                l.el.classList.remove('same-song-hidden-line');
                l.el.classList.add('revealed');
            }
        });
    };
    audio.addEventListener('play', onPlay);
    audio.addEventListener('pause', onPause);
    audio.addEventListener('seeking', onSeeking);
    audio.addEventListener('ended', onEnded);
    // --- Phase "réponse" : écrire ou chanter -------------------------------
    function openAnswerPanel() {
        var _a;
        const rl = lines[cutIndex];
        const expected = rl.masked.answer;
        panel = document.createElement('div');
        panel.className = 'answer-panel';
        panel.innerHTML = `
            <p class="answer-hint">🎵 La musique s'est arrêtée : complète la ligne !</p>
            <div class="answer-row">
                <input type="text" class="answer-input" placeholder="Les mots manquants…" autocomplete="off">
                <button type="button" class="answer-btn answer-submit">Valider</button>
                <button type="button" class="answer-btn answer-sing">🎤 Chanter</button>
            </div>
            <p class="answer-feedback" aria-live="polite"></p>`;
        ((_a = box.closest('.lyrics-container')) !== null && _a !== void 0 ? _a : box).after(panel);
        const input = panel.querySelector('.answer-input');
        const submit = panel.querySelector('.answer-submit');
        const sing = panel.querySelector('.answer-sing');
        const feedback = panel.querySelector('.answer-feedback');
        input.focus();
        const validate = () => {
            var _a;
            const given = normalize(input.value);
            // tolérant : accepte les mots attendus seuls OU une ligne entière chantée qui les contient
            const success = given.length > 0 && ` ${given} `.includes(` ${normalize(expected)} `);
            (_a = rl.blanks) === null || _a === void 0 ? void 0 : _a.forEach(b => {
                b.textContent = expected;
                b.classList.add(success ? 'is-right' : 'is-wrong');
            });
            feedback.textContent = success ? `✅ Bravo ! +${points} points` : `❌ C'était : « ${expected} »`;
            feedback.className = `answer-feedback ${success ? 'is-right' : 'is-wrong'}`;
            input.disabled = true;
            sing.hidden = true;
            // hook pour gamelogic.js : envoyer le score au serveur, passer au joueur suivant, etc.
            document.dispatchEvent(new CustomEvent('lyrics:answer', {
                detail: { success, points: success ? points : 0, expected, given: input.value },
            }));
            submit.textContent = '▶ Écouter la suite';
            submit.onclick = () => { panel === null || panel === void 0 ? void 0 : panel.remove(); panel = null; audio.play(); };
        };
        submit.onclick = validate;
        input.addEventListener('keydown', e => { if (e.key === 'Enter')
            validate(); });
        // Mode "chanter" : Web Speech API (Chrome / Edge / Safari, HTTPS requis)
        const SR = window.SpeechRecognition || window.webkitSpeechRecognition;
        if (!SR) {
            sing.hidden = true;
        }
        else {
            sing.onclick = () => {
                const rec = new SR();
                rec.lang = 'fr-FR';
                rec.interimResults = true;
                rec.onresult = (e) => {
                    input.value = Array.from(e.results).map((r) => r[0].transcript).join(' ');
                };
                rec.onend = () => { sing.classList.remove('is-listening'); if (input.value)
                    validate(); };
                sing.classList.add('is-listening');
                rec.start();
            };
        }
    }
    return () => {
        cancelAnimationFrame(raf);
        audio.removeEventListener('play', onPlay);
        audio.removeEventListener('pause', onPause);
        audio.removeEventListener('seeking', onSeeking);
        audio.removeEventListener('ended', onEnded);
        panel === null || panel === void 0 ? void 0 : panel.remove();
    };
}
// ---------------------------------------------------------------------------
// Initialisation d'un round
// ---------------------------------------------------------------------------
window.initLyrics = function (songFileName, points, targetId) {
    return __awaiter(this, void 0, void 0, function* () {
        var _a;
        try {
            currentRound === null || currentRound === void 0 ? void 0 : currentRound(); // nettoie le round précédent (rAF + listeners)
            currentRound = null;
            const response = yield fetch(`/api/lyrics/${encodeURIComponent(songFileName)}`);
            if (!response.ok)
                throw new Error("Erreur lors de la récupération du JSON");
            const data = yield response.json();
            const container = document.getElementById(targetId);
            if (!container)
                return;
            container.innerHTML = '';
            const sectionsKeys = Object.keys(data.parole);
            const refrains = sectionsKeys.filter(k => k.toLowerCase().includes('refrain'));
            const couplets = sectionsKeys.filter(k => k.toLowerCase().includes('couplet'));
            const firstRefrainKey = refrains[0] || null;
            const secondRefrainKey = refrains[1] || refrains[0] || null;
            const secondCoupletKey = couplets[1] || couplets[0] || null;
            let mode = 'points';
            let wordsToMask = 0;
            let targetSectionKey = null;
            const levelId = String(points);
            const pts = levelId === 'same-song' ? 0 : Number.parseInt(levelId.split('-')[0], 10);
            if (levelId === 'same-song') {
                mode = 'same-song';
            }
            else if (pts === 50) {
                wordsToMask = getRandomInt(8, 10);
                targetSectionKey = secondRefrainKey;
            }
            else if (pts === 40) {
                wordsToMask = getRandomInt(5, 7);
                targetSectionKey = firstRefrainKey;
            }
            else if (pts === 30) {
                wordsToMask = getRandomInt(4, 5);
                targetSectionKey = secondCoupletKey;
            }
            else if (pts === 20) {
                wordsToMask = 3;
                targetSectionKey = secondCoupletKey;
            }
            else if (pts === 10) {
                wordsToMask = 2;
                targetSectionKey = secondCoupletKey;
            }
            const sameSongCutLine = getRandomInt(3, 5);
            let lineCounter = 0;
            const roundLines = [];
            sectionsKeys.forEach(sectionName => {
                const lines = data.parole[sectionName];
                const sectionDiv = document.createElement('div');
                sectionDiv.className = 'lyric-section mb-3';
                const title = document.createElement('h5');
                title.textContent = sectionName;
                title.style.visibility = 'hidden';
                sectionDiv.appendChild(title);
                const isTarget = targetSectionKey !== null && sectionName === targetSectionKey;
                const maskIdx = isTarget && lines.length > 0 ? pickLineToMask(lines, wordsToMask) : -1;
                lines.forEach((lineObj, index) => {
                    const p = document.createElement('p');
                    p.className = 'lyric-line is-upcoming';
                    const rl = { el: p, start: parseTime(lineObj.time), text: lineObj.text, kind: 'normal' };
                    if (mode === 'same-song') {
                        if (lineCounter >= sameSongCutLine)
                            rl.kind = 'locked';
                        lineCounter++;
                    }
                    else if (index === maskIdx) {
                        rl.kind = 'masked';
                        rl.masked = buildMaskedLine(lineObj.text, wordsToMask);
                    }
                    renderLine(rl);
                    sectionDiv.appendChild(p);
                    roundLines.push(rl);
                });
                container.appendChild(sectionDiv);
            });
            roundLines.sort((a, b) => a.start - b.start); // sécurité si le JSON n'est pas ordonné
            const audioPlayer = document.getElementById(`audio-player-${points}`);
            if (!audioPlayer)
                return;
            const audioFileName = songFileName.replace('.json', '.mp3')
                .toLowerCase().normalize('NFD').replace(/[\u0300-\u036f]/g, '');
            audioPlayer.src = `/api/musiques/${encodeURIComponent(audioFileName)}`;
            audioPlayer.load();
            const box = (_a = container.closest('.lyrics-container')) !== null && _a !== void 0 ? _a : container;
            currentRound = startSync({ audio: audioPlayer, lines: roundLines, box, mode, points: pts });
            audioPlayer.play().catch(error => console.error("Erreur d'autoplay :", error));
        }
        catch (error) {
            console.error("Erreur globale lors de l'initialisation du round :", error);
        }
    });
};
document.playAudio = function (filename) {
    const audio = new Audio(`https://asset.nolp-jeu.fr/musiques/${encodeURIComponent(filename)}`);
    audio.play().catch(error => {
        console.error("Erreur lors de la lecture de l'audio:", error);
    });
};
export {};
