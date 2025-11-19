<template>
<div>
    <el-dialog v-model="data.dialogVisble" title="Host Setting" width="600" @open="onOpen" destroy-on-close>
        <div class="top-lines">
            <el-switch style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949; float: left;" v-model="data.hostValidSwitch" class="ml-2" inline-prompt active-text="Support" inactive-text="Not Support" @change="switchHostType" />
            <p style="text-align: left; margin-left: 10px;">{{ data.hostHelpMsg }}</p>
        </div>

        <el-table :data="hostList" height="300">
            <el-table-column prop="value" label="Host" />
            <el-table-column fixed="right" label="Operate" width="90">
                <template #default="scope">
                    <el-button v-if="data.hostValidSwitch" link type="danger" size="small" @click="removeHost(scope.row.value)">Remove</el-button>
                    <el-button v-if="data.hostValidSwitch == false" link type="primary" size="small" @click="addHost(scope.row.value)">Add</el-button>
                </template>
            </el-table-column>
        </el-table>
        <template #footer>
            <div class="dialog-footer">
                <el-button v-if="data.hostValidSwitch" type="primary" @click="addHostMsgBox">Add</el-button>
                <el-button type="warning" @click="resetHost">Add</el-button>
                <el-button @click="data.dialogVisble = false">Cancel</el-button>
            </div>
        </template>
    </el-dialog>
</div>
</template>

<script setup lang="ts">
import {
    ref,
    reactive,
    defineExpose
} from 'vue'
import {
    ElMessage,
    ElMessageBox
} from 'element-plus'
import {
    ProxyAddHost,
    ProxyRemoveHost,
    ProxyResetHost,
    ProxyHostList,
} from '../../../wailsjs/go/apps/App'

const show = () => {
    data.dialogVisble = true
}
defineExpose({
    show
})

const hostList = ref([])

const data = reactive({
    dialogVisble: false,
    hostValidSwitch: true,
    hostHelpMsg: '',
})

function onOpen() {
    data.hostValidSwitch = true
    switchHostType()
}

function removeHost(h) {
    ProxyRemoveHost(h).then(_ => {
        switchHostType()
        ElMessage.success("Remove success")
    })
}

function resetHost() {
    ProxyResetHost().then(_ => {
        switchHostType()
        ElMessage.success("Reset success")
    })
}

const addHostMsgBox = () => {
    ElMessageBox.prompt('Please input host(re start with #)', 'Add Host', {
            confirmButtonText: 'Confirm',
            cancelButtonText: 'Cancel',
        })
        .then(({
            value
        }) => {
            ProxyAddHost(value).then(_ => {
                switchHostType();
                ElMessage.success("Add success")
            })
        })
}

function addHost(h) {
    ProxyAddHost(h).then(_ => {
        switchHostType();
        ElMessage.success("Add success")
    })
}

function switchHostType() {
    ProxyHostList(data.hostValidSwitch).then(r => {
        hostList.value = r as any;
    })

    if (data.hostValidSwitch == true) {
        data.hostHelpMsg = 'Packet capture is supported for all hosts in the following list. A "#" at the beginning indicates a regular expression.'
    } else {
        data.hostHelpMsg = 'Packet capture is not supported on the following hosts. You can click "Add" to enable it.'
    }
}
</script>

<style scoped>
.top-lines {
    display: flex;
    align-items: center;
    margin-left: 8px;
}
</style>
