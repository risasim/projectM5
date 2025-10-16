<template>
  <div class="adminboard-page">
    <div class="adminboard-container">
      <div class="top-buttons">
        <div class="left-buttons">
          <button class="adminboard-btn">Add players</button>
          <button class="adminboard-btn">Manage players</button>
        </div>

        <h1 class="adminboard-title">Game Session Settings</h1>

        <router-link to="/leaderboard"><button class="adminboard-btn">Leaderboard</button></router-link>
      </div>

      <div class="gametype-select">
        <label for="gametype">Choose gametype:</label>
        <select id="gametype" class="adminboard-select">
          <option value="v">FreeFall</option>
          <option value="t">Infected</option>
          <option value="t">Team Deathmatch</option>

        </select>
      </div>

      <table class="player-table">
        <thead>
          <tr>
            <th>Players</th>
            <th>Team</th>
            <th>Hits</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>Orbay</td>
            <td>Blue</td>
            <td>31</td>
            <td class="alive">Alive</td>
          </tr>
          <tr>
            <td>Berk</td>
            <td>Red</td>
            <td>25</td>
            <td class="dead">Dead</td>
          </tr>
        </tbody>
      </table>

      <div class="session-buttons">
        <button class="start-session">Start session</button>
        <button class="end-session">End session</button>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'AdminBoard'
  }
</script>

<style scoped>.adminboard-page {
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
    background-color: #f2f2f2;
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
    border: 4px solid #000000;
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
</style>
