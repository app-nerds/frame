/*
 * Copyright © 2022 App Nerds LLC
 */

export const debounce = (fn, delay = 400) => {
  let id = null;

  return function() {
    let args = arguments;

    clearTimeout(id);

    id = setTimeout(() => {
      fn.apply(this, args);
    }, delay);
  };
}

