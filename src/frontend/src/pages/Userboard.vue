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

        <!-- SFX section -->
        <div class="sfx-section">
          <label for="deathSfx" class="sfx-label">Custom Death SFX:</label>

          <input
            id="deathSfx"
            ref="fileInput"
            type="file"
            accept=".mp3"
            class="sfx-input"
            @change="onFileSelected"
          />

          <div class="sfx-actions">
            <button
              class="upload-btn"
              @click="uploadSelected"
              :disabled="!selectedFile || uploading"
            >
              <span v-if="uploading">Uploading…</span>
              <span v-else>Upload</span>
            </button>

            <button
              class="play-btn"
              @click="playSound"
              :disabled="playing || (!hasSound && !selectedFile)"
            >
              <span v-if="playing">Playing…</span>
              <span v-else>Play sound</span>
            </button>
          </div>
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
      selectedFile: null,
      uploading: false,
      playing: false,
      hasSound: false,
      audioObjectUrl: null,
      previewUrl: null
    };
  },
  computed: {
    sessionStatusText() {
      if (this.sessionStatus === 'active') return 'Active';
      if (this.sessionStatus === 'waiting') return 'Waiting for players';
      return 'Inactive';
    }
  },
  mounted() {
    this.checkHasSound();
  },
  beforeUnmount() {
    if (this.audioObjectUrl) URL.revokeObjectURL(this.audioObjectUrl);
    if (this.previewUrl) URL.revokeObjectURL(this.previewUrl);
  },
  methods: {
    onFileSelected(e) {
      this.selectedFile = e.target.files[0] || null;
      console.log('Selected file:', this.selectedFile && this.selectedFile.name);
      if (this.previewUrl) {
        URL.revokeObjectURL(this.previewUrl);
        this.previewUrl = null;
      }
      if (this.selectedFile) {
        this.previewUrl = URL.createObjectURL(this.selectedFile);
      }
    },  async fetchSoundFromServer() {
  const token = localStorage.getItem('authToken');
  if (!token) return;

  try {
    const res = await fetch('/api/api/sound', {
      method: 'GET',
      headers: { Authorization: `Bearer ${token}` }
    });

    if (!res.ok) {
      console.warn('No uploaded sound found.');
      this.hasSound = false;
      return;
    }

    const blob = await res.blob();
    if (this.audioObjectUrl) URL.revokeObjectURL(this.audioObjectUrl);
    this.audioObjectUrl = URL.createObjectURL(blob);
    this.hasSound = true;
    console.log('Fetched uploaded sound successfully!');
  } catch (err) {
    console.error('Failed to fetch sound:', err);
    this.hasSound = false;
  }
},

    async uploadSelected() {
      if (!this.selectedFile) {
        alert('Choose a file first.');
        return;
      }

      const file = this.selectedFile;
      const token = localStorage.getItem('authToken');
      if (!token) {
        alert('You must log in again.');
        return;
      }

      const username = encodeURIComponent(this.username);
      const formData = new FormData();
      formData.append('sound', file);

      this.uploading = true;
      try {
        const res = await fetch(`/api/api/uploadSound`, {
          method: 'POST',
          headers: { Authorization: `Bearer ${token}` },
          body: formData
        });

        let parsed = null;
        try { parsed = await res.json(); } catch (e) {}

        if (res.ok && parsed && parsed.status === 'success') {
          alert('Audio uploaded successfully!');
          this.hasSound = true;
          if (this.$refs.fileInput) this.$refs.fileInput.value = '';
          this.selectedFile = null;
          if (this.previewUrl) {
            URL.revokeObjectURL(this.previewUrl);
            this.previewUrl = null;
          }
          await this.fetchSoundFromServer();
        } else {
          const msg = (parsed && (parsed.error || parsed.message)) || `Upload failed (${res.status})`;
          alert(msg);
        }
      } catch (err) {
        console.error('Upload failed', err);
        alert('Failed to upload sound. See console for details.');
      } finally {
        this.uploading = false;
      }
    },

    async playSound() {
      if (this.playing) return;
      this.playing = true;

      try {
        // play selected local file (preview)
        if (this.selectedFile && this.previewUrl) {
          console.log('Playing local preview:', this.selectedFile.name);
          const audio = new Audio(this.previewUrl);
          audio.onended = () => { this.playing = false; };
          audio.onerror = () => {
            this.playing = false;
            alert('Failed to play preview sound.');
          };
          await audio.play();
          return;
        }

        // otherwise, play uploaded sound from server
        const token = localStorage.getItem('authToken');
        if (!token) {
          alert('You must log in again.');
          this.playing = false;
          return;
        }

        const res = await fetch('/api/api/sound', {
          method: 'GET',
          headers: { Authorization: `Bearer ${token}` }
        });

        if (!res.ok) {
          if (res.status === 404) {
            this.hasSound = false;
            alert('No uploaded sound found.');
            return;
          }
          throw new Error(`Failed to fetch sound: ${res.status}`);
        }

        const blob = await res.blob();
        const url = URL.createObjectURL(blob);
        this.audioObjectUrl = url;

        const audio = new Audio(url);
        audio.onended = () => { this.playing = false; };
        audio.onerror = () => {
          this.playing = false;
          alert('Failed to play uploaded sound.');
        };
        await audio.play();
      } catch (err) {
        console.error('playSound error:', err);
        alert('Failed to play sound.');
        this.playing = false;
      }
    },

    async checkHasSound() {
      const token = localStorage.getItem('authToken');
      if (!token) { this.hasSound = false; return; }

      try {
        const res = await fetch('/api/api/sound', {
          method: 'GET',
          headers: { Authorization: `Bearer ${token}` }
        });
        this.hasSound = res.ok;
      } catch (err) {
        console.warn('checkHasSound failed', err);
        this.hasSound = false;
      }
    },

    async enterSession() {
      const token = localStorage.getItem('authToken');
      if (!token) {
        alert('You must log in first.');
        return;
      }

      try {
        const res = await fetch('/api/joinGame', {
          method: 'POST',
          headers: { Authorization: `Bearer ${token}` }
        });
        const data = await res.json().catch(() => ({}));
        if (res.ok && data.status === 'success') {
          alert('Joined game successfully!');
          this.sessionStatus = 'active';
        } else {
          const msg = data.error || data.message || `Failed to join (${res.status})`;
          alert(msg);
        }
      } catch (err) {
        console.error('Join game failed:', err);
        alert('Failed to join game.');
      }
    }
  }
};
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

/* keep your original styling; included minimal extras for sfx area */
.sfx-section {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  margin-bottom: 1.5vw;
  width: 100%;
}

.sfx-label { margin-right: 0.5rem; }

.sfx-actions {
  display: inline-flex;
  gap: 0.6rem;
  margin-left: 0.5rem;
}

.upload-btn,
.play-btn {
  padding: 0.5rem 0.8rem;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  font-weight: 600;
}

.upload-btn {
  background: #ff1111; color: white; border: 2px solid #000;
}
.play-btn {
  background: #ffffff; color: #000; border: 2px solid #000;
}

.upload-btn[disabled],
.play-btn[disabled] {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>


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
