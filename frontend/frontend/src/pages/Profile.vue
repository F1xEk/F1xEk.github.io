<script setup>
import { onMounted, ref } from "vue";

import api from "../api/axios";

const user = ref(null);

const sessions = ref([]);

const loadProfile = async () => {

  try {

    const token = localStorage.getItem("token");

    const profileResponse =
      await api.get("/profile", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

    user.value = profileResponse.data;

    const sessionsResponse =
      await api.get("/sessions", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

    sessions.value =
      sessionsResponse.data || [];

  } catch (error) {
    console.log(error);
  }
};

const logout = () => {
  localStorage.removeItem("token");
  window.location.href = "/login";
};

onMounted(() => {
  loadProfile();
});
</script>

<template>
  <div class="profile-page">
    <button
    class="back-btn"
    @click="$router.push('/map')"
  >
    ← 
  </button>
    <div class="profile-card">
      
      <div class="avatar">
        A
      </div>

      <h2>
        {{ user?.name }}
      </h2>

      <p>
        {{ user?.email }}
      </p>

    </div>

    <div class="history">

      <h3>
        История зарядных сессий
      </h3>

      <div
        v-for="session in sessions"
        :key="session.id"
        class="session-card"
      >

        <div class="row">
          <span>Сессия</span>

          <b>#{{ session.id }}</b>
        </div>

        <div class="row">
          <span>Энергия</span>

          <b>
            {{ session.energy.toFixed(2) }}
            кВтч
          </b>
        </div>

        <div class="row">
          <span>Стоимость</span>

          <b>
            {{ session.cost.toFixed(2) }}
            Руб
          </b>
        </div>

        <div class="row">
          <span>Длительность</span>

          <b>
            {{ session.duration }}
            сек
          </b>
        </div>

      </div>

    </div>

    <button
      class="logout-btn"
      @click="logout"
    >
      Выход
    </button>

  </div>
</template>

<style>
.profile-page {
  min-height: 100vh;

  background: #0f172a;

  padding: 20px;

  color: white;
}

.profile-card {
  background: #1e293b;

  border-radius: 24px;

  padding: 24px;

  text-align: center;

  margin-bottom: 24px;
}

.avatar {
  width: 80px;
  height: 80px;

  border-radius: 50%;

  background: #22c55e;

  margin: 0 auto 16px;

  display: flex;
  align-items: center;
  justify-content: center;

  font-size: 32px;
  font-weight: bold;
}

.history h3 {
  margin-bottom: 16px;
}

.session-card {
  background: #1e293b;

  border-radius: 20px;

  padding: 18px;

  margin-bottom: 14px;
}

.row {
  display: flex;
  justify-content: space-between;

  margin-bottom: 10px;
}

.logout-btn {
  width: 100%;

  margin-top: 20px;

  border: none;

  background: #ef4444;

  color: white;

  padding: 16px;

  border-radius: 16px;

  font-weight: bold;
}
.back-btn {

  background: #22c55e;

  color: white;

  border: none;

  padding: 10px 16px;

  border-radius: 12px;

  font-size: 14px;

  font-weight: 600;

  cursor: pointer;

  margin-bottom: 20px;

}

.back-btn:hover {

  opacity: 0.9;

}
</style>