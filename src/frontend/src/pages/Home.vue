<template>
  <div class="home-page">
    <div class="home-content">
      <div class="home-center">
        <img src="@/assets/phosho-coollogo_com-blackbg.png" alt="Logo" class="home-logo" />
        <p class="home-description">
        </p>

        <button class="userboard-btn">Userboard</button>

        <div class="info-section">
          <div class="info-card">
            <h2 class="info-title">Gamemodes</h2>
            <p class="info-content">some info</p>
          </div>
          <div class="info-card">
            <h2 class="info-title">How It Works</h2>
            <p class="info-content">some info </p>
          </div>
          <div class="info-card">
            <h2 class="info-title">Raspberry Pi</h2>
            <p class="info-content">pie</p>
          </div>
        </div>
      </div>
    </div>
  </div>

</template>

<script>
export default {
  name: 'HomePage'
}
</script>

<style scoped>
.home-page {
  position: relative;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  color: #111;
  overflow: hidden;
}

.home-content {
  flex: 1;
  display: flex;
  justify-content: center;
  text-align: center;
}

.home-logo {
  margin-bottom: 15px;
}

.home-title {
  font-size: 3rem;
  font-weight: 800;
  color: #b30000; /* strong red */
  margin-bottom: 10px;
}

.home-description {
  color: #333;
}

.userboard-btn {
  background: #ffffff;
  color: rgb(0, 0, 0);
  font-family: 'Trebuchet MS', Arial, sans-serif;
  font-weight: 600;
  border: 4px solid #000000;
  border-radius: 19px;
  padding: 12px 28px;
  cursor: pointer;
  font-size: 2rem;
  transition: 0.3s;
}

.userboard-btn:hover {
  background: #ad0000;
  color: rgb(0, 0, 0);
  transform: scale(1.05);
  border: 3px solid #000000;
}

.info-section {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 20px;
  margin-top: 50px;
}

.info-card {
  background: #fff;
  border: 4px solid #000000;
  border-radius: 14px;
  width: 280px;
  padding: 25px;
  box-shadow: 0 8px 18px rgba(0, 0, 0, 0.15);
  transition: 0.3s ease;
}

.info-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(179, 0, 0, 0.3);
}
  /*font-family: 'Trebuchet MS', Arial, sans-serif;*/
.info-title {
  color: #b30000;
  font-size: 1.3rem;
  font-weight: 700;
  font-family: 'Trebuchet MS', Arial, sans-serif;
  margin-bottom: 10px;
}

.info-content {
  color: #333;
  font-size: 1rem;
  font-family: 'Trebuchet MS', Arial, sans-serif;
  line-height: 1.5;
}
</style>
  