<template>
<div>
    <el-dialog v-model="data.dialogHeadersVisble" title="" width="800" destroy-on-close>
        <el-row justify="center">
            <el-descriptions :title="title" :border="true" :column="1">
                <template v-for="hVal, hKey in props.headers">
                    <el-descriptions-item :label="hKey">
                        <el-button :icon="DocumentCopy" link type="primary" @click="copy(hVal)" />
                        <div class="ellipsis">{{ hVal }}</div>
                    </el-descriptions-item>
                </template>
            </el-descriptions>
        </el-row>
    </el-dialog>
</div>
</template>

<script setup lang="ts">
import {
    computed,
    reactive,
    defineExpose
} from 'vue'
import {
    DocumentCopy
} from '@element-plus/icons-vue'
import {
    copy
} from '../../utils/copy'

const props = defineProps({
    headers: {},
    isRequest: Boolean,
})

const title = computed(() => {
    return props.isRequest ? 'Request Headers' : 'Response Headers'
})

const show = () => {
    data.dialogHeadersVisble = true
}
defineExpose({
    show
})

const data = reactive({
    dialogHeadersVisble: false,
})
</script>

<style>
.ellipsis {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    width: 350px;
}

.el-descriptions {
    padding: 3px 0px 0px 0px;
    text-align: center;
}
</style>
