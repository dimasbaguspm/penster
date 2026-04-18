import { createApp } from "vue";
import App from "./app.vue";
import router from "./router";
import "./style.css";

async function enableMocking() {
  if (import.meta.env.DEV) {
    const { worker } = await import("@/mocks/node");
    return worker.start();
  }
  return Promise.resolve();
}

enableMocking().then(() => {
  createApp(App).use(router).mount("#app");
});
