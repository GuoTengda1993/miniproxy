<!--
 * @Description: Cert
-->
<template>
<div>
    <el-dialog v-model="data.dialogCertVisble" title="Cert Download" width="500" :before-close="beforeClose" destroy-on-close>
        <el-row justify="center">
            <el-text tag="p" type="primary">
                Scan to download
            </el-text>
        </el-row>
        <el-row justify="center">
            <qrcode-vue :value="url" :size="280" level="H" class="qrcode"></qrcode-vue>
        </el-row>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="data.dialogCertVisble = false">Cancel</el-button>
                <el-button type="primary" @click="save">Save to local</el-button>
            </div>
        </template>
    </el-dialog>
</div>
</template>

<script setup lang="ts">
import {
    reactive,
    defineExpose,
    ref
} from 'vue'
import QrcodeVue from 'qrcode.vue'
import {
    ProxySaveCert,
    ProxyShowCert,
    ProxyCloseCert
} from '../../../wailsjs/go/apps/App'
import {
    Message
} from '../../utils/message'

const url = ref('')

const data = reactive({
    dialogCertVisble: false,
})

const show = () => {
    ProxyShowCert().then(r => {
        url.value = r
        data.dialogCertVisble = true
        return
    })
}
defineExpose({
    show
})

function save() {
    ProxySaveCert().then(r => {
        Message("Save success", r)
        data.dialogCertVisble = false
    })
}

function beforeClose(done) {
    ProxyCloseCert();
    done();
}
</script>

<style scoped>
.qrcode {
    display: block;
    width: 50%;
    height: 50%;
    margin: auto;
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100% 100%;
    background-origin: content-box;
}

.el-text {
    margin-bottom: 5px;
}
</style>
