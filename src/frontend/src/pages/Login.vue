<template>
  <div class="login-page">
    <div class="login-box">
      <h1 class="login-title">Login</h1>
      <form class="login-form">
        <input type="text" placeholder="Username" class="login-input" />
        <input type="password" placeholder="Password" class="login-input" />
        <button type="submit" class="login-button">Log In</button>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AppLogin'
}
</script>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh; 
  width: 100%;
}

.login-box {
  background: white;
  padding: 40px 50px;
  border-radius: 20px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
  width: 360px;
  text-align: center;
}

.login-title {
  margin-bottom: 30px;
  font-size: 32px;
  font-weight: 700;
  color: #ff4b4b;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.login-input {
  padding: 12px;
  border-radius: 8px;
  border: 1px solid #ccc;
  font-size: 16px;
  outline: none;
  transition: 0.2s ease;
}

.login-input:focus {
  border-color: #ff4b4b;
  box-shadow: 0 0 6px rgba(255, 75, 75, 0.4);
}

.login-button {
  background: #ff4b4b;
  color: white;
  font-size: 18px;
  font-weight: 600;
  padding: 12px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.login-button:hover {
  background: #ff6666;
  transform: scale(1.03);
}
</style>
