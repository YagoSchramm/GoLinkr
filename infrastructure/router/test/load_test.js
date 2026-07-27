import http from 'k6/http';
import { check } from 'k6';

export const options = {
    vus: 100,
    duration: '30s',
};

export default function () {
    const res = http.get('http://localhost:8080/abc123', {
        redirects: 0, 
    });

    check(res, {
        'status é 301 ou 302': (r) => r.status === 301 || r.status === 302,
    });
}