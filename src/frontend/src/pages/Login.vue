<template>
  <div class="page-container">
    <div class="login-page">
      <div class="login-box">
        <h1 class="login-title">Login</h1>
        <form class="login-form">


          <input type="text" placeholder="Username" class="login-input" />


          <input type="password" placeholder="Password" class="login-input" />


          <router-link to="/userboard"><button type="submit" class="login-button">Log In</button></router-link>
          <router-link to="/adminboard"><button type="submit" class="login-button">Log In (admin)</button></router-link>
        </form>
      </div>
    </div>
  </div>  
</template>

<script>
export default {
  name: 'AppLogin',
  methods: {
    loginUser() {
      this.errorMessage = '';
      this.successMessage = '';
      this.isLoading = true;
      // SECURITY
      const usernameResult = InputSanitizer.sanitizeUsername(this.username);
      if (!usernameResult.valid) {
        this.errorMessage = usernameResult.error;
        this.isLoading = false;
        return;
      }

      if (!this.password || this.password.trim().length === 0) {
        this.errorMessage = 'Password cannot be empty';
        this.isLoading = false;
        return;
      }

      if (this.password.length > 50) {
        this.errorMessage = 'Password too long';
        this.isLoading = false;
        return;
      }

      if (InputSanitizer.containsCommandInjection(this.password)) {
        this.errorMessage = 'Password contains dangerous characters';
        this.isLoading = false;
        return;
      }
    }
  }
};
</script>

<style scoped>
.page-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  display: flex;
  justify-content: center;
  align-items: center;
}

.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  width: 100%;
}

.login-box {
  background: white;
  padding: 2rem;
  border-radius: 20px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
  width: 360px;
  max-width: 90%;              
  text-align: center;
  border: 4px solid #000;
}

.login-title {
  margin-bottom: 1.5rem;
  font-size: 2rem;
  font-weight: 700;
  color: #ff4b4b;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.login-input {
  padding: 0.75rem;
  border-radius: 8px;
  border: 1px solid #ccc;
  font-size: 1rem;
  outline: none;
  transition: 0.2s ease;
}

.login-input:focus {
  border-color: #ff4b4b;
  box-shadow: 0 0 6px rgba(255, 75, 75, 0.4);
}

.login-button {
  background: #ff1111;
  color: white;
  font-size: 1.1rem;
  font-weight: 600;
  padding: 0.8rem;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 4px solid #000;
  width: 100%;                
}

.login-button:hover {
  background: #ff6666;
  transform: scale(1.03);
}

@media (max-width: 600px) {
  .login-box {
    padding: 1.5rem;
    width: 80%;
  }

  .login-title {
    font-size: 1.6rem;
  }

  .login-input {
    font-size: 0.95rem;
    padding: 0.6rem;
  }

  .login-button {
    font-size: 1rem;
    padding: 0.7rem;
  }
}
</style>
