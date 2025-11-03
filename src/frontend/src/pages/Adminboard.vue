<template>
  <div class="page-container">
    <div class="adminboard-page">
      <div class="adminboard-container">

        <div class="top-buttons">
          <h1 class="adminboard-title">Game Session Settings</h1>
          <div class="button-group">
            <button class="adminboard-btn" @click="$router.push('/adminadduser')">Add User</button>
            <button class="adminboard-btn" @click="goToLeaderboard">Leaderboard</button>
          </div>
        </div>

        <div class="gametype-select">
          <label for="gametype">Choose gametype:</label>
          <select
            id="gametype"
            class="adminboard-select"
            v-model="gameMode"
            @change="onGameModeChange"
            :disabled="isGameActive"
          >
            <option value="Freefall">FreeFall</option>
            <option value="Infected">Infected</option>
            <option value="TeamDeathmatch">Team Deathmatch</option>
          </select>
          <p v-if="isGameActive" style="color: red; font-weight: 600; margin-top: 0.5rem;">
            Game active — cannot change mode
          </p>
        </div>

        <table class="player-table">
          <thead>
            <tr>
              <th>Player</th>
              <th>Team</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="players.length === 0">
              <td colspan="3">No players connected yet.</td>
            </tr>
            <tr v-for="player in players" :key="player.username">
              <td>{{ player.username }}</td>
              <td>{{ player.team || '-' }}</td>
              <td :class="{ online: player.online, offline: !player.online }">
                {{ player.online ? 'Online' : 'Offline' }}
              </td>
            </tr>
          </tbody>
        </table>

        <div class="session-buttons">
          <button class="start-session" @click="createGame">Create Game</button>
          <button class="start-session" @click="startGame">Start Game</button>
          <button class="end-session" @click="stopGame">Stop Game</button>
        </div>

        <div v-if="message" class="message-box">{{ message }}</div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AdminBoard',
  data() {
    return {
      message: '',
      gameMode: sessionStorage.getItem('selectedGameMode') || 'Freefall',
      players: [],
      isGameActive: false,
      serverGameStatus: 'Idle',
      pollingInterval: null
    };
  },

  methods: {
    /** 🔐 Retrieve token or redirect if not logged in */
    getAuthToken() {
      const token = sessionStorage.getItem('authToken');
      console.log('[getAuthToken] token:', token);
      if (!token) {
        alert('You must log in first.');
        console.warn('[getAuthToken] No token found, redirecting to /login');
        this.$router.push('/login');
        return null;
      }
      return token;
    },

    /** 🎮 Handle mode change */
    onGameModeChange() {
      console.log('[onGameModeChange] triggered:', this.gameMode);
      if (this.isGameActive) return;
      this.message = `Game mode changed to ${this.gameMode}`;
      sessionStorage.setItem('selectedGameMode', this.gameMode);
    },

    /** 🔁 Poll both game status and session players */
    async pollGameStatusAndPlayers() {
      const token = this.getAuthToken();
      if (!token) return;

      try {
        const statusRes = await fetch('/api/api/gameStatus', {
          method: 'GET',
          headers: { Authorization: `Bearer ${token}` }
        });

        console.log('[pollGameStatusAndPlayers] GameStatus status:', statusRes.status);
        const statusData = await statusRes.json().catch(() => ({}));
        console.log('[pollGameStatusAndPlayers] GameStatus data:', statusData);

        if (statusRes.ok && statusData.status === 'success') {
          const raw = statusData.Game_Status || '';
          const status = raw.charAt(0).toUpperCase() + raw.slice(1).toLowerCase();
          this.serverGameStatus = status;
          this.isGameActive = (status === 'Created' || status === 'Started');
        } else {
          console.warn('[pollGameStatusAndPlayers] Non-success response:', statusData);
          this.serverGameStatus = 'Idle';
          this.isGameActive = false;
        }

        if (this.isGameActive) {
          console.log('[pollGameStatusAndPlayers] Game active — fetching session players');
          await this.fetchSessionPlayers();
        } else {
          this.players = this.players.map(p => ({ ...p, online: false }));
        }

      } catch (err) {
        console.error('[pollGameStatusAndPlayers] Error:', err);
        this.isGameActive = false;
      }
    },

    async fetchAllUsers() {
      const token = this.getAuthToken();
      if (!token) return;

      const url = '/api/api/users';
      console.log('[fetchAllUsers] fetching from:', url);

      try {
        const res = await fetch(url, {
          method: 'GET',
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        });

        console.log('[fetchAllUsers] response status:', res.status);
        const text = await res.text();
        let data;
        try { data = JSON.parse(text); }
        catch (err) {
          console.error('[fetchAllUsers] JSON parse failed:', err, text);
          this.message = 'Invalid JSON response from backend';
          return;
        }

        const users = Array.isArray(data.data) ? data.data : [];
        console.log('[fetchAllUsers] users:', users);

        this.players = users.map((u, i) => ({
          username: u.username || `User#${i}`,
          team: u.team || '-',
          online: false
        }));

        console.log('[fetchAllUsers] players array updated:', this.players);
      } catch (err) {
        console.error('[fetchAllUsers] network or exception:', err);
        this.message = 'Network error while fetching users.';
      }
    },

    async fetchSessionPlayers() {
      const token = this.getAuthToken();
      if (!token) {
        return;
      }

      const url = '/api/api/sessionPlayers';
      try {
        const res = await fetch(url, {
          method: 'GET',
          headers: { Authorization: `Bearer ${token}` },
          credentials: 'include'
        });

        const text = await res.text();
        console.log('Raw response:', text.substring(0, 300));

        let data;
        try { data = JSON.parse(text); }
        catch (err) {
          return;
        }

        if (res.ok && data.status === 'success') {
          const sessionUsers = new Set(data.Players);
          this.players = this.players.map(p => ({
            ...p,
            online: sessionUsers.has(p.username)
          }));
        } else if (res.status === 400 && data.error) {
          console.warn(data.error);
          this.players = this.players.map(p => ({ ...p, online: false }));
          this.isGameActive = false;
          this.message = data.error;
        } else if (res.status === 401 || res.status === 403) {
          this.message = 'Session expired. Please log in again.';
          this.$router.push('/login');
        }
      } catch (err) {
        this.message = 'Network error while fetching session players.';
      }
    },

    async createGame() {
      const token = this.getAuthToken();
      if (!token) return;

      try {
        const res = await fetch('/api/api/createGame', {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ game_type: this.gameMode })
        });

        console.log('[createGame] response:', res.status);
        const data = await res.json().catch(() => ({}));

        if (res.ok && data.status === 'success') {
          console.log(data);
          this.message = data.message;
          this.isGameActive = true;
        } else {
          console.warn(data);
          this.message = data.error || 'Failed to create game.';
        }
      } catch (err) {
        console.error(err);
        this.message = 'Error creating game.';
      }
    },

    async startGame() {
      console.log('[startGame] called');
      const token = this.getAuthToken();
      if (!token) return;

      try {
        const res = await fetch('/api/api/startGame', {
          method: 'POST',
          headers: { Authorization: `Bearer ${token}` }
        });
        const data = await res.json();
        console.log('[startGame] response:', data);

        this.message = data.message || 'Game started';
      } catch (err) {
        console.error('[startGame] error:', err);
        this.message = 'Error starting game';
      }
    },

    async stopGame() {
      console.log('[stopGame] called');
      const token = this.getAuthToken();
      if (!token) return;

      try {
        const res = await fetch('/api/api/stopGame', {
          method: 'POST',
          headers: { Authorization: `Bearer ${token}` }
        });
        const data = await res.json();
        console.log('[stopGame] response:', data);

        this.message = data.message || 'Game stopped.';
        this.isGameActive = false;
        this.players = this.players.map(p => ({ ...p, online: false }));
      } catch (err) {
        console.error('[stopGame] error:', err);
        this.message = 'Error stopping game.';
      }
    },

    goToLeaderboard() {
      const routes = {
        Freefall: '/leaderboard-ffa',
        Infected: '/leaderboard-inf',
        TeamDeathmatch: '/leaderboard-tdm'
      };
      this.$router.push(routes[this.gameMode]);
    }
  },

  mounted() {
    console.log('[mounted] AdminBoard mounted');
    this.fetchAllUsers();
    this.pollGameStatusAndPlayers();
    this.pollingInterval = setInterval(this.pollGameStatusAndPlayers, 2500);
  },

  beforeUnmount() {
    if (this.pollingInterval) clearInterval(this.pollingInterval);
  }
};
</script>

