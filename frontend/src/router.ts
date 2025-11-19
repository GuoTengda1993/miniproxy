/*
 * @Author: guotengda guotengda@xiaomi.com
 * @Date: 2025-05-20 19:12:18
 * @LastEditors: guotengda guotengda@xiaomi.com
 * @LastEditTime: 2025-11-06 15:16:10
 * @FilePath: /cartools/frontend/src/router.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const Proxy = ()=> import('./views/proxy/Proxy.vue')


const routes: Array<RouteRecordRaw> = [
    {
        path: '/',
        name: 'proxy',
        component: Proxy,
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

export default router