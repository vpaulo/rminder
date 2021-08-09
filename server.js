import path from 'path';
import express from 'express';
import { fileURLToPath } from 'url';

const port = process.env.PORT || 5500;
const app = express();

// we need to change up how __dirname is used for ES6 purposes
const __dirname = path.dirname(fileURLToPath(import.meta.url));

app.use(express.static(path.join(__dirname, 'docs')));
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

app.listen(port, async () => {});