<style>
.page-container {
  position: relative;
  width: 100%;
  min-height: 100vh;
  overflow-y: auto;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 2rem 0;
}
.adminboard-page {
  display: flex;
  justify-content: center;
  width: 100%;
}
</style>

<style scoped>
  .adminboard-container {
    width: 55%;
    max-width: 800px;
    background: #fff;
    border: 0.25vw solid #000;
    border-radius: 1vw;
    padding: 2.5vw 2vw;
    box-shadow: 0 0.8vw 1.5vw rgba(0, 0, 0, 0.25);
    display: flex;
    flex-direction: column;
    justify-content: center;  
    align-items: center;      
  }

  .gametype-select {
    margin-bottom: 2rem;
  }

  .top-buttons {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2vw;
  }

  .left-buttons {
    display: flex;
    flex-direction: column;
    gap: 0.7vw;
  }

  .adminboard-title {
    font-size: 1.8vw;
    font-weight: 700;
    text-align: center;
  }

  .adminboard-btn {
    background-color: #ffffff;
    border: 0.15vw solid #000;
    color: black;
    font-weight: 600;
    padding: 0.6vw 1.2vw;
    border-radius: 0.5vw;
    cursor: pointer;
    font-size: 1vw;
    transition: all 0.25s ease;
    box-shadow: 0 0.4vw 0.8vw rgba(0, 0, 0, 0.1);
  }

  .adminboard-btn:hover {
    background-color: #dac3c3;
    transform: translateY(-0.1vw);
  }

  .adminboard-btn:active { 
    background-color: #e6e6e6; 
    box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.2); 
    transform: translateY(0); 
    transition: all 0.1s ease-in; 
  } 

  .adminboard-btn:active { 
    transform: translateY(0); 
    box-shadow: 0 3px 6px rgba(0, 0, 0, 0.15); 
  }

  .gametype-select {
    text-align: center;
    margin-bottom: 1.5vw;
  }

  .adminboard-select {
    padding: 0.4vw 0.8vw;
    border: 0.1vw solid #000;
    margin-left: 0.6vw;
    border-radius: 0.3vw;
  }

  .player-table {
    width: 100%;
    border-collapse: collapse;
    margin-bottom: 2vw;
    text-align: center;
    font-size: 1vw;
  }

  .player-table th,
  .player-table td {
    border: 0.1vw solid #ddd;
    padding: 0.8vw;
  }

  .player-table thead {
    background-color: #f7f7f7;
    font-weight: 600;
  }

  .alive {
    color: green;
    font-weight: bold;
  }

  .dead {
    color: red;
    font-weight: bold;
  }

  .session-buttons {
    display: flex;
    justify-content: center;
    gap: 2vw;
    width: 100%;
  }

  .start-session,
  .end-session {
    font-weight: 600;
    font-size: 1vw;
    border-radius: 0.5vw;
    padding: 0.8vw 1.8vw;
    cursor: pointer;
    transition: all 0.25s ease;
    border: none;
    box-shadow: 0 0.4vw 0.8vw rgba(0, 0, 0, 0.15);
  }

  .start-session {
    background-color: #28a745;
    color: white;
  }

  .start-session:hover {
    background-color: #218838;
    transform: translateY(-0.2vw);
  }

  .start-session:active { 
    transform: translateY(0); 
    box-shadow: 0 3px 6px rgba(0, 0, 0, 0.15); 
  }

  .end-session {
    background-color: #dc3545;
    color: white;
  }

  .end-session:hover {
    background-color: #b02a37;
    transform: translateY(-0.2vw);
  }

  .online { 
    color: green; 
    font-weight: bold; 
  }
  
  .offline { 
    color: red; 
    font-weight: bold; 
  }

  .end-session:active { 
    transform: translateY(0); 
    box-shadow: 0 3px 6px rgba(81, 7, 7, 0.15); 
  }

  .message-box {
    margin-top: 1vh;
    padding: 1.2vh 2vw;
    font-weight: 600;
    text-align: center;
    font-size: 1vw;
  }

