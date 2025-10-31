<template>
  <div class="userlist-page">
    <div class="userlist-container">
      <div class="top-section">
        <h1 class="userlist-title">Connected Players</h1>
      </div>

      <table class="userlist-table">
        <thead>
          <tr>
            <th>Username</th>
            <th>Pi Serial</th>
            <th>Death Sound</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="player in players" :key="player.id || player.username">
            <td>{{ player.username }}</td>
            <td>{{ player.pi_sn || '-' }}</td>
            <td>{{ player.death_sound || 'None' }}</td>
            <td>
              <button class="edit-btn" @click="editUser(player)">Edit</button>
              <button class="delete-btn" @click="deleteUser(player)">Delete</button>
            </td>
          </tr>
          <tr v-if="players.length === 0">
            <td colspan="4">No players found.</td>
          </tr>
        </tbody>
      </table>

      <router-link to="/adminboard">
        <button class="back-btn">Back to Adminboard</button>
      </router-link>

      <div v-if="message" class="message-box">{{ message }}</div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AdminEdit',
  data() {
    return {
      players: [],
      message: ''
    };
  },
  methods: {
    async fetchUsers() {
      const token = localStorage.getItem('authToken');
      if (!token) {
        alert('You must log in first.');
        return;
      }

      try {
        const response = await fetch('/api/api/users', {
          headers: { Authorization: `Bearer ${token}` }
        });
        if (!response.ok) {
          throw new Error(`Failed to fetch users: ${response.status}`);
        }

        const data = await response.json();
        console.log('Fetched user list:', data);

        this.players = Array.isArray(data) ? data : [];
      } catch (error) {
        console.error('Error loading user list:', error);
        this.message = 'Error loading user list.';
      }
    },

    editUser(player) {
      alert(`Editing ${player.username}`);
      this.$router.push('/userboard');
    },

    async deleteUser(player) {
      const token = localStorage.getItem('authToken');
      if (!token) {
        alert('You must log in first.');
        return;
      }

      if (!confirm(`Are you sure you want to delete ${player.username}?`)) return;

      try {
        const response = await fetch('/api/api/user', {
          method: 'DELETE',
          headers: { Authorization: `Bearer ${token}` }
        });

        const data = await response.json();
        if (response.ok && data.status === 'success') {
          alert(`Deleted ${player.username} successfully.`);
          this.fetchUsers();
        } else {
          alert(data.error || 'Failed to delete user.');
        }
      } catch (error) {
        console.error('Delete failed:', error);
        alert('Network error while deleting user.');
      }
    }
  },
  mounted() {
    this.fetchUsers();
  }
};
</script>

<style scoped>
.userlist-page {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100%;
  background: none;
}

.userlist-container {
  width: 55%;
  max-width: 800px;
  background: #fff;
  border: 0.25vw solid #000;
  border-radius: 1vw;
  padding: 2.5vw 2vw;
  box-shadow: 0 0.8vw 1.5vw rgba(0, 0, 0, 0.25);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.top-section {
  text-align: center;
  margin-bottom: 2vw;
}

.userlist-title {
  font-size: 1.8vw;
  font-weight: 700;
  text-align: center;
}

.userlist-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 2vw;
  text-align: center;
  font-size: 1vw;
}

.userlist-table th,
.userlist-table td {
  border: 0.1vw solid #ddd;
  padding: 0.8vw;
}

.userlist-table thead {
  background-color: #f7f7f7;
  font-weight: 600;
}

.userlist-table tbody tr:nth-child(even) {
  background-color: #f2f2f2;
}

.edit-btn {
  background-color: #4CAF50;
  color: black;
  font-weight: 600;
  padding: 0.4vw 1vw;
  border: none;
  border-radius: 0.3vw;
  margin-right: 1vw;
  cursor: pointer;
  transition: 0.2s ease;
}

.edit-btn:hover {
  background-color: #45a049;
  color: white;
}

.delete-btn {
  background-color: #8B0000;
  color: white;
  font-weight: 600;
  padding: 0.4vw 1vw;
  border: none;
  border-radius: 0.3vw;
  cursor: pointer;
  transition: 0.2s ease;
}

.delete-btn:hover {
  background-color: #b22222;
}

.back-btn {
  font-weight: 600;
  font-size: 1vw;
  border-radius: 0.5vw;
  padding: 0.8vw 1.8vw;
  cursor: pointer;
  transition: all 0.25s ease;
  border: 0.15vw solid #000;
  box-shadow: 0 0.4vw 0.8vw rgba(0, 0, 0, 0.1);
  background-color: #ffffff;
  color: black;
  border: 4px solid #000000;
}

.back-btn:hover {
  background-color: #f2f2f2;
  transform: translateY(-0.1vw);
}

.back-btn:active {
  transform: translateY(0); 
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.15); 
}

@media (max-width: 768px) {
  .userlist-container {
    width: 90%;
    padding: 5vw 4vw;
    border-radius: 3vw;
  }

  .userlist-title {
    font-size: 5vw;
    margin-bottom: 3vw;
  }

  .userlist-table {
    font-size: 3.5vw;
    border-spacing: 0;
  }

  .userlist-table th,
  .userlist-table td {
    padding: 3vw 2vw;
    border-width: 0.3vw;
  }

  .userlist-container {
    overflow-x: auto;
  }

  .edit-btn,
  .delete-btn {
    font-size: 3.5vw;
    padding: 1.5vw 3vw;
    border-radius: 2vw;
    margin: 1vw 0.5vw;
  }

  .back-btn {
    font-size: 3.5vw;
    padding: 2vw 5vw;
    border-radius: 3vw;
  }
}

@media (max-width: 480px) {
  .userlist-container {
    margin-top: -100px;
    width: 80%;
    padding: 6vw 5vw;
  }

  .userlist-title {
    font-size: 6vw;
  }

  .userlist-table {
    font-size: 4vw;
  }

  .edit-btn,
  .delete-btn {
    font-size: 3.8vw;
    padding: 2vw 4vw;
  }

  .back-btn {
    font-size: 4vw;
    padding: 2.5vw 6vw;
  }
}

</style>
