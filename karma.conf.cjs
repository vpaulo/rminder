// Karma configuration
// Generated on Wed Apr 07 2021 17:20:59 GMT+0100 (British Summer Time)

module.exports = function (config) {
  config.set({
    // base path that will be used to resolve all patterns (eg. files, exclude)
    basePath: './src',

    // frameworks to use
    // available frameworks: https://npmjs.org/browse/keyword/karma-adapter
    frameworks: ['mocha', 'chai', 'testdouble'],

    plugins: [
      'karma-chai',
      'karma-mocha',
      'karma-mocha-reporter',
      'karma-chrome-launcher',
      'karma-coverage',
      'karma-testdouble',
    ],

    // list of files / patterns to load in the browser
    files: [
      { pattern: 'js/workers/dbw.js', included: false },
      { pattern: '**/*.js', type: 'module' },
      { pattern: '../test/js/*.spec.js', type: 'module' },
      { pattern: '../test/App.spec.js', type: 'module' },
    ],

    // list of files / patterns to exclude
    exclude: [],

    // preprocess matching files before serving them to the browser
    // available preprocessors: https://npmjs.org/browse/keyword/karma-preprocessor
    preprocessors: {
      '**/*.js': ['coverage'],
    },

    // test results reporter to use
    // possible values: 'dots', 'progress'
    // available reporters: https://npmjs.org/browse/keyword/karma-reporter
    reporters: ['mocha', 'coverage'],

    // web server port
    port: 5500,

    // enable / disable colors in the output (reporters and logs)
    colors: true,

    // level of logging
    // possible values: config.LOG_DISABLE || config.LOG_ERROR || config.LOG_WARN || config.LOG_INFO || config.LOG_DEBUG
    logLevel: config.LOG_INFO,

    // enable / disable watching file and executing tests whenever any file changes
    autoWatch: false,

    // start these browsers
    // available browser launchers: https://npmjs.org/browse/keyword/karma-launcher
    browsers: ['ChromeHeadless'],

    // Continuous Integration mode
    // if true, Karma captures browsers, runs the tests and exits
    singleRun: true,
    restartOnFileChange: true,

    // Concurrency level
    // how many browser should be started simultaneous
    concurrency: Infinity,

    // optionally, configure the reporter
    coverageReporter: {
      dir: '../coverage/',
      reporters: ['html', 'lcovonly', 'text-summary'],
      check: {
        global: {
          statements: 85,
          branches: 85,
          functions: 85,
          lines: 85,
        },
      },
    },
    client: {
      mocha: {
        timeout: 10000, // 10 seconds - upped from 2 seconds
      },
    },
    verbose: true,
  });
};
