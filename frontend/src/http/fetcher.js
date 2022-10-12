/*
 * Copyright © 2022 App Nerds LLC
 */

export async function fetcher(url, options, spinner, msBeforeShowSpinner = 1000) {
  let timerID;

  if (spinner) {
    timerID = setTimeout(() => {
      spinner.show();
    }, msBeforeShowSpinner);
  }

  const response = await fetch(url, options);

  if (spinner) {
    clearTimeout(timerID);
    spinner.hide();
  }

  return response;
}

