<script setup>
if (localStorage.getItem("token")) {
  router.push("/map");
}
import { ref } from "vue";
import { useRouter } from "vue-router";

import api from "../api/axios";

const router = useRouter();

const name = ref("");
const email = ref("");
const password = ref("");

const register = async () => {
  try {

    await api.post("/register", {
      name: name.value,
      email: email.value,
      password: password.value,
    });

    router.push("/login");

  } catch (error) {
    alert("Register error");
  }
};
</script>

<template>
  <div class="container">
    <h1>Регистрация</h1>

    <input v-model="name" placeholder="Имя" />

    <input v-model="email" placeholder="Email" />

    <input
      v-model="password"
      type="password"
      placeholder="Пароль"
    />

    <a href="/login">
        Уже есть аккаунт? Войти
      </a>


    <button @click="register">
      Зарегистрироваться
    </button>
  </div>
</template>

<style>
.container {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 20px;
}
</style>