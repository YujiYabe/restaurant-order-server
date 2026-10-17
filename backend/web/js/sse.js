document.addEventListener('alpine:init', () => {
  const kind = encodeURIComponent(document.getElementById('kind')?.getAttribute('name') ?? '');
  const placeName = encodeURIComponent(document.getElementById('place_name')?.getAttribute('name') ?? '');
  const urlString = `/v1/sse/app?kind=${kind}&place_name=${placeName}`;
  const connection = Alpine.store('connection');
  const eventSource = new EventSource(urlString);

  eventSource.addEventListener('open', () => {
    connection.connected = true;
    connection.error = null;
  });

  eventSource.addEventListener('message', (event) => {
    try {
      buildOnSse(event.data);
      connection.lastMessageAt = new Date();
      connection.error = null;
    } catch (error) {
      connection.error = error;
      console.error('SSE message handling failed:', error);
    }
  });

  eventSource.addEventListener('error', (error) => {
    connection.connected = false;
    connection.error = error;
    console.error('SSE connection failed:', error);
  });

  window.addEventListener('beforeunload', () => eventSource.close(), { once: true });
});
