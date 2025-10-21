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
            <th>Deaths</th>
            <th>Score (Hits)</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(player, index) in players" :key="player.name">
            <td>{{ index + 1 }}</td>
            <td>{{ player.name }}</td>
            <td>{{ player.deaths }}</td>
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
  name: 'LeaderboardFFA',
  data() {
    return {
      players: [
        { name: 'Berk', deaths: 3, score: 27 },
        { name: 'Orbay', deaths: 2, score: 25 },
        { name: 'Richard', deaths: 9, score: 14 },
        { name: 'Peter', deaths: 3, score: 9 },
        { name: 'Muhammed', deaths: 7, score: 8 },
        { name: 'Marciano', deaths: 2, score: 4 }
      ]
    }
  },
  methods: {
    goBack() {
      this.$router.go(-1);
    }
  }
}
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
</style>
