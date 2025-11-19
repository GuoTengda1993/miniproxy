/*
 * @Author: guotengda guotengda@xiaomi.com
 * @Date: 2025-02-15 15:31:10
 * @LastEditors: guotengda guotengda@xiaomi.com
 * @LastEditTime: 2025-11-06 16:31:06
 * @FilePath: /cartools/frontend/src/main.js
 * @Description: main
 */
import {createApp} from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import JsonViewer from "vue3-json-viewer"
import './style.css';
import "vue3-json-viewer/dist/vue3-json-viewer.css";

const app = createApp(App)
app.use(router)
app.use(ElementPlus)
app.use(JsonViewer)
app.mount('#app')
