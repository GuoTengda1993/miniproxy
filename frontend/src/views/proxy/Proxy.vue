<template>
<main>
    <el-row type="flex" :gutter="10" style="display: flex; align-items: center;">
        <el-button style="margin-right: auto;" v-if="data.isRunning == false" type="primary" @click="start">Start</el-button>
        <el-button style="margin-right: auto;" v-if="data.isRunning" type="danger" @click="stop">Stop</el-button>
        <el-tag v-if="data.isRunning" size="large" type="success" style="margin-left: 10px;">{{ data.addr }}</el-tag>

        <div style="flex-grow: 1; display: flex; justify-content: right;">
            <el-select v-if="remoteAddrs.length > 1" v-model="selectRemoteAddr" placeholder="设备筛选" style="width: 150px; margin-right: 10px;" @change="queryFlowAll">
                <el-option key="" label="全部" value="" />
                <el-option v-for="item in remoteAddrs" :key="item" :label="item" :value="item" />
            </el-select>
            <el-select v-if="catchHosts.length > 1" v-model="selectHost" placeholder="Host筛选" style="width: 150px; margin-right: 10px;" @change="queryFlowAll">
                <el-option key="" label="全部" value="" />
                <el-option v-for="item in catchHosts" :key="item" :label="item" :value="item" />
            </el-select>

            <el-dropdown :disabled="data.isRunning == false">
                <el-button type="primary" :icon="Setting" :disabled="data.isRunning == false" style="margin: 0px 5px;">
                    Setting<el-icon class="el-icon--right">
                        <arrow-down />
                    </el-icon>
                </el-button>
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item>
                            <el-button text @click="dialogHostShow">Host Setting</el-button>
                        </el-dropdown-item>
                        <el-dropdown-item>
                            <el-button text @click="dialogCertShow">Cert Download</el-button>
                        </el-dropdown-item>
                    </el-dropdown-menu>
                </template>
            </el-dropdown>
        </div>
    </el-row>

    <template v-if="data.isRunning">
        <el-row align="middle" style="margin-top: 15px;" :gutter="10">
            <el-col :span="17">
                <el-input v-model="filter" placeholder="Filter by path or http code" clearable @input="flowFilter" />
            </el-col>
            <el-col :span="2">
                <el-checkbox size="small" border v-model="filterOnlyError" label="Show only error" @change="queryFlowAll" style="background-color: rgba(0,191,255,0.1);" />
            </el-col>
            <el-col :span="5" style="flex-grow: 1; display: flex; justify-content: center; align-items:center">
                <el-button type="primary" plain size="small" @click="collapseDetail">Collapse detail</el-button>
                <el-button type="danger" plain size="small" @click="clearFlowList">Clear</el-button>
            </el-col>
        </el-row>
        <el-row>
            <el-col :span="24" class="table-wrapper" @mouseenter="moustInTable = true" @mouseleave="moustInTable = false">
                <el-table :data="flowList" @row-contextmenu="handleRightClick" @row-click="showFlowInfo" :height="tableHeight" style="width: 100%" highlight-current-row :current-row-key="currRowId" scrollbar-always-on row-key="id" ref="tableRef">
                    <el-table-column prop="host" label="Host" width="280" />
                    <el-table-column prop="path" label="Path" />
                    <el-table-column prop="method" label="Method" width="90" />
                    <el-table-column prop="code" label="HttpCode" width="120">
                        <template #default="scope">
                            <el-tag v-if="scope.row.code >= 200 && scope.row.code < 300" type="success">{{ scope.row.code }}</el-tag>
                            <el-tag v-else-if="scope.row.code == 0" type="info">{{ scope.row.code }}</el-tag>
                            <el-tag v-else type="danger">{{ scope.row.code }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="return_code" label="ReturnCode" width="120">
                        <template #default="scope">
                            <div v-if="scope.row.return_code != ''">
                                <el-tag v-if="scope.row.return_code == '200' || scope.row.return_code == '0'" type="success">{{ scope.row.return_code }}</el-tag>
                                <el-tag v-else type="danger">{{ scope.row.return_code }}</el-tag>
                            </div>
                        </template>
                    </el-table-column>
                    <el-table-column prop="response_time_str" label="ResponseTime" width="130">
                        <template #default="scope">
                            <div v-if="scope.row.response_time_str != ''">{{ scope.row.response_time_str.split(" ")[1] }}</div>
                            <div v-else>{{ scope.row.response_time_str }}</div>
                        </template>
                    </el-table-column>
                </el-table>
                <el-backtop target=".table-wrapper .el-scrollbar__wrap" :visibility-height="70" :right="15" :bottom="15"></el-backtop>
            </el-col>
        </el-row>
        <el-row>
            <div style="flex-grow: 1; display: flex; justify-content: right; margin-top: 5px;">
                <el-pagination v-model:current-page="data.currentPage" :page-size="pageSize" :background="true" size="small" layout="total, prev, pager, next" :total="data.total" @current-change="handlePageChange" />
            </div>
        </el-row>
    </template>

    <template v-if="data.showFlowDetail && data.isRunning">
        <el-tabs v-model="tabActiveName" type="border-card">
            <!-- request detail -->
            <el-tab-pane label="Request Info" name="req">
                <el-descriptions :border="true" :column="2">
                    <template #extra>
                        <el-button type="primary" link @click="dialogHeadersShow(true)">Headers</el-button>
                        <el-button type="primary" link @click="copyFlowInfo('curl')">Copy Curl</el-button>
                    </template>
                    <el-descriptions-item label="Host">
                        {{ data.flowDetail['scheme'] }}://{{ data.flowDetail['host'] }}
                    </el-descriptions-item>
                    <el-descriptions-item label="Path">
                        <el-button :icon="DocumentCopy" link type="primary" @click="copyFlowInfo('path')" />
                        {{ data.flowDetail['path'] }}
                    </el-descriptions-item>
                    <el-descriptions-item label="Method">{{ data.flowDetail['method'] }}</el-descriptions-item>
                    <el-descriptions-item label="Size">{{ data.flowDetail['req_content_length'] }}</el-descriptions-item>
                    <el-descriptions-item v-if="data.flowDetail['query'] != ''" label="Query">
                        <el-button :icon="DocumentCopy" link type="primary" @click="copyFlowInfo('query')" />
                        <el-tooltip :content="data.flowDetail['query']" placement="top">
                            <div class="ellipsis">{{ data.flowDetail['query'] }}</div>
                        </el-tooltip>
                    </el-descriptions-item>
                    <el-descriptions-item v-if="data.flowDetail['method'] != 'GET'" label="Body">
                        <div class="json-box">
                            <json-viewer :value="JSON.parse(data.flowDetail['data'])" copyable boxed sort expanded :expand-depth="5" />
                        </div>
                    </el-descriptions-item>
                </el-descriptions>
            </el-tab-pane>
            <!-- Response Info -->
            <el-tab-pane label="Response Info" name="resp">
                <el-descriptions :border="true">
                    <template #extra>
                        <el-button type="primary" link @click="dialogHeadersShow(false)">Headers</el-button>
                        <el-button type="primary" link @click="copyResponse">Copy Data</el-button>
                    </template>
                    <el-descriptions-item label="Http Code">{{ data.flowDetail['code'] }}</el-descriptions-item>
                    <el-descriptions-item label="Response Time">{{ data.flowDetail['response_time_str'] }}</el-descriptions-item>
                    <el-descriptions-item label="Cost Time">{{ data.flowDetail['duration'] }}ms</el-descriptions-item>
                    <el-descriptions-item v-if="data.flowDetail['resp_content_length'] != ''" label="Size">{{ data.flowDetail['resp_content_length'] }}</el-descriptions-item>
                </el-descriptions>
                <div class="json-box">
                    <Vue3JsonEditor v-model="data.flowResponse" :show-btns="false" :expandedOnStart="false" mode="view" :modes="jsonModeList" />
                </div>
            </el-tab-pane>
        </el-tabs>
    </template>

    <!-- right click menu -->
    <div v-show="contextMenu.visible" class="context-menu" ref="menuRef" :style="{
            left: contextMenu.x + 'px',
            top: contextMenu.y + 'px',
            position: 'fixed'
          }">
        <div class="menu-item" @click="copyURL">Copy URL</div>
        <div class="menu-item" @click="copyFlowInfo('curl')">Copy Curl</div>
        <div class="menu-item" @click="copyHeader('Cookie')">Copy Cookie</div>
        <div class="menu-item" @click="copyResponse">Copy Response</div>
        <el-divider style="margin-top: 2px; margin-bottom: 2px;"></el-divider>
        <div class="menu-item" @click="dialogHeadersShow(true)">Req Headers</div>
        <div class="menu-item" @click="dialogHeadersShow(false)">Resp Headers</div>
    </div>

    <DialogCert ref="dialogCertRef" />
    <DialogHeaders ref="dialogHeadersRef" :headers="dialogHeadersInfo" :isRequest="dialogHeadersIsRequest" />
    <DialogHost ref="dialogHostRef" />
