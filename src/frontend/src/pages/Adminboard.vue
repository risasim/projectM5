<template>
  <div class="page-container">
    <div class="adminboard-page">
      <div class="adminboard-container">
        <div class="top-buttons">
          <div class="left-buttons">
            <button class="adminboard-btn" @click="$router.push('/adminedit')">
              Manage Players
            </button>
          </div>

          <h1 class="adminboard-title">Game Session Settings</h1>

          <button class="adminboard-btn" @click="goToLeaderboard">Leaderboard</button>
        </div>

        <div class="gametype-select">
          <label for="gametype">Choose gametype:</label>
          <select
            id="gametype"
            class="adminboard-select"
            v-model="gameMode"
            @change="handleGameModeChange"
          >
            <option value="Freefall">FreeFall</option>
            <option value="Infected">Infected</option>
            <option value="TeamDeathmatch">Team Deathmatch</option>
          </select>
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
            <tr v-for="player in players" :key="player.username">
              <td>{{ player.username }}</td>
              <td>{{ player.team || '-' }}</td>
              <td :class="player.online ? 'alive' : 'dead'">{{ player.online ? 'Online' : 'Offline' }}</td>
            </tr>
            <tr v-if="players.length === 0">
              <td colspan="3">No players connected yet.</td>
            </tr>
          </tbody>
        </table>

        <div class="session-buttons">
          <button class="start-session" @click="createGame">Create Game</button>
          <button class="start-session" @click="startGameSession">Start Game</button>
          <button class="end-session" @click="endGameSession">Stop Game</button>
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
      gameMode: 'Freefall',
      players: [],
      websocket: null,
      gameStatus: 'Idle',
      statusInterval: null
    };
  },
  methods: {
    getToken() {
      const token = localStorage.getItem('authToken');
      if (!token) {
        alert('You must log in first.');
        return null;
      }
      return token;
    },

    async createGame() {
      const token = this.getToken();
      if (!token) return;
      try {
        const response = await fetch('/api/createGame', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ game_type: this.gameMode })
        });
        const data = await response.json();
        if (data.status === 'success') {
          this.message = data.message || 'Game created successfully';
          alert(this.message);
        } else {
          this.message = data.error || 'Failed to create game';
          alert(this.message);
        }
      } catch (error) {
        console.error('Create game error:', error);
        this.message = 'Network error while creating game.';
      }
    },

    async startGameSession() {
      const token = this.getToken();
      if (!token) return;
      try {
        const response = await fetch('/api/startGame', {
          method: 'POST',
          headers: { 'Authorization': `Bearer ${token}` }
        });
        const data = await response.json();
        if (data.status === 'success') {
          this.message = data.message || 'Game started successfully';
          alert(this.message);
          this.connectPisSocket();
        } else {
          this.message = data.error || 'Failed to start game';
          alert(this.message);
        }
      } catch (error) {
        console.error('Start game error:', error);
        this.message = 'Network error while starting game.';
      }
    },

    async endGameSession() {
      const token = this.getToken();
      if (!token) return;
      try {
        const response = await fetch('/api/stopGame', {
          method: 'POST',
          headers: { 'Authorization': `Bearer ${token}` }
        });
        const data = await response.json();
        if (data.status === 'success') {
          this.message = data.message || 'Game stopped successfully';
          alert(this.message);
          this.disconnectPisSocket();
        } else {
          this.message = data.error || 'Failed to stop game';
          alert(this.message);
        }
      } catch (error) {
        console.error('Stop game error:', error);
        this.message = 'Network error while stopping game.';
      }
    },

    async pollGameStatus() {
      const token = this.getToken();
      if (!token) return;
      try {
        const response = await fetch('/api/gameStatus', {
          headers: { 'Authorization': `Bearer ${token}` }
        });
        const data = await response.json();
        if (data.Game_Status) {
          this.gameStatus = data.Game_Status;
          this.message = `Game status: ${data.Game_Status}`;
        }
      } catch (error) {
        console.error('Status poll error:', error);
      }
    },

    connectPisSocket() {
      if (this.websocket) this.websocket.close();
      const websocketURL = 'ws://116.203.97.62:8080/api/wsPis';
      this.websocket = new WebSocket(websocketURL);
      this.websocket.onopen = () => console.log('Connected to Pi WebSocket');
      this.websocket.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          if (message.devices) {
            this.players = message.devices.map(device => ({
              username: device.username,
              team: device.team || '-',
              online: device.connected
            }));
          }
        } catch (error) {
          console.error('WebSocket parse error:', error);
        }
      };
      this.websocket.onerror = (error) => console.error('WebSocket error:', error);
      this.websocket.onclose = () => {
        console.log('WebSocket closed, retrying in 5s...');
        setTimeout(() => {
          if (this.gameStatus === 'Started') this.connectPisSocket();
        }, 5000);
      };
    },

    disconnectPisSocket() {
      if (this.websocket) {
        this.websocket.close();
        this.websocket = null;
      }
      this.players = [];
    },

    handleGameModeChange(event) {
      this.gameMode = event.target.value;
      this.message = `Game mode changed to ${this.gameMode}`;
      alert(`Game mode changed to ${this.gameMode}`);
    },

    goToLeaderboard() {
      const map = {
        Freefall: '/leaderboard-ffa',
        Infected: '/leaderboard-inf',
        TeamDeathmatch: '/leaderboard-tdm'
      };
      this.$router.push(map[this.gameMode]);
    }
  },

  mounted() {
    this.pollGameStatus();
    this.statusInterval = setInterval(this.pollGameStatus, 3000);
  },

  beforeUnmount() {
    clearInterval(this.statusInterval);
    this.disconnectPisSocket();
  }
};
</script>

<style>
.page-container {
  position: fixed;       
  top: 5%;
  left: 0;
  width: 100vw;
  height: 100vh;
  overflow: hidden;      
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>

<style scoped>
  .adminboard-page {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh; 
    width: 100%;
    background: none;
  }

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