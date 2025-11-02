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
        @click="toggleSettingsOverlay" 
        class="settings-button"
      >
        Settings
      </button>
    </div>

    <div 
      v-if="showSettingsOverlay" 
      class="sure-overlay" 
      @click.self="toggleSettingsOverlay"
    >
      <div class="sure-content settings-content">
        <h2>Account Settings</h2>
        <p>Manage your account options below.</p>
        <div class="sure-buttons settings-buttons">
          
          <button 
            class="action-btn logout-btn" 
            @click="promptLogoutConfirmation" 
          >
            Log Out
          </button>
          
          <button 
            class="action-btn delete-btn" 
            @click="promptDeleteConfirmation" >
            Delete Account
          </button>
          
          <button class="cancel-btn close-btn" @click="toggleSettingsOverlay">Close</button>
        </div>
      </div>
    </div>
    
    <div 
      v-if="showLogoutConfirmation" 
      class="sure-overlay" 
      @click.self="cancelLogoutConfirmation"
    >
      <div class="sure-content logout-confirm-content">
        <h2>Confirm Logout</h2>
        <p>Are you sure you want to log out?</p>
        <div class="sure-buttons confirm-buttons">
          <button class="confirm-btn" @click="confirmLogout">Yes, Log Out</button>
          <button class="cancel-btn" @click="cancelLogoutConfirmation">Cancel</button>
        </div>
      </div>
    </div>

    <div 
      v-if="showDeleteConfirmation" 
      class="sure-overlay" 
      @click.self="cancelDeleteConfirmation"
    >
      <div class="sure-content delete-confirm-content">
        <h2>⚠️ Permanent Account Deletion</h2>
        <p>This action is **irreversible**. Are you absolutely sure you want to permanently delete your account?</p>
        <div class="sure-buttons confirm-buttons">
          <button class="confirm-btn delete-confirm-btn" @click="confirmDeleteAccount">Yes, Delete Account</button>
          <button class="cancel-btn" @click="cancelDeleteConfirmation">Cancel</button>
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
      showSettingsOverlay: false, 
      showLogoutConfirmation: false, 
      showDeleteConfirmation: false,
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
    toggleSettingsOverlay() {
      this.showSettingsOverlay = !this.showSettingsOverlay;
      if (!this.showSettingsOverlay) {
        this.showLogoutConfirmation = false;
        this.showDeleteConfirmation = false;
      }
    },
    
    promptLogoutConfirmation() {
        this.showSettingsOverlay = false;
        this.showLogoutConfirmation = true;
    },
    cancelLogoutConfirmation() {
        this.showLogoutConfirmation = false;
        this.showSettingsOverlay = true; 
    },
    confirmLogout() {
      localStorage.clear();
      this.isAuthenticated = false;
      this.showSettingsOverlay = false; 
      this.showLogoutConfirmation = false; 
      console.log('User logged out');
      this.$router.push('/');
      location.reload();
    },

    promptDeleteConfirmation() {
        this.showSettingsOverlay = false;
        this.showDeleteConfirmation = true;
    },
    cancelDeleteConfirmation() {
        this.showDeleteConfirmation = false;
        this.showSettingsOverlay = true; 
    },

    async confirmDeleteAccount() {
      const token = localStorage.getItem('authToken');
      if (!token) {
        alert('You are not logged in.');
        this.showDeleteConfirmation = false;
        this.showSettingsOverlay = false;
        this.$router.push('/login');
        return;
      }

      try {
        const res = await fetch('/api/user', { 
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${token}`, 
          },
        });

        if (res.ok) {
          alert('Account successfully deleted. Goodbye!');
          localStorage.clear();
          this.isAuthenticated = false;
          this.showSettingsOverlay = false;
          this.showDeleteConfirmation = false;
          this.$router.push('/');
          location.reload();
        } else {
          const data = await res.json().catch(() => ({}));
          const errorMessage = data.error || data.message || `Failed to delete account (Status: ${res.status}).`;
          alert(errorMessage);
          this.showDeleteConfirmation = false;
          this.showSettingsOverlay = true;
        }
      } catch (error) {
        console.error('Account deletion failed:', error);
        alert('An error occurred while trying to delete the account.');
        this.showDeleteConfirmation = false;
        this.showSettingsOverlay = true;
      }
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
.settings-button {
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
.settings-button:hover {
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

.settings-buttons {
    flex-direction: column; 
    gap: 10px !important;
}

.action-btn {
  width: 100%;
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: 0.25s ease;
  border: 2px solid #000;
}

.logout-btn,
.delete-btn {
    background-color: #ff1500; 
    color: white; 
}
.logout-btn:hover,
.delete-btn:hover {
    background-color: #620d0a; 
    transform: scale(1.02);
}

.close-btn {
    width: 100%; 
}
.confirm-buttons {
    flex-direction: row;
    justify-content: center;
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
  width: auto; 
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
  width: auto; 
}

.cancel-btn:hover {
  background-color: #d1d1d1;
  transform: scale(1.05);
}

.delete-confirm-btn {
    background-color: #b30000;
}
.delete-confirm-btn:hover {
    background-color: #800000;
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
  .settings-button {
    position: absolute;
    right: 10px;
    top: 20px;
    padding: 5px 12px;
    font-size: 0.9rem;
  }
  
  .confirm-buttons {
    flex-direction: column;
  }
  
  .confirm-btn,
  .cancel-btn {
      width: 100%;
      margin: 5px 0;
  }
}

@media (max-width: 450px) {
    .sure-content {
        padding: 20px;
    }
}
</style>