</main>
</template>

<script setup lang="ts">
import {
    ref,
    reactive,
    onMounted,
    onBeforeUnmount
} from 'vue'
import {
    Setting,
    DocumentCopy,
    ArrowDown,
} from '@element-plus/icons-vue'
import type {
    TableInstance
} from 'element-plus'
import {
    ProxyStart,
    ProxyStop,
    ProxyInfo,
    ProxyFlowListByPage,
    ProxyFlowClear,
    WindowGetHeight,
    WindowGetWidth,
    ProxyRemoteAddrs,
    ProxyCatchHosts,
    LogPrint,
} from '../../../wailsjs/go/apps/App'
import {
    EventsOn,
    EventsOff
} from '../../../wailsjs/runtime/runtime'
import {
    Vue3JsonEditor
} from 'vue3-json-editor'
import {
    copy
} from '../../utils/copy'
import DialogHost from './DialogHost.vue'
import DialogHeaders from './DialogHeaders.vue'
import DialogCert from './DialogCert.vue'

const tabActiveName = ref('req')
const pageSize = ref(50)

const flowList = ref([])

const tableHeight = ref(400)

const tableRef = ref < TableInstance > ()

const dialogHostRef = ref()
const dialogHostShow = () => {
    dialogHostRef.value.show()
}

