export function logger(o, ...n) {
  if (window.location.hash === '#debug') {
    console.log(o, n);
  }
}
