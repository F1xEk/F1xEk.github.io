<script setup>
import { onMounted, ref } from "vue";

import L from "leaflet";

import {
  greenIcon,
  orangeIcon,
  redIcon,
} from "../utils/icons";

import api from "../api/axios";

const markers = {};

const charging = ref(false);

const sessionId = ref(null);

const ocppData = ref(null);

const selectedStation = ref(null);

const chargingInfo = ref({
  station: "",
  seconds: 0,
  energy: 0,
  cost: 0,
});

let timer = null;

const logout = () => {
  localStorage.removeItem("token");
  window.location.href = "/login";
};

const connectWebSocket = () => {

  const ws = new WebSocket(
    "ws://localhost:8080/api/ws"
  );

  ws.onmessage = (event) => {

    const data = JSON.parse(event.data);

    ocppData.value = data;

    if (
      charging.value &&
      data.status === "Charging"
    ) {
      chargingInfo.value.energy =
        data.energy;
    }

    const marker =
      markers[data.station_id];

    if (!marker) return;

    if (data.status === "Available") {
      marker.setIcon(greenIcon);
    }

    if (data.status === "Charging") {
      marker.setIcon(orangeIcon);
    }

    if (data.status === "Offline") {
      marker.setIcon(redIcon);
    }
  };
};

