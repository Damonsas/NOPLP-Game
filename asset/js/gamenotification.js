export function showNotification(message, type = 'success') {
    const notif = document.createElement('div');
    notif.className = `notification ${type}`;
    notif.textContent = message;
    notif.style.cssText = `
        position: fixed; top: 20px; right: 20px; padding: 15px 20px;
        border-radius: 5px; color: white; font-weight: bold; z-index: 1000;
        max-width: 300px; word-wrap: break-word;
    `;
    const colors = { success: '#28a745', error: '#dc3545', warning: '#ffc107', info: '#17a2b8' };
    notif.style.backgroundColor = colors[type] || colors.info;
    if (type === 'warning')
        notif.style.color = '#212529';
    document.body.appendChild(notif);
    setTimeout(() => notif.remove(), 4000);
}
