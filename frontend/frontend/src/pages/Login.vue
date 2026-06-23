<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";

import api from "../api/axios";

const router = useRouter();

const email = ref("");
const password = ref("");

const login = async () => {

  try {

    const response = await api.post("/login", {
      email: email.value,
      password: password.value,
    });

    localStorage.setItem(
      "token",
      response.data.token
    );

    router.push("/map");

  } catch (error) {
    alert("Login error");
  }
};
</script>

<template>
  <div class="auth-page">

    <div class="card">

      <h1>Zaрядись!</h1>

      <p class="subtitle">
        Войти что бы продолжить
      </p>

      <input
        v-model="email"
        placeholder="Email"
      />

      <input
        v-model="password"
        type="password"
        placeholder="Пароль"
      />

      <button @click="login">
        Войти
      </button>

      <a href="/register">
        Создать аккаунт
      </a>

    </div>

  </div>
</template>

<style>
.auth-page {
  width: 100%;
  height: 100vh;

  display: flex;
  justify-content: center;
  align-items: center;

  background:
    linear-gradient(
      180deg,
      #0f172a,
      #111827
    );
}

.card {
  width: 90%;
  max-width: 360px;

  background: #1e293b;

  border-radius: 24px;

  padding: 28px;

  display: flex;
  flex-direction: column;

  gap: 14px;

  box-shadow: 0 4px 30px rgba(0,0,0,0.3);
}

.card h1 {
  color: white;
}

.subtitle {
  color: #94a3b8;
  margin-bottom: 10px;
}

input {
  border: none;

  background: #334155;

  padding: 14px;

  border-radius: 14px;

  color: white;
}

button {
  border: none;

  background: #22c55e;

  color: white;

  padding: 14px;

  border-radius: 14px;

  font-weight: bold;

  cursor: pointer;
}

a {
  text-align: center;
  color: #94a3b8;
}
</style>