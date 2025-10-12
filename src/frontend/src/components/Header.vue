<template>
  <header class="header-gradient">
    <!--<router-link to="/" class="header-title">PhoSho</router-link>-->
    <router-link to="/"> <img src="@/assets/phosho-coollogo_com.png" alt="Logo" class="header-title" /> </router-link>

    <router-link to="/login" class="login-button">Log In</router-link>
  </header>
</template>

<script>
export default {
  name: 'AppHeader'
}
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
  position: fixed; /* stays at top always */
  top: 0;
  left: 0;
  z-index: 10; /* above background but below modals if any */
}



.header-title {
  color: white;
  width: 350px;
}

.login-button {
  background: white;
  color: #000000;
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  text-decoration: none;
  border: 2px solid #000000;
  transition: 0.25s ease;
}

.login-button:hover {
  background: #ff0000;
  color: white;
  transform: scale(1.05);
}
</style>
