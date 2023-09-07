/** @typedef { object & { http: fetcher, tokenGetterFunction: function, expiredTokenCallback: function, spinner: object, navigateTo: function } } GraphQLOptions */
/**
 * This class is a wrapper around the fetcher function
 * that makes it easy to execute GraphQL queries and mutations.
 * It also handles expired tokens.
 * @class GraphQL
 * @param {string} queryURL The URL to the GraphQL API
 * @param {GraphQLOptions} options Options for the GraphQL class
 */
export class GraphQL {
    constructor(queryURL: any, options?: {
        http: typeof fetcher;
        tokenGetterFunction: any;
        expiredTokenCallback: any;
        spinner: any;
        navigateTo: any;
    });
    queryURL: any;
    http: typeof fetcher;
    tokenGetterFunction: any;
    expiredTokenCallback: any;
    spinner: any;
    navigateTo: any;
    /**
     * Executes a query against a GraphQL API
     * @param query string A graphql query. Omit the "query {}" portion.
     * @returns {Promise<object>} A promise that resolves to the fetch response
     */
    query(query: any): Promise<object>;
    /**
     * Executes a mutation against a GraphQL API
     * @param query string A graphql mutation. Omit the "mutation {}" portion
     * @returns {Promise<object>} A promise that resolves to the fetch response
     */
    mutation(query: any): Promise<object>;
}
export type GraphQLOptions = object & {
    http: typeof fetcher;
    tokenGetterFunction: Function;
    expiredTokenCallback: Function;
    spinner: object;
    navigateTo: Function;
};
import { fetcher } from "./fetcher.js";
