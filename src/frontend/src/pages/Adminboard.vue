<template>
  <div class="page-container">
    <div class="adminboard-page">
      <div class="adminboard-container">
        <div class="top-buttons">
          <div class="left-buttons">
            <button class="adminboard-btn" @click="$router.push('/adminedit')"> Manage players </button>
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
            <option value="FreeFall">FreeFall</option>
            <option value="Infected">Infected</option>
            <option value="Team Deathmatch">Team Deathmatch</option>
          </select>
        </div>

        <table class="player-table">
          <thead>
            <tr>
              <th>Joined Players</th>
              <th>Team</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td>Orbay</td>
              <td>Blue</td>
            </tr>
            <tr>
              <td>Berk</td>
              <td>Red</td>
            </tr>
          </tbody>
        </table>

        <div class="session-buttons">
          <button class="start-session" @click="startGameSession">Start session</button>
          <button class="end-session" @click="endGameSession">End session</button>
        </div>

        <div v-if="message" class="message-box">
          {{ message }}
        </div>
      </div>
    </div>
  </div>  
</template>

<script>
  export default {
    name: 'AdminBoard',
    data(){
      return{
      message: "",
      gameMode: 'FreeFall'
      }
    },

    methods: {
      startGameSession(){
        this.message = 'Session has started.'
        alert('Session has started.')
      },
      endGameSession(){
        this.message = 'Session has ended.'
        alert('Session has ended.')
      },
      handleGameModeChange(event) {
        const mode = event.target.value;
        this.message = `Game mode changed to ${mode}`
        alert(`Game mode has changed to ${mode}`)
      },
      goToLeaderboard() {
        if (this.gameMode === 'FreeFall') {
          this.$router.push('/leaderboard-ffa');
        } else if (this.gameMode === 'Infected') {
          this.$router.push('/leaderboard-inf');
        } else if (this.gameMode === 'Team Deathmatch') {
          this.$router.push('/leaderboard-tdm');
        }
      } 
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