/*
 * Copyright © 2022 App Nerds LLC
 */

export async function fetcher(url, options, spinner, msBeforeShowSpinner = 1000) {
  let timerID;
  let response;

  if (spinner) {
    timerID = setTimeout(() => {
      spinner.show();
    }, msBeforeShowSpinner);
  }

  try {
    response = await fetch(url, options);
  } finally {
    if (spinner) {
      clearTimeout(timerID);
      spinner.hide();
    }
  }

  return response;
}

