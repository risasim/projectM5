<template>
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
  background: #ffffff;
  border: 0.2vw solid #000000;
  border-radius: 1.5vw;
  padding: 2vw;
  margin-bottom: 0;
  box-shadow: 0 0.5vw 1vw rgba(0, 0, 0, 0.15);
  width: 25%;
  text-align: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  border: 4px solid #000000;

}

.login-title {
  margin-top: -1vw;
  margin-bottom: 1.5vw;
  font-size: 2vw;
  font-weight: 700;
  color: #b30000;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1vw;
  width: 100%; 
}

.login-input {
  padding: 2.6%;
  border-radius: 0.5vw;
  border: 0.08vw solid #ccc;
  font-size: 1vw;
  outline: none;
  transition: 0.2s ease;
  margin-bottom: 0.5vw; 
}

.login-input:focus {
  border-color: #b30000;
  box-shadow: 0 0 0.4vw rgba(255, 75, 75, 0.4);
}

.login-button {
  background: #b30000;
  background: #ff1111;
  color: white;
  font-size: 1.1vw;
  font-weight: 600;
  padding: 2.6%;
  border: none;
  border-radius: 0.5vw;

  cursor: pointer;
  transition: all 0.2s ease;
  width: 360px;
  border: 4px solid #000000;
}

.login-button:hover {
  background: #b30000;
  transform: scale(1.03);
}
</style>