@media (max-width: 768px) {
  .adminboard-container {
    width: 90%;
    padding: 5vw 4vw;
    border-radius: 3vw;
  }

  .top-buttons {
    flex-direction: column;
    align-items: center;
    gap: 2vw;
  }

  .left-buttons {
    flex-direction: row;
    justify-content: center;
    gap: 3vw;
  }

  .adminboard-title {
    font-size: 5vw;
    text-align: center;
  }

  .adminboard-btn {
    font-size: 3.5vw;
    padding: 1.5vw 3vw;
    border-radius: 2vw;
  }

  .gametype-select label {
    font-size: 3.8vw;
  }

  .adminboard-select {
    font-size: 3.5vw;
    padding: 1.2vw 2vw;
  }

  .adminboard-container {
    overflow-x: auto;
  }

  .player-table {
    font-size: 3.2vw;
    border-spacing: 0;
  }

  .player-table th,
  .player-table td {
    padding: 2vw 3vw;
    border-width: 0.3vw;
  }

  .session-buttons {
    flex-direction: column;
    align-items: center;
    gap: 3vw;
  }

  .start-session,
  .end-session {
    font-size: 3.5vw;
    padding: 2vw 5vw;
    border-radius: 3vw;
  }

  .message-box {
    font-size: 3.5vw;
    margin-top: 3vw;
  }
}

@media (max-width: 480px) {
  .adminboard-container {
    margin-top: -100px;
    width: 80%;
    padding: 6vw 5vw;
  }

  .adminboard-title {
    font-size: 6vw;
  }

  .adminboard-btn,
  .start-session,
  .end-session {
    font-size: 4vw;
    padding: 2vw 6vw;
  }

  .player-table {
    font-size: 3.8vw;
  }
}
</style>