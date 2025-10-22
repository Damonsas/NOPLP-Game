
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

document.addEventListener('DOMContentLoaded',() => {
    const menu = document.getElementById('sidebarmenu');
    const button = document.getElementById('menutogglebtn');
    button.addEventListener('click', () => {
        menu.classList.toggle('is-open');

    });
});


window.toggleElement = toggleElement;
