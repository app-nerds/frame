export function application(targetElement: HTMLElement, routes: Array<Route>, pageNotFoundView?: BaseView): Application;
export type Route = import("./router.js").Route;
export type Application = object & {
    routes: Array<Route>;
    targetElement: HTMLElement;
    router: Router;
    afterRoute: Function;
    beforeRoute: Function;
    injectParams: Function;
    onRenderComplete: Function;
    go: Function;
};
import { Router } from "./router.js";
