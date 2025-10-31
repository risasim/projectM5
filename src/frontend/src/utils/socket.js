export function useLeaderboardSocket(onMessage) {
  const socket = new WebSocket("ws://localhost:8080/ws/leaderboard");

  socket.onopen = () => console.log("Connected to leaderboard WebSocket");
  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    onMessage(data);
  };
  socket.onerror = (err) => console.error("WebSocket error:", err);
  socket.onclose = () => console.log("WebSocket closed");

  return socket;
}



// idk this is just staying here for now
export function useLeaderboardSocket(onMessage) {
  const token = localStorage.getItem('authToken');
  if (!token) {
    console.error("No auth token found, redirecting to login...");
    return null;
  }

  const socket = new WebSocket(`ws://116.203.97.62:8080/ws/leaderboard?token=${token}`);

  socket.onopen = () => console.log("Connected to leaderboard WebSocket");
  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log("Received:", data);
    onMessage(data);
  };
  socket.onerror = (err) => console.error("WebSocket error:", err);
  socket.onclose = () => console.log("WebSocket closed");

  return socket;
}
