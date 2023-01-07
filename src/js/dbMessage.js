import { logger } from './logger.js';

export function dbMessage(app, data, db) {
  if (app[data.type]) {
    app[data.type](data, db);
    return;
  }

  logger(`Error running db worker - type: ${data.type} does not exist`);
}
