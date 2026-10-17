(() => {
  'use strict';

  const previewSize = 240;
  const previewGap = 18;
  const previewPadding = 26;

  function updatePreviewPosition(preview, event) {
    let left = event.clientX + previewGap;
    let top = event.clientY + previewGap;
    const outerSize = previewSize + previewPadding;

    if (left + outerSize > window.innerWidth) {
      left = event.clientX - outerSize - previewGap;
    }
    if (top + outerSize > window.innerHeight) {
      top = event.clientY - outerSize - previewGap;
    }

    preview.style.left = `${Math.max(8, left)}px`;
    preview.style.top = `${Math.max(8, top)}px`;
  }

  function showPreview(preview, link, event) {
    const absoluteURL = new URL(link.href, window.location.href).href;
    preview.replaceChildren();
    window.$(preview).qrcode({
      width: previewSize,
      height: previewSize,
      text: absoluteURL,
    });
    preview.hidden = false;
    preview.setAttribute('aria-hidden', 'false');
    updatePreviewPosition(preview, event);
    requestAnimationFrame(() => preview.classList.add('is-visible'));
  }

  function hidePreview(preview) {
    preview.classList.remove('is-visible');
    preview.setAttribute('aria-hidden', 'true');
    window.setTimeout(() => {
      if (!preview.classList.contains('is-visible')) {
        preview.hidden = true;
      }
    }, 160);
  }

  function initIndex() {
    const preview = document.getElementById('qr-preview');
    if (!preview || typeof window.$ !== 'function' || !window.$.fn || typeof window.$.fn.qrcode !== 'function') {
      return;
    }

    document.querySelectorAll('.portal-card, .place-links a').forEach((link) => {
      link.addEventListener('mouseenter', (event) => showPreview(preview, link, event));
      link.addEventListener('mousemove', (event) => updatePreviewPosition(preview, event));
      link.addEventListener('mouseleave', () => hidePreview(preview));
    });
  }

  window.addEventListener('load', initIndex);
})();