const startCharging = async (station) => {

  try {

    const token = localStorage.getItem("token");

    const response = await api.post(
      "/start-session",
      {
        station_id: station.id,
      },
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    sessionId.value =
      response.data.session_id;

    charging.value = true;
    selectedStation.value = null;

    connectWebSocket();

    chargingInfo.value.station =
      station.title;

    chargingInfo.value.seconds = 0;
    chargingInfo.value.energy = 0;
    chargingInfo.value.cost = 0;

    timer = setInterval(() => {

      chargingInfo.value.seconds += 1;

      chargingInfo.value.energy += 0.02;

      chargingInfo.value.cost += 0.1;

    }, 1000);

  } catch (error) {
    console.log(error);
  }
};

const stopCharging = async () => {

  try {

    const token = localStorage.getItem("token");

    await api.post(
      "/stop-session",
      {
        session_id: sessionId.value,
        energy: chargingInfo.value.energy,
        cost: chargingInfo.value.cost,
        duration: chargingInfo.value.seconds,
      },
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    clearInterval(timer);

    charging.value = false;

  } catch (error) {
    console.log(error);
  }
};

const loadMap = async () => {

  const map = L.map("map", {
    zoomControl: false,
  }).setView([55.801713, 49.112377], 12);

  L.tileLayer(
    "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png",
    {
      attribution:
        "&copy; OpenStreetMap contributors",
    }
  ).addTo(map);

  try {

    const token = localStorage.getItem("token");

    const response = await api.get("/stations", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    const stations = response.data || [];

    stations.forEach((station) => {

      const marker = L.marker([
        station.latitude,
        station.longitude,
      ], {
      icon: greenIcon,
      }).addTo(map);

      markers[station.id] = marker;

     marker.on("click", () => {

      selectedStation.value = station;

     });

      // marker.bindPopup(`
      //   <div>
      //     <h3>${station.title}</h3>

      //     <p>${station.address}</p>

      //     <p>Status: ${station.status}</p>

      //     <p>Connector: ${station.connector}</p>

      //     <p>Tariff: ${station.tariff}</p>

      //     <button>
      //       Start Session
      //     </button>
      //   </div>
      // `);

    });

  } catch (error) {
    console.log(error);
  }
};

onMounted(() => {

  loadMap();

  connectWebSocket();
});
</script>

<template>
  <div class="app">

    <div class="topbar">

      <div>
        <h2>Zaрядись!</h2>
        <p>Зарядные станции рядом</p>
      </div>

      <div class="avatar">
        A
      </div>

    </div>

    <div id="map"></div>

    <div
      v-if="selectedStation !== null && !charging"
      class="station-card"
    >

      <div class="station-header">

        <h2>
          {{ selectedStation.title }}
        </h2>

        <button
          class="close-btn"
          @click="selectedStation = null"
        >
          ✕
        </button>

      </div>

      <p>
        📍 {{ selectedStation.address }}
      </p>

      <p>
        🔌 {{ selectedStation.connector }}
      </p>

      <p>
        ⚡ {{ selectedStation.tariff }}
      </p>

      <p>
        🟢 {{ selectedStation.status }}
      </p>

      <button
        class="start-btn"
        @click="startCharging(selectedStation)"
      >
        Начать сессию
      </button>

    </div>

    <div
      v-if="charging"
      class="charging-card"
    >

      <h2>Зарядка...</h2>

      <div
        v-if="ocppData"
        class="ocpp-live"
      >

        <div class="live-row">
          <span>Статус</span>
          <b>{{ ocppData.status }}</b>
        </div>

        <div class="live-row">
          <span>Мощность</span>
          <b>{{ ocppData.power }} kW</b>
        </div>

        <div class="live-row">
          <span>Напряжение</span>
          <b>{{ ocppData.voltage }} V</b>
        </div>

        <div class="live-row">
          <span>Сила тока</span>
          <b>{{ ocppData.current }} A</b>
        </div>

      </div>

      <p>
        {{ chargingInfo.station }}
      </p>

      <div class="stats">

        <div class="stat">
          <span>Время</span>

          <b>
            {{ chargingInfo.seconds }} сек
          </b>
        </div>

        <div class="stat">
          <span>Энергия</span>

          <b>
            {{
              chargingInfo.energy.toFixed(2)
            }}
            кВтч
          </b>
        </div>

        <div class="stat">
          <span>Стоимость</span>

          <b>
            {{
              chargingInfo.cost.toFixed(2)
            }}
            Руб
          </b>
        </div>

      </div>

      <button
        class="stop-btn"
        @click="stopCharging"
      >
        Остановить сессию
      </button>

    </div>

    <div class="bottom-nav">

      <a
        href="/map"
        class="nav-item active"
      >
        Карта
      </a>

      <a
        href="/profile"
        class="nav-item"
      >
        Профиль
      </a>
      <a
      href="/analytics"
      class="nav-item"
      >
        Аналитика
      </a>

      <button
        class="logout-btn"
        @click="logout"
      >
        Выйти
      </button>

    </div>

  </div>
</template>

<style>
.app {
  width: 100%;
  height: 100vh;
  overflow: hidden;
  position: relative;
}

#map {
  width: 100%;
  height: 100vh;
}

.topbar {
  position: absolute;
  top: 20px;
  left: 20px;
  right: 20px;

  z-index: 1000;

  background: rgba(15, 23, 42, 0.95);

  backdrop-filter: blur(10px);

  border-radius: 20px;

  padding: 16px;

  display: flex;
  justify-content: space-between;
  align-items: center;

  color: white;

  box-shadow: 0 4px 20px rgba(0,0,0,0.3);
}

.topbar p {
  opacity: 0.7;
  margin-top: 4px;
}

.avatar {
  width: 44px;
  height: 44px;

  border-radius: 50%;

  background: #22c55e;

  display: flex;
  align-items: center;
  justify-content: center;

  font-weight: bold;
}

.bottom-nav {
  position: absolute;

  bottom: 20px;
  left: 20px;
  right: 20px;

  z-index: 1000;

  background: rgba(15, 23, 42, 0.95);

  backdrop-filter: blur(10px);

  border-radius: 24px;

  padding: 14px;

  display: flex;
  justify-content: space-around;
  align-items: center;

  box-shadow: 0 4px 20px rgba(0,0,0,0.3);
}

.nav-item {
  color: white;
  opacity: 0.7;
  font-size: 15px;
  margin-right: 10px;
}

.active {
  opacity: 1;
  font-weight: bold;
}

.logout-btn {
  border: none;
  background: #ef4444;
  color: white;

  padding: 10px 14px;

  border-radius: 12px;

  cursor: pointer;
}

.leaflet-popup-content-wrapper {
  border-radius: 18px;
}

.popup {
  min-width: 220px;
}

.popup h3 {
  margin-bottom: 10px;
}

.popup p {
  margin-bottom: 12px;
}

.info-row {
  display: flex;
  justify-content: space-between;

  margin-bottom: 8px;
}

.start-btn {
  width: 100%;

  margin-top: 12px;

  border: none;

  background: #22c55e;
  color: white;

  padding: 12px;

  border-radius: 12px;

  cursor: pointer;

  font-weight: bold;
}
.charging-card {
  position: absolute;

  left: 20px;
  right: 20px;
  bottom: 110px;

  z-index: 1000;

  background: rgba(15, 23, 42, 0.96);

  backdrop-filter: blur(10px);

  border-radius: 24px;

  padding: 20px;

  color: white;

  box-shadow: 0 4px 30px rgba(0,0,0,0.3);
}

.stats {
  margin-top: 16px;

  display: flex;
  justify-content: space-between;
}

.stat {
  display: flex;
  flex-direction: column;
}

.stat span {
  opacity: 0.7;
  margin-bottom: 6px;
}

.stop-btn {
  width: 100%;

  margin-top: 20px;

  border: none;

  background: #ef4444;

  color: white;

  padding: 14px;

  border-radius: 14px;

  font-weight: bold;
}

.ocpp-live {
  margin-top: 16px;

  background: rgba(255,255,255,0.05);

  border-radius: 16px;

  padding: 14px;
}

.live-row {
  display: flex;
  justify-content: space-between;

  margin-bottom: 8px;
}
.station-card {

  position: fixed;

  bottom: 300px;

  left: 0;

  right: 0;

  background: white;

  border-radius: 24px 24px 0 0;

  padding: 20px;

  box-shadow: 0 -5px 20px rgba(0,0,0,0.15);

  z-index: 999;

  animation: slideUp .3s ease;
  color: #000;

}

.station-header {

  display: flex;

  justify-content: space-between;

  align-items: center;

}

.close-btn {

  border: none;

  background: none;

  font-size: 20px;

  cursor: pointer;
  color: #000;

}

.start-btn {

  width: 100%;

  margin-top: 15px;

  padding: 14px;

  border: none;

  border-radius: 12px;

  background: #16a34a;

  color: white;

  font-size: 16px;

  font-weight: bold;

  cursor: pointer;

}

@keyframes slideUp {

  from {

    transform: translateY(100%);

  }

  to {

    transform: translateY(0);

  }

}
</style>