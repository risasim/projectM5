<template>
  <div class="userboard-container">
    <div class="userboard">
      <!-- Header -->
      <div class="userboard-header">
        <h2>Userboard</h2>
      </div>

      <!-- Main content area -->
      <div class="userboard-content">
        <div class="user-info">

          <span class="label">Username:</span>
        </div>

        <div class="user-sfx">
          <span class="label">Custom Death SFX: </span>
          <button class="insert-btn">Insert</button>
        </div>

        <div class="session-info">
          <div class="session-status">
            <span class="label">Session status: </span>
            <span class="value">Waiting for players</span>
          </div>
        </div>

        <div class="session-btn">
          <button class="enter-btn">Enter current game session</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Userboard"
};
</script>

<style scoped>
.userboard-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: transparent;
}

.userboard {
  background: rgba(255, 255, 255, 0.9);
  border: 2px solid black;
  border-radius: 8px;
  width: 1200px;
  height: 500px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.userboard-header {
  background: white;
  border: 2px solid black;
  border-radius: 6px 6px 0 0;
  padding: 8px 24px;
  margin-top: 24px;
  z-index: 2;
}

.userboard-header h2 {
  margin: 2;
  font-size: 1.5rem;
  font-weight: 800;
  color: black;
}

.userboard-content {
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
  padding: 2rem;
  width: 100%;
}

.user-info{
    margin-left: 2%;
    font-size: x-large;
}

.session-info{
margin-left: 38.4%;
}

.session-btn {
  display: flex;
  justify-content: space-between;
  align-items: right;
margin-left: 40%;
}

.user-sfx{
margin-top: 9%;
margin-left: 2%;

}

.label {
  font-weight: bold;
  
  color: #111;
}

.value {
  color: #333;
}

.insert-btn {
  background: #ddd;
  padding: 0.5% 4%;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
  transition: 0.2s;
}

.insert-btn:hover {
  background: #bbb;
}

.enter-btn {
  border: 2px solid black;
  background: transparent;
  padding: 10px 16px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.25s ease;
}

.enter-btn:hover {
  background: black;
  color: white;
}
</style>
