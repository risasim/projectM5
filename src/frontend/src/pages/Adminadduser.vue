<template>
  <div class="page-container">
    <div class="create-account-page">
      <div class="create-account-box">
        <button class="admin-panel-btn" @click="$router.push('/adminboard')">
          Back to Admin Panel
        </button>

        <h1 class="create-account-title">Add New User</h1>
        <h3 class="create-account-subtitle">Enter User Details</h3>

        <div class="add-user-section">
          <input
            v-model="newUser.username"
            type="text"
            placeholder="Username"
            class="create-account-input"
          />
          <input
            v-model="newUser.password"
            type="password"
            placeholder="Password"
            class="create-account-input"
          />
          <input
            v-model="newUser.pi_sn"
            type="text"
            placeholder="Pi Serial Number"
            class="create-account-input"
          />

          <button class="create-account-button" @click="addUser">Add User</button>

          <p v-if="message" class="status-message">{{ message }}</p>
          <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>
          <p v-if="successMessage" class="success-message">{{ successMessage }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "CreateAccount",
  data() {
    return {
      newUser: { username: "", password: "", pi_sn: "" },
      message: "",
      errorMessage: "",
      successMessage: ""
    };
  },
  methods: {
    getAuthToken() {
      const token = localStorage.getItem("authToken");
      if (!token) {
        alert("You must log in first.");
        this.$router.push("/login");
        return null;
      }
      return token;
    },

    async addUser() {
      console.log("[addUser] called", this.newUser);
      const token = this.getAuthToken();
      if (!token) return;

      if (!this.newUser.username || !this.newUser.password || !this.newUser.pi_sn) {
        this.message = "Please fill all fields before adding a user.";
        return;
      }

      if (this.newUser.password.length < 8) {
        this.message = "Password must be at least 8 characters long.";
        return;
      }

      try {
        const res = await fetch("/api/addUser", {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json"
          },
          body: JSON.stringify(this.newUser)
        });

        const text = await res.text();
        let data;
        try {
          data = JSON.parse(text);
        } catch {
          data = { error: text };
        }

        console.log("[addUser] response:", res.status, data);

        if (res.ok && data.status === "success") {
          this.message = data.message || "User added successfully.";
          this.newUser = { username: "", password: "", pi_sn: "" };
          if (this.fetchAllUsers) await this.fetchAllUsers();
        } else {
          this.message = data.error || "Failed to add user.";
        }
      } catch (err) {
        console.error("[addUser] exception:", err);
        this.message = "Network error adding user.";
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
.create-account-page {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  width: 100%;
}
.create-account-box {
  position: relative;
  background: white;
  padding: 2.5rem 2rem 2.8rem;
  border-radius: 20px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
  width: 380px;
  max-width: 90%;
  text-align: center;
  border: 4px solid #000;
}
.admin-panel-btn {
  position: absolute;
  top: 15px;
  left: 15px;
  background: white;
  border: 2px solid #000;
  border-radius: 9999px;
  padding: 6px 14px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}
.admin-panel-btn:hover {
  background: #f5f5f5;
  transform: scale(1.04);
}
.create-account-title {
  margin-top: 3rem;
  margin-bottom: 1rem;
  font-size: 2rem;
  font-weight: 700;
  color: #ff4b4b;
}
.create-account-subtitle {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 1.5rem;
}
.add-user-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}
.create-account-input {
  width: 80%;
  padding: 0.75rem;
  border-radius: 8px;
  border: 1px solid #ccc;
  font-size: 1rem;
  outline: none;
  transition: 0.2s ease;
  text-align: center;
}
.create-account-input:focus {
  border-color: #ff4b4b;
  box-shadow: 0 0 6px rgba(255, 75, 75, 0.4);
}
.create-account-button {
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
  width: 82%;
}
.create-account-button:hover {
  background: #ff6666;
  transform: scale(1.03);
}
.status-message {
  color: #000;
  font-weight: 500;
  margin-top: 10px;
}
.error-message {
  color: red;
  margin-top: 10px;
}
.success-message {
  color: green;
  margin-top: 10px;
}
</style>