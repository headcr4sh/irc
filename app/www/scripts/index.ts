import {IRCServiceWorkerElement} from "./IRCServiceWorkerElement";

(async () => {
    const swElement = document.querySelector("irc-service-worker") as IRCServiceWorkerElement;
    try {
        const registration = await navigator.serviceWorker.register('./sw.js', {scope: './'});
        swElement.scope = registration.scope;
        swElement.status = "ok";
    } catch (e) {
        swElement.status = "error"
    }
})();
