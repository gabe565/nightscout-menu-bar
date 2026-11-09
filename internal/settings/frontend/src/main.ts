import { createApp } from "vue";
import App from "@/App.vue";
import { options } from "@/settings";
import "@/style.css";

document.documentElement.dataset.platform = options.platform;
createApp(App).mount("#app");