const dialogHeadersRef = ref()
const dialogHeadersInfo = ref()
const dialogHeadersIsRequest = ref(true)
const dialogHeadersShow = (isRequest: boolean) => {
    if (isRequest) {
        dialogHeadersIsRequest.value = true
        dialogHeadersInfo.value = data.flowDetail['headers']
    } else {
        dialogHeadersIsRequest.value = false
        dialogHeadersInfo.value = data.flowDetail['response_headers']
    }
    dialogHeadersRef.value.show()
}

const dialogCertRef = ref()
const dialogCertShow = () => {
    dialogCertRef.value.show()
}

const menuRef = ref(null)
const contextMenu = reactive({
    visible: false,
    x: 0,
    y: 0,
})

const currRowId = ref(0)

const jsonModeList = ["view", "code"]

const moustInTable = ref(false)

const data = reactive({
    isRunning: false,
    addr: '',
    showFlowDetail: false,
    flowDetail: {},
    flowResponse: {},

    currentPage: 1,
    total: 0,
})

const filter = ref('')
const filterOnlyError = ref(false)

const selectRemoteAddr = ref('')
const remoteAddrs = ref < string[] > ([])

const selectHost = ref('')
const catchHosts = ref < string[] > ([])

const windowHeight = ref(0)
const windowWidth = ref(0)

onMounted(() => {
    document.addEventListener('click', closeContextMenu)
    document.addEventListener('contextmenu', closeContextMenu)
    loadProxyInfo();
    queryFlowAll();
    WindowGetHeight().then(h => {
        windowHeight.value = h
        tableHeight.value = h - 230
    })
    WindowGetWidth().then(w => {
        windowWidth.value = w
    })
    EventsOn("newRemoteAddr", (_) => {
      ProxyRemoteAddrs().then(r => {
        if (r.length > remoteAddrs.value.length) {
          remoteAddrs.value = r
        }
      });
    })
    EventsOn("newHost", (_) => {
      ProxyCatchHosts().then(r => {
        if (r.length > catchHosts.value.length) {
          catchHosts.value = r
        }
      });
    })
    EventsOn("reloadFlowList", (_) => {
        queryFlowAll()
        if (moustInTable.value == false) {
            tableRef.value?.scrollTo(0, 0)
        }
    })
})

onBeforeUnmount(() => {
    EventsOff('newRemoteAddr');
    EventsOff('newHost');
    EventsOff('reloadFlowList');
    document.removeEventListener('click', closeContextMenu);
    document.removeEventListener('contextmenu', closeContextMenu);
    dialogCertRef.value = null;
    dialogHeadersRef.value = null;
    dialogHostRef.value = null;
    flowList.value = [];
})

