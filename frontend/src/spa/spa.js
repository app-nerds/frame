/*
 * Copyright © 2022 App Nerds LLC
 */

import { Router } from "./router.js";
import { DefaultPageNotFound } from "./base-view.js";

export const application = (
  targetElement,
  routes,
  pageNotFoundView = DefaultPageNotFound
) => {
  window._router = new Router(targetElement, routes, pageNotFoundView);
  window.navigateTo = window._router.navigateTo.bind(window._router);

  window.addEventListener("popstate", (e) => {
    window._router._route({
      state: e.state,
    });
  });

  return {
    routes: routes,
    targetElement: targetElement,
    router: window._router,

    afterRoute: (f) => {
      window._router.afterRoute = f.bind(window._router);
    },

    beforeRoute: (f) => {
      window._router.beforeRoute = f.bind(window._router);
    },

    injectParams: (f) => {
      window._router.injectParams = f.bind(window._router);
    },

    onRenderComplete: (f) => {
      window._router.onRenderComplete = f.bind(window._router);
    },

    go: () => {
      window._router._route({});
    },
  };
};
