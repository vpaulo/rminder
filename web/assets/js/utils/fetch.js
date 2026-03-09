export function csrfToken() {
  return document.querySelector('meta[name="csrf-token"]')?.content || "";
}

/**
 * @param {string} method
 * @param {string} url
 * @param {FormData|URLSearchParams|string|null} body
 * @returns {Promise<string>}
 */
export async function apiHTML(method, url, body = null) {
  const headers = { "X-CSRF-Token": csrfToken() };
  const options = { method, headers };

  if (body instanceof FormData) {
    options.body = body;
  } else if (body) {
    headers["Content-Type"] = "application/x-www-form-urlencoded";
    options.body = body instanceof URLSearchParams ? body.toString() : body;
  }

  const resp = await fetch(url, options);
  return resp.text();
}

/**
 * @param {string} selector
 * @param {string} html
 * @param {"innerHTML"|"outerHTML"} mode
 */
export function swapHTML(selector, html, mode = "innerHTML") {
  const el = document.querySelector(selector);
  if (!el) return;
  if (mode === "outerHTML") {
    el.outerHTML = html;
  } else {
    el.innerHTML = html;
  }
}
