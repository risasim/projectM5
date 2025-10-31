<template>
  <div class="leaderboard-selection-page">
    <h1 class="info-text">
      This page is temporary until we figure out sessions.<br />
      When the server starts, the leaderboard button will direct you to the 
      respective leaderboard of the session.
    </h1>

    <div class="button-container">
      <router-link to="/leaderboard-ffa">
        <button class="userboard-btn">Free For All</button>
      </router-link>
      <router-link to="/leaderboard-tdm">
        <button class="userboard-btn">Team Deathmatch</button>
      </router-link>
      <router-link to="/leaderboard-inf">
        <button class="userboard-btn">Infected</button>
      </router-link>
    </div>
  </div>
</template>

<style scoped>


.info-text {
  max-width: 700px;
  font-size: 1.5rem;
  font-weight: 600;
  color: #ffffff;
  margin-bottom: 3rem;
  line-height: 1.5;
}

.button-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.userboard-btn {
  background: #ffffff;
  color: #000000;
  font-family: 'Trebuchet MS', Arial, sans-serif;
  font-weight: 600;
  border: 4px solid #000000;
  border-radius: 20px;
  padding: 12px 28px;
  cursor: pointer;
  font-size: 1.5rem;
  transition: 0.3s;
  width: 250px;
}

.userboard-btn:hover {
  background-color: #e3e3e3;
  transform: translateY(-3px);
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
