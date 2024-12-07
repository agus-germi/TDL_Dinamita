import http from 'k6/http'; // Módulo para hacer peticiones HTTP
import { check } from 'k6'; // Módulo para realizar comprobaciones
import { sleep } from 'k6'; // Módulo para hacer pausas entre peticiones
import { Trend } from 'k6/metrics'; // Módulo para crear métricas personalizadas

// Definir métricas personalizadas
let myTrend = new Trend('waiting_time');

// Configurar el escenario de carga
export let options = {
  stages: [
    { duration: '30s', target: 50 }, // 50 usuarios durante 30 segundos
    { duration: '1m', target: 100 }, // 100 usuarios durante 1 minuto
    { duration: '30s', target: 0 },  // Reducir a 0 usuarios durante 30 segundos
  ],
};


export default function () {
  const user_id = Math.floor(Math.random() * 100) + 1;  // Genera un id de usuario aleatorio entre 1 y 100
  const table_number = Math.floor(Math.random() * 10) + 1;  // Genera un número de mesa aleatorio entre 1 y 10

  const url = 'http://api:8080/api/v1/reservations';
  const payload = JSON.stringify({
    user_id: user_id,
    table_number: table_number,
    date: '2024-12-10',
    time_slot: '19:00',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  // Hacer una petición HTTP al endpoint de la reserva
  let response = http.post(url, payload, params);

  // Comprobar que la petición ha tenido éxito
  check(response, {
    'is status 200': (r) => r.status === 200,
    'response time is less than 200ms': (r) => r.timings.duration < 200,
  });

  // Registrar métricas personalizadas
  myTrend.add(response.timings.duration);

  // Esperar entre 1 y 3 segundos antes de la siguiente petición
  sleep(Math.random() * 2 + 1);
}
