<template>
  <div class="page-container">
    <div class="userboard-page">
      <div class="userboard-container">
        <div class="floating-userboard-title">
          <h1>Userboard</h1>
        </div>

        <div class="userboard-top">
          <div class="user-info">
            <h2 class="username">Username: <span class="value">{{ username }}</span></h2>
            <p class="team">Team: <span class="value team">{{ team }}</span></p>
          </div>

          <router-link to="/leaderboard">
            <button class="leaderboard-btn">Leaderboard</button>
          </router-link>
        </div>

        <div class="stats-section">
          <p>Your Total Victories: <span class="value">{{ victories }}</span></p>
          <p>Total Deaths: <span class="value">{{ deaths }}</span></p>
        </div>

        <div class="sfx-section">
          <label for="deathSfx" class="sfx-label">Custom Death SFX:</label>
          <input id="deathSfx" type="file" accept=".mp3, .ogg, .wav" class="sfx-input" @change="handleFileUpload" />
        </div>

        <div class="session-status">
          <p>Session status: <span :class="['status', sessionStatus]">{{ sessionStatusText }}</span></p>
        </div>

        <button class="enter-session-btn" @click="enterSession">Enter current game session</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'UserBoard',
  data() {
    return {
      username: localStorage.getItem('username') || 'Unknown',
      team: 'Unknown',
      victories: 0,
      deaths: 0,
      sessionStatus: 'waiting',
    };
  },
  computed: {
    sessionStatusText() {
      if (this.sessionStatus === 'active') return 'Active';
      if (this.sessionStatus === 'waiting') return 'Waiting for players';
      return 'Inactive';
    }
  },
  methods: {
    // upload sound to backend
    async handleFileUpload(event) {
      const file = event.target.files[0];
      if (!file) return;

      const maxSizeMB = 2;
      if (file.size > maxSizeMB * 1024 * 1024) {
        alert('File size must be less than 2 MB.');
        event.target.value = '';
        return;
      }

      try {
        const audioURL = URL.createObjectURL(file);
        const audio = new Audio(audioURL);
        await new Promise((resolve, reject) => {
          audio.onloadedmetadata = () => {
            if (audio.duration > 5) {
              alert('Audio must be less than 5 seconds.');
              event.target.value = '';
              URL.revokeObjectURL(audioURL);
              reject();
            } else resolve();
          };
        });
        URL.revokeObjectURL(audioURL);

        // upload to backend
        const token = localStorage.getItem('token');
        if (!token) {
          alert('You must log in again.');
          return;
        }
        const formData = new FormData();
        formData.append('sound', file);
        const response = await fetch(`/api/uploadSound?username=${encodeURIComponent(this.username)}`, {
          method: 'POST',
          headers: { Authorization: `Bearer ${token}` },
          body: formData
        });
        const data = await response.json();
        if (data.status === 'success') alert('Sound uploaded successfully!');
        else alert(data.error || 'Upload failed.');


        alert('Audio uploaded.');
      } catch (err) {
        console.error('Audio validation failed:', err);
      }
    },

    // join the game session
    async enterSession() {
      const token = localStorage.getItem('token');
      if (!token) {
        alert('You must log in first.');
        return;
      }

      try {
        const response = await fetch(`/api/joinGame?username=${encodeURIComponent(this.username)}`, {
          method: 'POST',
          headers: { Authorization: `Bearer ${token}` }
        });
        const data = await response.json();
        if (data.status === 'success') {
          alert('Joined game successfully!');
          this.sessionStatus = 'active';
        } else {
          alert(data.error || 'Failed to join game.');
        }
      } catch (err) {
        console.error('Join game failed:', err);
      }
    }
  }
}
</script>

