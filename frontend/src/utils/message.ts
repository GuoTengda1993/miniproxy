import { ElMessage, ElNotification } from 'element-plus'

export function Message(success: string, error: string) {
    if (error != "") {
        ElMessage.error(error)
    } else if (success != "") {
        ElMessage.success(success)
    }
}

export function Notification(success: string, error: string) {
    let options = {
        title: '成功',
        message: success,
        offset: 80,
    }
    if (error != "") {
        options.title = 'Fail'
        options.message = error
        ElNotification.error(options)
    } else if (success != "") {
        ElNotification.success(options)
    }
}