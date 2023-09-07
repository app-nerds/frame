export const ErrTokenExpired: "token expired";
/**
 * SessionService is a service used to manage session data. It makes use
 * of the browser's sessionStorage object. This class provides
 * only static methods.
 * @class SessionService
 */
export default class SessionService {
    /**
     * Clears the member from session storage.
     * @static
     * @memberof SessionService
     * @returns {void}
     */
    static clearMember(): void;
    /**
     * Clears the token from session storage.
     * @static
     * @memberof SessionService
     * @returns {void}
     */
    static clearToken(): void;
    /**
     * Gets the member from session storage.
     * @static
     * @memberof SessionService
     * @returns {object} The member object
     */
    static getMember(): object;
    /**
     * Gets the token from session storage.
     * @static
     * @memberof SessionService
     * @returns {object} The token object
     */
    static getToken(): object;
    /**
     * Determines if the session has a member.
     * @static
     * @memberof SessionService
     * @returns {boolean} True if the session has a member, otherwise false.
     */
    static hasMember(): boolean;
    /**
     * Determines if the session has a token.
     * @static
     * @memberof SessionService
     * @returns {boolean} True if the session has a token, otherwise false.
     */
    static hasToken(): boolean;
    /**
     * Navigates to the specified path if the token has expired.
     * This is determined by examing the error message. If the
     * error message is "token expired", the user is navigated
     * to the specified path.
     * @static
     * @param {object} e The error object
     * @param {string} path The path to navigate to
     * @param {function} navigateTo The function to use to navigate
     * @memberof SessionService
     */
    static navigateOnTokenExpired(e: object, path: string, navigateTo: Function): void;
    /**
     * Sets the member in session storage.
     * @static
     * @param {object} member The member object
     * @memberof SessionService
     * @returns {void}
     */
    static setMember(member: object): void;
    /**
     * Sets the token in session storage.
     * @static
     * @param {object} token The token object
     * @memberof SessionService
     * @returns {void}
     */
    static setToken(token: object): void;
    /**
     * This is a function that can be used as a callback for the fetcher
     * class. It will check the response for a 401 status code. If it
     * finds one, it will clear the token and navigate to the specified
     * path.
     * @static
     * @param {object} httpResponse The HTTP response object
     * @param {string} path The path to navigate to
     * @param {function} navigateTo The function to use to navigate
     * @memberof SessionService
     * @returns {boolean} True if the token is valid, otherwise false.
     */
    static tokenExpireFunc(httpResponse: object, path: string, navigateTo: Function): boolean;
}
