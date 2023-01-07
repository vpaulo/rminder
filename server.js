import path from 'path';
import { fileURLToPath } from 'url';
import Fastify from 'fastify';
import fastifyStatic from '@fastify/static';
import fastifyCors from '@fastify/cors';

const fastify = Fastify({
  logger: true,
});

const port = process.env.PORT || 5500;

// we need to change up how __dirname is used for ES6 purposes
const __dirname = path.dirname(fileURLToPath(import.meta.url));

fastify.register(fastifyCors, {
  origin: '*',
  methods: ['GET'],
});

// Serve static files from the "public" directory.
fastify.register(fastifyStatic, {
  root: path.join(__dirname, 'dist'),
});

// Run the server!
fastify.listen({ port }, (err) => {
  if (err) {
    fastify.log.error(err);
    process.exit(1);
  }
  fastify.log.info(`server listening on ${fastify.server.address().port}`);
});
