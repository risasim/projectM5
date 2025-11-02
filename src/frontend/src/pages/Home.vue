<template>
  <div class="home-page">
    <div class="home-content">
      <div class="home-center">
        <img src="@/assets/phosho-coollogo_com-blackbg.png" alt="Logo" class="home-logo" />

        <div class="button-group">
          <button class="userboard-btn" @click="goToBoard">
            {{ isAdmin ? "Adminboard" : "Userboard" }}
          </button>

          <button class="leaderboard-btn" @click="goToLeaderboard">
            Leaderboard
          </button>
        </div>

        <div class="info-section">
          <div class="info-card1">
            <h2 class="info-title">Gamemodes</h2>
            <p class="info-content">
              We have several gamemodes you can play in such as Free-For-All, Infected and Team Deathmatch.
            </p>
          </div>

          <div class="info-card2">
            <h1 class="info-title">Welcome to PhoSho</h1>
            <p class="info-content">
              PhoSho is a laser-tag game brought to you by Group 29! It works with infrared lasers so it is completely safe to play with.
              Shoot those photons and gain those points on your profile! Hope you have fun!
            </p>
          </div>

          <div class="info-card3">
            <h2 class="info-title">Raspberry Pi</h2>
            <p class="info-content">
              We have a Raspberry Pi connected to the sensor on your vest and to the infrared gun.
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HomePage',
  data() {
    return {
      isAdmin: false, 
      selectedGameMode: sessionStorage.getItem("selectedGameMode") || "FreeFall" //default FFA
    };
  },
  mounted() {
    const token = sessionStorage.getItem("authToken");
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        this.isAdmin = payload.username === "admin" || payload.role === "admin";
      } catch {
        console.warn("Invalid token format");
      }
    }
  },
  methods: {
    goToBoard() {
      if (this.isAdmin) {
        this.$router.push('/adminboard');
      } else {
        this.$router.push('/userboard');
      }
    },
    goToLeaderboard() {
      // will be replaced with live socket logic, placeholder
      const gameMode = this.selectedGameMode;
      if (gameMode === 'FreeFall') {
        this.$router.push('/leaderboard-ffa');
      } else if (gameMode === 'Infected') {
        this.$router.push('/leaderboard-inf');
      } else if (gameMode === 'Team Deathmatch') {
        this.$router.push('/leaderboard-tdm');
      } else {
        this.$router.push('/leaderboard');
      }
    }
  }
};
</script>


<style scoped>
.home-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: url('@/assets/red-bg.png') center/cover no-repeat;
  color: #111;
  overflow-x: hidden;
}

.home-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  width: 100%;
  padding: 2rem 1rem;
}

.home-logo {
  display: block;
  max-width: 380px;
  width: 80%;
  height: auto;
  margin: 0 auto 1.5rem auto;
}

.userboard-btn,
.leaderboard-btn {
  background: #ffffff;
  color: #000;
  font-family: 'Trebuchet MS', Arial, sans-serif;
  font-weight: 600;
  border: 3px solid #000;
  border-radius: 25px;
  padding: 0.7rem 1.5rem;
  cursor: pointer;
  font-size: 1.2rem;
  transition: 0.3s;
  margin: 0.5rem;
  box-shadow: 2px 3px 8px rgba(0, 0, 0, 0.2);
}

.userboard-btn:hover,
.leaderboard-btn:hover {
  background: #dac3c3;
  transform: scale(1.05);
}

.info-section {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 20px;
  margin-top: 30px;
}

.info-card1,
.info-card2,
.info-card3 {
  background: #fff;
  border: 3px solid #000;
  border-radius: 14px;
  width: 280px;
  padding: 25px;
  box-shadow: 0 8px 18px rgba(0, 0, 0, 0.15);
  transition: 0.3s ease;
}

.info-card1:hover,
.info-card2:hover,
.info-card3:hover {
  transform: translateY(-5px);
}

.info-title {
  color: #b30000;
  font-size: 1.4rem;
  font-weight: 700;
  font-family: 'Trebuchet MS', Arial, sans-serif;
  margin-bottom: 10px;
}

.info-content {
  color: #333;
  font-size: 1rem;
  font-family: 'Trebuchet MS', Arial, sans-serif;
  line-height: 1.5;
}

@media (max-width: 768px) {
  .home-logo {
    max-width: 500px;
  }

  .userboard-btn,
  .leaderboard-btn {
    width: 50%;
    font-size: 1rem;
    padding: 0.6rem;
  }

  .info-section {
    flex-direction: column;
    align-items: center;
    gap: 15px;
  }

  .info-card1,
  .info-card2,
  .info-card3 {
    width: 90%;
  }
}

@media (max-width: 480px) {
  .home-logo {
    max-width: 450px;
  }

  .userboard-btn,
  .leaderboard-btn {
    font-size: 0.9rem;
    padding: 0.5rem 0.8rem;
  }
 
  .info-card1,
  .info-card2,
  .info-card3 {
    width: 80%;
    padding: 15px;
  }

  .info-title {
    font-size: 1rem;
  }

  .info-content {
    font-size: 0.9rem;
  }
}
</style>