const handleRightClick = (row, column, event) => {
    event.preventDefault();
    event.stopPropagation();
    contextMenu.x = event.clientX
    contextMenu.y = event.clientY

    let hGap = 170
    if (event.clientY + hGap > windowHeight.value) {
        contextMenu.y = event.clientY - hGap
    }
    let wGap = 120
    if (event.clientX + wGap > windowWidth.value) {
        contextMenu.x = event.clientX - wGap
    }
    setCurrentFlow(row)
    contextMenu.visible = true
}

const closeContextMenu = () => {
    contextMenu.visible = false
}

const handlePageChange = (val: number) => {
    data.currentPage = val
    ProxyFlowListByPage(filter.value, selectRemoteAddr.value, selectHost.value, val, filterOnlyError.value).then(r => {
        flowList.value = r['list']
        data.total = r['total']
    })
    collapseDetail()
}

function queryFlowAll() {
    ProxyFlowListByPage(filter.value, selectRemoteAddr.value, selectHost.value, data.currentPage, filterOnlyError.value).then(r => {
        if (r['total'] == -1) {
            return
        }
        if (data.currentPage <= 1) {
            flowList.value = r['list']
        }
        data.total = r['total']
    })
}

function clearFlowList() {
    ProxyFlowClear().then(_ => {
        flowList.value = []
        collapseDetail()
    })
    data.currentPage = 1
    data.total = 0
    selectRemoteAddr.value = ""
    remoteAddrs.value = []
    catchHosts.value = []
}

function loadProxyInfo() {
    ProxyInfo().then(result => {
        data.isRunning = result['isRunning']
        if (result['isRunning']) {
            data.addr = result['addr']
        }
    })
}

function start() {
    ProxyStart().then(_ => {
        ProxyInfo().then(result => {
            data.isRunning = result['isRunning']
            data.addr = result['addr']
        })
    })
}

function stop() {
    ProxyStop().then(_ => {
        data.isRunning = false
        flowList.value = []
        data.showFlowDetail = false
        data.flowDetail = {}
        data.flowResponse = {}
    })
}

function showFlowInfo(row, column) {
    if (contextMenu.visible) {
        contextMenu.visible = false
        return
    }
    setCurrentFlow(row)
    data.showFlowDetail = true
    tableHeight.value = 245
}

function setCurrentFlow(row) {
    currRowId.value = row.id
    data.flowResponse = {}
    data.flowDetail = row
    try {
        data.flowResponse = JSON.parse(row.response_body)
    } catch (error) {
        LogPrint(error)
        data.flowResponse = row.response_body
    }
}

function copyFlowInfo(key) {
    copy(data.flowDetail[key])
}

function copyHeader(key) {
    copy(data.flowDetail['headers'][key])
}

function copyURL() {
    let url = data.flowDetail['scheme'] + "://" + data.flowDetail['host'] + data.flowDetail['path']
    copy(url)
}

function copyResponse() {
    try {
        let respObj = JSON.parse(data.flowDetail['response_body'])
        let respStr = JSON.stringify(respObj, null, 4)
        copy(respStr)
    } catch (error) {
        LogPrint(error)
        copy(data.flowDetail['response_body'])
    }
}

function flowFilter(value) {
    data.currentPage = 1
    filter.value = value.toLowerCase()
    queryFlowAll()
}

function collapseDetail() {
    data.showFlowDetail = false
    WindowGetHeight().then(h => {
        tableHeight.value = h - 230
    })
}
</script>

<style scoped>
.el-row {
    margin-top: 3px;
    margin-bottom: 3px;
}

.min-latency {
    margin-right: 10px;
}

.divider {
    margin-left: 10px;
    margin-right: 10px;
    width: 2px;
    height: 30px;
    background: #d0cbcb;
}

.form-hint {
    font-size: 12px;
    color: #909399;
    margin-top: 2px;
}

.el-descriptions {
    /* background-color: rgb(215, 237, 245); */
    padding: 3px 0px 0px 0px;
    text-align: center;
}

.json-box {
    text-align: left;
    background-color: white;
}

.ellipsis {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    width: 350px;
}

.table-wrapper {
    position: relative;
}

.el-backtop {
    position: absolute;
}

.el-text {
    text-align: left;
    margin: 10px 15px 10px 15px;
}

.context-menu {
    position: fixed;
    background: #fff;
    border: 1px solid #ebeef5;
    border-radius: 4px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    z-index: 9999;
}

.menu-item {
    padding: 6px 16px;
    cursor: pointer;
    font-size: 12px;
    color: #606266;
    display: flex;
}

.menu-item:hover {
    background: #f5f7fa;
    color: #409eff;
}
</style>
