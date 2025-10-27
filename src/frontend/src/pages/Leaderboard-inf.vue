<template>
  <div class="leaderboard-page">
    <div class="leaderboard-container">
      <div class="top-section">
        <h1 class="leaderboard-title">Leaderboard (Infected)</h1>
      </div>

      <table class="leaderboard-table">
        <thead>
          <tr>
            <th>Rank</th>
            <th>Player</th>
            <th>Status</th>
            <th>Score</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(player, index) in players" :key="player.username">
            <td>{{ index + 1 }}</td>
            <td>{{ player.username }}</td>
            <td :style="{ color: player.infected ? 'red' : 'blue' }">
              {{ player.infected ? 'Infected' : 'Survivor' }}
            </td>
            <td>{{ player.score }}</td>
          </tr>
        </tbody>
      </table>

      <div>
        <button class="back-btn" @click="goBack">Back</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'LeaderboardINF',
  data() {
    return {
      players: [],
      ws: null
    };
  },
  methods: {
    goBack() {
      this.$router.go(-1);
    },
    connectLeaderboard() {

      // connect to websocket
      const wsUrl = `ws://116.203.97.62:8080/api/wsLeaderboard`;
      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        console.log('Connected to leaderboard WebSocket (INFECTED)');
      };

      this.ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);

          // checks game type
          if (message.game_type && message.game_type.toLowerCase() === 'infected') {
            const infectedNames = message.data?.infected?.map(p => p.username) || [];
            const playerList = message.players || [];

            // combines survivors and infected
            this.players = playerList.map((p, i) => ({
              username: p.username,
              infected: infectedNames.includes(p.username),
              score: p.score || Math.floor(Math.random() * 30)
            }));

            // sorting, survivors first
            this.players.sort((a, b) => Number(a.infected) - Number(b.infected));
          }
        } catch (err) {
          console.error('WS parse error:', err);
        }
      };

      this.ws.onerror = (err) => {
        console.error('WebSocket error:', err);
      };

      this.ws.onclose = () => {
        console.log('WebSocket closed. Reconnecting in 5s...');
        setTimeout(() => this.connectLeaderboard(), 5000);
      };
    }
  },
  mounted() {
    this.connectLeaderboard();
  },
  beforeUnmount() {
    if (this.ws) this.ws.close();
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

