
function toggleElement(elementId) {
    const element = document.getElementById(elementId);
    if (element) {
        if (element.style.display === 'none' || element.classList.contains('hidden-element')) {
            element.style.display = 'block';
            element.classList.remove('hidden-element');
            element.classList.add('visible');
        } else {
            element.style.display = 'none';
            element.classList.add('hidden-element');
            element.classList.remove('visible');
        }
    }
}

menutogglebtn.onclick = () => {
    sidebarmenu.classList.toggle('is-open');
};

document.onclick = (e) => {
    if (!sidebarmenu.contains(e.target) && !menutogglebtn.contains(e.target)) {
        sidebarmenu.classList.remove('is-open');
    }
};

window.toggleLevelSongs = (levelId) => {
    const target = document.getElementById(levelId);
    if (!target) return;
    const isVisible = target.style.display === 'block';
    document.querySelectorAll('.songs-for-level').forEach(div => {
        div.style.display = 'none';
    });
    target.style.display = isVisible ? 'none' : 'block';
};

window.toggleElement = toggleElement;
window.toggleLevelSongs = toggleLevelSongs;