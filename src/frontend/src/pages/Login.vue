<template>
  <div class="page-container">
    <div class="login-page">
      <div class="login-box">
        <h1 class="login-title">Login</h1>
        <form class="login-form" @submit.prevent>

              <input v-model="username" type="text" placeholder="Username" class="login-input" />
              <input v-model="password" type="password" placeholder="Password" class="login-input" />

              <button type="button" class="login-button" @click="loginUser(false)">Log In</button>
              <button type="button" class="login-button" @click="loginUser(true)">Log In (admin)</button>

              <p v-if="errorMessage" style="color:red">{{ errorMessage }}</p>
        </form>
      </div>
    </div>
  </div>  
</template>

<script>
export default {
  name: 'AppLogin',
  data() {
    return {
      username: '',
      password: '',
      errorMessage: '',
      isLoading: false
    };
  },
  methods: {
    async loginUser(isAdmin = false) {
      this.errorMessage = '';
      this.isLoading = true;

      try {
        //get the token
        const response = await fetch('/api/auth', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            username: this.username,
            password: this.password
          })
        });

        if (!response.ok) throw new Error('Login failed');
        const data = await response.json();

        // Store token
        localStorage.setItem('authToken', data.token);
        console.log('Token stored:', data.token);

        // Redirect
        if (isAdmin) {
          this.$router.push('/adminboard');
        } else {
          this.$router.push('/userboard');
        }
      } catch (err) {
        this.errorMessage = 'Invalid credentials or server error.';
        console.error(err);
      } finally {
        this.isLoading = false;
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
  color: #b30000;
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
  border-color: #b30000;
  box-shadow: 0 0 0.4vw rgba(255, 75, 75, 0.4);
}

.login-button {
  background: #b30000;
  background: #ff1111;
  color: white;
  font-size: 1.1rem;
  font-weight: 600;
  padding: 0.8rem;
  border: none;
  border-radius: 0.5vw;

  cursor: pointer;
  transition: all 0.2s ease;
  border: 4px solid #000;
  width: 100%;                
}

.login-button:hover {
  background: #b30000;
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
