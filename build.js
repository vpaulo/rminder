import { readFileSync, writeFileSync } from 'fs';
import mkdirp from 'mkdirp';
import pkg from 'glob';
import sass from 'sass';
import { dirname } from 'path';
import { minify } from 'terser';
import htmlMinifier from 'html-minifier';

const { sync: globSync } = pkg;

const files = globSync('./src/**/*.js');
const options = {
  parse: {},
  compress: false,
  mangle: true,
  output: {
    ast: true,
    code: true,
    keep_quoted_props: true,
  },
  toplevel: true,
};

files.map(async (file) => {
  // console.log(`Minifying ${file}`);
  const terserResult = await minify(readFileSync(file, 'utf8'), options);
  if (terserResult.error) {
    console.log(`Minifying ${file} error.`, terserResult.error);
  } else {
    const path = file.replace('src/', 'dist/');
    writeFile(path, terserResult.code);
  }
});

function writeFile(path, content) {
  mkdirp.sync(dirname(path));
  writeFileSync(path, content, 'utf8');
  console.log(`Minifying ${path} success.`);
}

// Create main styles file
(() => {
  const distPath = './dist/css/styles.css';
  const result = sass.renderSync({
    file: './src/scss/styles.scss',
    outputStyle: 'compressed',
  });

  if (result) {
    writeFile(distPath, result.css.toString());
  }
})();

// Create and minify html main page
(() => {
  const distPath = './dist/index.html';
  const html = readFileSync('./src/index.html', 'utf8');
  const result = htmlMinifier.minify(html, {
    collapseWhitespace: true,
    removeComments: true,
  });
  if (result) {
    writeFile(distPath, result);
  }
})();
