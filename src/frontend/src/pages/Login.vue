
<template>
  <div class="page-container">
    <div class="login-page">
      <div class="login-box">
        <h1 class="login-title">Login</h1>
        <form class="login-form" @submit.prevent>

              <input v-model="username" type="text" placeholder="Username" class="login-input" />
              <input v-model="password" type="password" placeholder="Password" class="login-input" />

              <button type="button" class="login-button" @click="loginUser(false)">Log In</button>

              <p v-if="errorMessage" style="color:red">{{ errorMessage }}</p>
        </form>
      </div>
    </div>
  </div>  
</template>

<script>
import { sanitizeUsername, sanitizePassword } from '@/utils/loginSanitizer'

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
      this.errorMessage = ''
      this.isLoading = true

      // sanitize the input before sending
      const user = sanitizeUsername(this.username)
      const pass = sanitizePassword(this.password)

      if (!user.valid || !pass.valid) {
        this.errorMessage = user.error || pass.error
        this.isLoading = false
        return
      }

      try {
        // get the token
        const response = await fetch('/api/api/auth', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            username: user.value,
            password: pass.value
          })
        });
        console.log("Res is: ",response)

        if (!response.ok) {
          let msg = `${response.status} ${response.statusText}`
          try { const e = await response.json(); msg = e.error || e.message || msg } catch {}
          throw new Error(msg)
        }

        const data = await response.json()
        console.log('Auth response data:', data)

        // store token
        if (!data.token) throw new Error('No token returned from server')
        localStorage.setItem('authToken', data.token)    
        localStorage.setItem('userRole', data.role)
        localStorage.setItem('username', user.value)

        // route according to role
        if (data.role === 'admin') {
          this.$router.push('/adminboard')
        } else {
          this.$router.push('/userboard')
        }
      } catch (err) {
        this.errorMessage = 'Invalid credentials or server error.'
        console.error('Login error:', err)
      } finally {
        this.isLoading = false
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
