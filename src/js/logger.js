export function logger(message, ...args) {
  if (window.location.hash === '#debug') {
    console.log(message, args);
  }
}
