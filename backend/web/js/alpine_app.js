document.addEventListener('alpine:init', () => {
  Alpine.store('connection', {
    connected: false,
    lastMessageAt: null,
    error: null,
  });

  Alpine.data('restaurantOrderServerPage', () => ({
    async init() {
      if (typeof window.commonInit === 'function') {
        await window.commonInit();
      }

      const initializerName = document.body.dataset.pageInitializer || 'init';
      const initializer = window[initializerName];
      if (typeof initializer === 'function') {
        await initializer();
      }
    },
  }));
});
