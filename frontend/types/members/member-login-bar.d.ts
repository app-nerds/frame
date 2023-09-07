export default class MemberLoginBar extends HTMLElement {
    static get observedAttributes(): string[];
    memberService: MemberService;
    loginPath: string;
    attributedChangedCallback(name: any, _: any, newValue: any): void;
    connectedCallback(): Promise<void>;
    /*******************************************************************************
     * Event methods
     ******************************************************************************/
    /*******************************************************************************
     * UI elements
     ******************************************************************************/
    createContainerEl(): HTMLDivElement;
    createAvatarEl(container: any, member: any): void;
    createTextEl(container: any, member: any): void;
    createPopupMenu(container: any, member: any): void;
    /*******************************************************************************
     * Private methods
     ******************************************************************************/
    onMyAccountClick(): void;
    onLogOutClick(): void;
}
import { MemberService } from "./member-service.js";
