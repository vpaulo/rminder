/** @typedef {[null, T]} Success */
/** @typedef {[E]} Failure */
/** @typedef {Success<T> | Failure<E>} Result */

/**
 * @param {Promise<T>} promise
 * @returns {Promise<Result<T,E>}
 */
export async function tryCatch(promise) {
  try {
    const data = await promise;
    return [null, data];
  } catch (error) {
    return [error];
  }
}
