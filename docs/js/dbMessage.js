import { logger as e } from './logger.js';
export function dbMessage(r, t, o) {
  if (r[t.type]) {
    r[t.type](t, o);
    return;
  }
  e(`Error running db worker - type: ${t.type} does not exist`);
}
