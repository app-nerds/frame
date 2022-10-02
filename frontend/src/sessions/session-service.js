/*
 * Copyright © 2022 App Nerds LLC
 */

export const ErrTokenExpired = "token expired";

export default class SessionService {
  static clearMember() {
    window.sessionStorage.removeItem("member");
  }

  static clearToken() {
    window.sessionStorage.removeItem("token");
  }

  static getMember() {
    return JSON.parse(window.sessionStorage.getItem("member"));
  }

  static getToken() {
    return JSON.parse(window.sessionStorage.getItem("token"));
  }

  static hasMember() {
    return window.sessionStorage.getItem("member") !== null;
  }

  static hasToken() {
    return window.sessionStorage.getItem("token") !== null;
  }

  static navigateOnTokenExpired(e, path, navigateTo) {
    if (e.message === ErrTokenExpired) {
      SessionService.clearToken();
      navigateTo(path);
    }
  }

  static setMember(member) {
    window.sessionStorage.setItem("member", JSON.stringify(member));
  }

  static setToken(token) {
    window.sessionStorage.setItem("token", JSON.stringify(token));
  }

  static tokenExpireFunc(httpResponse, path, navigateTo) {
    if (httpResponse && httpResponse.status === 401) {
      SessionService.clearToken();
      SessionService.navigateOnTokenExpired({ message: ErrTokenExpired }, path, navigateTo);
      return false;
    }

    if (!SessionService.hasToken()) {
      SessionService.navigateOnTokenExpired({ message: ErrTokenExpired }, path, navigateTo);
      return false;
    }

    return true;
  };
}

