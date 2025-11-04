<template>
  <div class="leaderboard-page">
    <div class="leaderboard-container">
      <div class="top-section">
        <h1 class="leaderboard-title">Leaderboard (Team Deathmatch)</h1>
      </div>
      
      <p v-if="serverGameStatus !== 'Started' && teams.length === 0" style="text-align: center; margin-bottom: 1rem; font-weight: 600;">
        Waiting for game to start...
      </p>
      
      <div v-if="teams.length === 0 && serverGameStatus === 'Started'" style="text-align: center; margin-bottom: 1rem;">
        Waiting for player data...
      </div>

      <div v-for="(team, index) in sortedTeams" :key="team.name" class="team-section">
        <h2 class="team-name">{{ index + 1 }}. {{ team.name }} — Score: {{ team.score }}</h2>

        <table class="leaderboard-table">
          <thead>
            <tr>
              <th>Player</th>
              <th>Status</th>
              <th>Team</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="member in team.sortedMembers" :key="member.username">
              <td>{{ member.username }}</td>
              <td>{{ member.deaths }}</td>
              <td>{{ member.score }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <button class="back-btn" @click="goBack">Back</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'LeaderboardTDM',
  data() {
    return {
      teams: [],
      websocket: null,
      serverGameStatus: 'Idle',
      gameStatusPolling: null
    };
  },
  computed: {
    sortedTeams() {
      return this.teams
        .slice()
        .sort((a, b) => b.score - a.score)
        .map(team => ({
          ...team,
          sortedMembers: team.members
            .slice()
            .sort((a, b) => b.score - a.score || a.deaths - b.deaths)
        }));
    }
  },
  methods: {
    goBack() {
      this.$router.go(-1);
    },

    async getGameStatus() {
        const token = localStorage.getItem('authToken');
        if (!token) {
            console.warn('[GameStatus] No token for getGameStatus. Cannot poll.');
            this.serverGameStatus = 'Inactive';
            return;
        }
        try {
            const res = await fetch('/api/gameStatus', {
                method: 'GET',
                headers: { Authorization: `Bearer ${token}` }
            });
            
            if (res.status === 401) {
                console.warn('[GameStatus] Token expired, stopping polling.');
                this.serverGameStatus = 'Inactive'; 
                if (this.gameStatusPolling) clearInterval(this.gameStatusPolling);
                this.gameStatusPolling = null;
                return;
            }

            const data = await res.json().catch(() => ({}));
            
            if (res.ok && data.status === 'success') {
                const rawStatus = data.Game_Status;
                
                if (typeof rawStatus === 'string' && rawStatus.length > 0) {
                    const lowerStatus = rawStatus.toLowerCase();
                    const newStatus = lowerStatus.charAt(0).toUpperCase() + lowerStatus.slice(1);
                    
                    const oldStatus = this.serverGameStatus; 
                    this.serverGameStatus = newStatus;

                   if (newStatus === 'Started' && oldStatus !== 'Started') {
                        console.log('[GameStatus] Game has started, connecting to WebSocket.');
                        this.connectLeaderboard();
                    } else if (newStatus !== 'Started' && oldStatus === 'Started') {
                        console.log('[GameStatus] Game has stopped, disconnecting WebSocket.');
                        if (this.websocket) {
                            this.websocket.close();
                        }
                        this.teams = [];
                    }

                } else {
                    console.warn('[GameStatus] Server response missing or invalid Game_Status:', rawStatus);
                    this.serverGameStatus = 'Idle'; 
                    if (this.websocket) this.websocket.close();
                    this.teams = [];
                }

            } else {
                console.warn('[GameStatus] Failed (non-success response):', data.error || data.message || res.statusText);
                this.serverGameStatus = 'Inactive';
                if (this.websocket) this.websocket.close();
                this.teams = [];
            }
        } catch (err) {
            console.error('[GameStatus] Poll failed (network error):', err);
            this.serverGameStatus = 'Inactive'; 
            if (this.websocket) this.websocket.close();
            this.teams = [];
        }
    },

    connectLeaderboard() {
      const token = localStorage.getItem("authToken");
      const websocketURL = `ws://116.203.97.62:8080/api/wsLeaderboard?token=${token}`;
      this.websocket = new WebSocket(websocketURL);

      this.websocket.onopen = () => {
        console.log('Connected to leaderboard WebSocket (TDM)');
      };

      this.websocket.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);

          // only to process the team deathmatch
          if (message.game_type?.toLowerCase() === 'teamdeathmatch' && Array.isArray(message.teams)) {
            this.teams = message.teams.map(team => ({
              name: team.name,
              status: player.status || 'unknown',
              members: team.members?.map(m => ({
                username: m.username,
                deaths: m.deaths || 0,
                score: m.score || 0
              })) || []
            }));
          }
        } catch (error) {
          console.error('WebSocket parse error:', error);
        }
      };

      this.websocket.onerror = (error) => {
        console.error('WebSocket error:', error);
      };

      this.websocket.onclose = () => {
        console.log('WebSocket closed.');
        if (this.serverGameStatus === 'Started') {
          console.log('Game is active. Reconnecting in 5s...');
          setTimeout(this.connectLeaderboard, 5000);
        } else {
          console.log('Game is not active. Not reconnecting.');
        }
      };
    }
  },
  mounted() {
    this.getGameStatus();
    this.gameStatusPolling = setInterval(this.getGameStatus, 2500);
  },
  beforeUnmount() {
    if (this.gameStatusPolling) clearInterval(this.gameStatusPolling);
    if (this.websocket) this.websocket.close();
  }
};
</script>

