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
</style>