<style>
.page-container {
  position: fixed;       
  width: 100vw;
  height: 100vh;
  overflow: hidden;      
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>

<style scoped>
.userboard-page {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh; 
  width: 100%;
  background: none;
}

.userboard-container {
  position: relative;
  top: 3%;
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

.floating-userboard-title {
  position: absolute;
  top: -1.8vw;
  left: 50%;
  transform: translateX(-50%);
  background: #fff;
  border: 0.25vw solid #000;
  border-radius: 0.5vw;
  padding: 0.4vw 1.2vw;
  box-shadow: 0 0.4vw 1vw rgba(0, 0, 0, 0.25);
}

.floating-userboard-title h1 {
  margin: 0;
  font-size: 1.6vw;
  font-weight: 700;
  text-align: center;
}

.userboard-top {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2vw;
}

.user-info {
  text-align: left;
}

.username, .team {
  font-size: 1.5vw;
  margin: 0.2vw 0;
}

.userboard-title {
  font-size: 1.8vw;
  font-weight: 700;
  text-align: center;
}

.leaderboard-btn {
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

.leaderboard-btn:hover {
  background-color: #f2f2f2;
  transform: translateY(-0.1vw);
}

.stats-section {
  width: 100%;
  text-align: left;
  margin-bottom: 1.5vw;
  font-size: 1.3vw;
}

.stats-section .value {
  font-weight: 600;
}

.sfx-section {
  display: flex;
  align-items: center;
  margin-bottom: 1.5vw;
  width: 100%;
}

.sfx-label {
  font-size: 1.1vw;
  margin-right: 1vw;
}

.sfx-input {
  border: 0.15vw solid #000;
  padding: 0.4vw;
  border-radius: 0.3vw;
  cursor: pointer;
}

.session-status {
  margin-bottom: 1.5vw;
  font-size: 1.1vw;
}

.status {
  font-weight: 600;
}

.status.active {
  color: green;
}

.status.waiting {
  color: orange;
}

.status.inactive {
  color: red;
}

.value.team{
  color: red;
  font: bolder;
}

.enter-session-btn {
  font-weight: 600;
  font-size: 1vw;
  border-radius: 0.5vw;
  padding: 0.8vw 1.8vw;
  cursor: pointer;
  transition: all 0.25s ease;
  border: none;
  box-shadow: 0 0.4vw 0.8vw rgba(0, 0, 0, 0.15);
  background-color: #28a745;
  color: white;
  border: 4px solid #000000;
}

.enter-session-btn:hover {
  background-color: #218838;
  transform: translateY(-0.2vw);
}

.enter-session-btn:active { 
  transform: translateY(0); 
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.15); 
}

@media (max-width: 768px) {
  .userboard-container {
    width: 90%;
    padding: 5vw 4vw;
    border-width: 0.5vw;
    border-radius: 3vw;
    top: 5%;
  }

  .floating-userboard-title {
    top: -4vw;
    padding: 1vw 3vw;
    border-radius: 2vw;
  }

  .floating-userboard-title h1 {
    font-size: 5vw;
  }

  .userboard-top {
    flex-direction: column;
    align-items: center;
    gap: 3vw;
  }

  .username,
  .team {
    font-size: 4vw;
    text-align: center;
  }

  .leaderboard-btn {
    font-size: 3.5vw;
    padding: 1vw 3vw;
    border-radius: 2vw;
  }

  .stats-section {
    font-size: 3.5vw;
    text-align: center;
  }

  .sfx-section {
    flex-direction: column;
    align-items: center;
    gap: 2vw;
  }

  .sfx-label {
    font-size: 3.5vw;
    margin: 0;
  }

  .sfx-input {
    width: 80%;
    font-size: 3.3vw;
    padding: 2vw;
    border-radius: 2vw;
  }

  .session-status {
    font-size: 3.5vw;
    text-align: center;
  }

  .enter-session-btn {
    font-size: 3.8vw;
    padding: 2vw 5vw;
    border-radius: 3vw;
  }
}

@media (max-width: 480px) {
  .userboard-container {
    margin-top: -100px;
    width: 80%;
    padding: 6vw 5vw;
  }

  .floating-userboard-title h1 {
    font-size: 6vw;
  }

  .enter-session-btn {
    font-size: 4vw;
    padding: 2.5vw 6vw;
  }
}

</style>
