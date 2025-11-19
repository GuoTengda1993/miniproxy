import useClipboard from 'vue-clipboard3'
import { ElMessage } from 'element-plus'

const { toClipboard } = useClipboard()

export const copy = async (msg: string) => {
  try {
    await toClipboard(msg)
    ElMessage.success("Copy success")
  } catch (e) {
    ElMessage.error(e)
  }
}