/*
 * Copyright © 2022 App Nerds LLC
 */

export async function fetcher(url, options, spinner) {
  let timerID;

  if (spinner) {
    timerID = setTimeout(() => {
      spinner.show();
    }, 1000);
  }

  const response = await fetch(url, options);

  if (spinner) {
    clearTimeout(timerID);
    spinner.hide();
  }

  return response;
}