<style scoped>
.leaderboard-page {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh; 
  width: 100%;
  background: none;
}

.leaderboard-container {
  width: 55%;
  max-width: 800px;
  background: #fff;
  border: 0.25vw solid #000;
  border-radius: 1vw;
  padding: 2.5vw 2vw;
  box-shadow: 0 0.8vw 1.5vw rgba(0, 0, 0, 0.25);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.top-section {
  text-align: center;
  margin-bottom: 2vw;
}

.leaderboard-title {
  font-size: 1.8vw;
  font-weight: 700;
  text-align: center;
}

.leaderboard-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 2vw;
  text-align: center;
  font-size: 1vw;
}

.leaderboard-table th,
.leaderboard-table td {
  border: 0.1vw solid #ddd;
  padding: 0.8vw;
}

.leaderboard-table thead {
  background-color: #f7f7f7;
  font-weight: 600;
}

.leaderboard-table tbody tr:nth-child(even) {
  background-color: #f2f2f2;
}

.back-btn {
  font-weight: 600;
  font-size: 1vw;
  border-radius: 0.5vw;
  padding: 0.8vw 1.8vw;
  cursor: pointer;
  transition: all 0.25s ease;
  border: none;
  box-shadow: 0 0.4vw 0.8vw rgba(0, 0, 0, 0.15);
  background-color: #ff1500;
  color: white;
  border: 4px solid #000000;
}

.back-btn:hover {
  background-color: #620d0a;
  transform: translateY(-0.2vw);
}

.back-btn:active { 
  transform: translateY(0); 
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.15); 
}

@media (max-width: 768px) {
  .leaderboard-container {
    width: 90%;
    padding: 5vw 3vw;
    border-radius: 3vw;
  }

  .leaderboard-title {
    font-size: 5vw;
    margin-bottom: 3vw;
  }

  .leaderboard-table {
    font-size: 3.8vw;
    border-spacing: 0;
  }

  .leaderboard-table th,
  .leaderboard-table td {
    padding: 3vw 1vw;
    border-width: 0.3vw;
  }

  .leaderboard-table thead {
    background-color: #f7f7f7;
  }

  .leaderboard-container {
    overflow-x: auto;
  }

  .back-btn {
    font-size: 3.5vw;
    padding: 2vw 5vw;
    border-radius: 3vw;
  }
}

@media (max-width: 480px) {
  .leaderboard-container {
    width: 80%;
    margin-top: -100px;
    padding: 6vw 4vw;
  }

  .leaderboard-title {
    font-size: 6vw;
  }

  .leaderboard-table {
    font-size: 4vw;
  }

  .back-btn {
    font-size: 4vw;
    padding: 2.5vw 6vw;
  }
}
</style>
