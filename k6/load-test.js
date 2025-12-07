import http from 'k6/http';
import { check } from 'k6';
import { Rate } from 'k6/metrics';

const successRate = new Rate('success_rate');

export const options = {
  thresholds: {
    http_req_duration: ['p(95)<50', 'p(99)<100'], 
    success_rate: ['rate>0.99'], 
  },

  stages: [
    { duration: '30s', target: 50 },  // Ramp up to 50 concurrent backends
    { duration: '1m', target: 50 },   // Stability check: 50 backends
    { duration: '30s', target: 100 }, // Ramp up to 100 concurrent backends (Stress check)
    { duration: '1m', target: 100 },  // Peak load: 100 backends
    { duration: '30s', target: 0 },   // Ramp down
  ],
};

function getRandomIP() {
  const r = () => Math.floor(Math.random() * 255);
  return `${r()}.${r()}.${r()}.${r()}`;
}

export default function () {
  const ip = getRandomIP();
  const url = `http://localhost:8080/lookup/${ip}`;

  const res = http.get(url);

  const result = check(res, {
    'status is 200': (r) => r.status === 200,
    'content type is json': (r) => r.headers['Content-Type'] && r.headers['Content-Type'].includes('application/json'),
    'has valid location': (r) => r.body.includes('location'),
  });

  successRate.add(result);
}
