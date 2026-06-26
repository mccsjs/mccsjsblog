import { createApp } from "vue";
import zhCn from "element-plus/es/locale/lang/zh-cn";
import ElementPlus from 'element-plus';
import App from "./App.vue";
import router from "./router";
import "./assets/css/main.scss";

const app = createApp(App);

app.use(router);
app.use(ElementPlus, { locale: zhCn });

app.mount("#app");
