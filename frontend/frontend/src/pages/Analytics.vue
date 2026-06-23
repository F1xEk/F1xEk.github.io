<script setup>
import { onMounted, ref } from "vue";

import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend,
} from "chart.js";

import { Line } from "vue-chartjs";

import api from "../api/axios";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend
);

const analytics = ref(null);

const chartData = ref(null);

const loadAnalytics = async () => {

  try {

    const token = localStorage.getItem("token");

    const response = await api.get(
      "/analytics",
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    analytics.value = response.data;

    chartData.value = {

      labels:
        response.data.chart.map(
          item =>
            new Date(
              item.date
            ).toLocaleDateString()
        ),

      datasets: [
        {
          label: "Energy (kWh)",

          data:
            response.data.chart.map(
              item => item.energy
            ),
          borderColor: "#22c55e",

          backgroundColor:
            "rgba(34, 197, 94, 0.2)",

          pointBackgroundColor:
            "#22c55e",

          pointBorderColor:
            "#22c55e",

          fill: true,
  
          tension: 0.4,
        },
      ],
    };

  } catch (error) {
    console.log(error);
  }
};

onMounted(() => {
  loadAnalytics();
});
</script>

<template>
  <div class="analytics-page">

    <button
    class="back-btn"
    @click="$router.push('/map')"
  >
    ←
  </button>

    <h1>
      Аналитика
    </h1>

    <div
      v-if="analytics"
      class="stats-grid"
    >

      <div class="card">
        <span>Сессии</span>

        <h2>
          {{
            analytics.total_sessions
          }}
        </h2>
      </div>

      <div class="card">
        <span>Энергия</span>

        <h2>
          {{
            analytics.total_energy.toFixed(2)
          }}
          кВтч
        </h2>
      </div>

      <div class="card">
        <span>Стоимость</span>

        <h2>
          {{
            analytics.total_cost.toFixed(2)
          }}
          Руб
        </h2>
      </div>

      <div class="card">
        <span>Средняя длительность</span>

        <h2>
          {{
            analytics.avg_duration.toFixed(0)
          }}
          сек
        </h2>
      </div>

    </div>

    <div class="chart-card">

      <Line
        v-if="chartData"
        :data="chartData"
      />

    </div>

  </div>
</template>

<style>
.analytics-page {
  min-height: 100vh;

  background: #0f172a;

  color: white;

  padding: 20px;
}

.stats-grid {
  display: grid;

  grid-template-columns:
    repeat(2, 1fr);

  gap: 14px;

  margin-top: 20px;
}

.card {
  background: #1e293b;

  border-radius: 20px;

  padding: 20px;
}

.card span {
  opacity: 0.7;
}

.card h2 {
  margin-top: 10px;
}

.chart-card {
  margin-top: 24px;

  background: #1e293b;

  border-radius: 24px;

  padding: 20px;
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