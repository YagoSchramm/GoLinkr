import http from 'k6/http';

export const options = {
  vus: 1,
  iterations: 1,
};

export default function () {
  const res = http.get('http://localhost:8080/link/abc123', {
    redirects: 0,
  });

  console.log(`Status: ${res.status}`);
  console.log(`Location: ${res.headers.Location}`);
}