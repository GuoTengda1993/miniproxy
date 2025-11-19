<template>
    <el-container>
        <el-header>
            <el-menu :default-active="activeIndex" mode="horizontal" :ellipsis="false" :router="true">
                <el-menu-item index="">
                    <img style="width: 45px" src="./assets/images/logo-universal.png" alt="miniproxy" />
                </el-menu-item>
                <el-menu-item index="/">Proxy</el-menu-item>
            </el-menu>
        </el-header>
        <el-main>
            <RouterView :key="$route.fullPath" />
        </el-main>
    </el-container>
</template>

<script setup>
import {
    ref,
    onMounted,
    onBeforeUnmount
} from 'vue'
import {
    ElMessage,
    ElMessageBox
} from 'element-plus'
import {
    FrontPause,
    FrontResume,
    LogPrint
} from '../wailsjs/go/apps/App'
import {
    EventsOn,
    EventsOff
} from '../wailsjs/runtime/runtime'

const activeIndex = ref('/')

onMounted(() => {
    EventsOn("sendMessage", (msg) => {
        if (msg == "" || msg == undefined || msg == null) {
            return
        }
        if (msg.use_box) {
            ElMessageBox.alert(msg.message, msg.title, {
                confirmButtonText: 'OK',
                dangerouslyUseHTMLString: true,
            })
            return
        }
        ElMessage({
            message: msg.message,
            type: msg.type,
        })
    })

    if (typeof document.visibilityState !== 'undefined') {
        document.addEventListener('visibilitychange', function () {
            if (document.visibilityState === 'visible') {
                FrontResume();
            } else {
                FrontPause();
            }
        });
    } else {
        window.addEventListener('focus', function () {
            LogPrint('窗口激活，可能在前台');
            FrontResume();
        });
        window.addEventListener('blur', function () {
            LogPrint('窗口失活，可能在后台');
            FrontPause();
        });
    }

})

onBeforeUnmount(() => {
    EventsOff('sendMessage');
})
</script>

<style scoped>
.el-menu--horizontal>.el-menu-item:nth-child(1) {
    margin-right: auto;
}

.el-header {
    padding: 0;
    top: 0;
    width: 100%;
    position: fixed;
    z-index: 1000;
}

.el-main {
    margin-top: 60px;
}
</style>
