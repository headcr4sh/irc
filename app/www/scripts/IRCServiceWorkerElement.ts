export type Status = "not-started" | "ok" | "error";

export class IRCServiceWorkerElement extends HTMLElement {
    public get status(): Status {
        return this.getAttribute("status") as Status || "not-started";
    }

    public set status(status: Status) {
        this.setAttribute("status", status);
    }

    public get scope(): string {
        return this.getAttribute("scope") || "";
    }

    public set scope(scope: string) {
        this.setAttribute("scope", scope);
    }

    constructor() {
        super();
    }
}

window.customElements.define("irc-service-worker", IRCServiceWorkerElement);
