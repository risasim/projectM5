<template>
  <header class="header-gradient">
    <router-link to="/"> 
      <img src="@/assets/phosho-coollogo_com.png" alt="Logo" class="header-title" /> 
    </router-link>

    <div v-if="isHomePage">
      <router-link 
        v-if="!isAuthenticated" 
        to="/login" 
        class="login-button"
      >
        Log In
      </router-link>

      <button 
        v-else 
        @click="toggleLogoutSure" 
        class="logout-button"
      >
        Log Out
      </button>
    </div>

    <div 
      v-if="showLogoutSure" 
      class="sure-overlay" 
      @click.self="toggleLogoutSure"
    >
      <div class="sure-content">
        <h2>Confirm Logout</h2>
        <p>Are you sure you want to log out?</p>
        <div class="sure-buttons">
          <button class="confirm-btn" @click="confirmLogout">Yes, Log Out</button>
          <button class="cancel-btn" @click="toggleLogoutSure">Cancel</button>
        </div>
      </div>
    </div>
  </header>
</template>

<script>
export default {
  name: 'AppHeader',
  data() {
    return {
      showLogoutSure: false,
      isAuthenticated: false
    };
  },
  computed: {
    currentPath() {
      return this.$route.path;
    },
    isHomePage() {
      return this.currentPath === '/';
    }
  },
  mounted() {
    this.checkAuthStatus();
    window.addEventListener('storage', this.checkAuthStatus);
  },
  beforeUnmount() {
    window.removeEventListener('storage', this.checkAuthStatus);
  },
  methods: {
    checkAuthStatus() {
      const token = localStorage.getItem('authToken');
      this.isAuthenticated = !!token;
    },
    toggleLogoutSure() {
      this.showLogoutSure = !this.showLogoutSure;
    },
    confirmLogout() {
      localStorage.clear();
      this.isAuthenticated = false;
      this.showLogoutSure = false;
      console.log('User logged out');
      this.$router.push('/');
      location.reload();
    }
  },
  watch: {
    $route() {
      this.checkAuthStatus();
    }
  }
};
</script>

<style scoped>
.header-gradient {
  width: 100%;
  height: 90px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 40px;
  background: linear-gradient(to right, #ff5e5e, #ffc8c8, #ff4e4e);
  box-shadow: 0 9px 19px rgb(0, 0, 0);
  font-family: 'Trebuchet MS', 'Lucida Sans', Arial;
  box-sizing: border-box;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 10;
}

.header-title {
  color: white;
  width: 350px;
}

.login-button,
.logout-button {
  background: white;
  color: #000000;
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  text-decoration: none;
  border: 2px solid #000000;
  transition: 0.25s ease;
  cursor: pointer;
}

.login-button:hover,
.logout-button:hover {
  background: #dac3c3;
  color: white;
  transform: scale(1.05);
}

.sure-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 1000;
  display: flex;
  justify-content: center;
  align-items: center;
}

.sure-content {
  background: white;
  padding: 30px 40px;
  border-radius: 12px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.3);
  text-align: center;
  width: 90%;
  max-width: 400px;
}

.sure-content h2 {
  margin-bottom: 10px;
  font-size: 1.5rem;
  color: #333;
}

.sure-content p {
  margin-bottom: 20px;
  color: #555;
}

.sure-buttons {
  display: flex;
  justify-content: center;
  gap: 15px;
}

.confirm-btn {
  background-color: #ff1500;
  color: white;
  border: 2px solid #000;
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: 0.25s ease;
}

.confirm-btn:hover {
  background-color: #620d0a;
  transform: scale(1.05);
}

.cancel-btn {
  background-color: #f0f0f0;
  color: black;
  border: 2px solid #000;
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: 0.25s ease;
}

.cancel-btn:hover {
  background-color: #d1d1d1;
  transform: scale(1.05);
}

@media (max-width: 600px) {
  .header-gradient {
    justify-content: center;
    padding: 0;
    height: 70px;
  }

  .header-title {
    width: 250px;
  }

  .login-button,
  .logout-button {
    position: absolute;
    right: 10px;
    top: 20px;
    padding: 5px 12px;
    font-size: 0.9rem;
  }
}
</style>
