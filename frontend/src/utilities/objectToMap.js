/*
 * Copyright © 2022 App Nerds LLC
 */

/**
 * Converts a classic JS object to a Map
 * @param o object The object to convert
 */
export const objectToMap = (o = {}) => {
  let result = new Map();

  for (const key in o) {
    result.set(key, o[key]);
  }

  return result;
};

