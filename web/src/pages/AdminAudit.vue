<template>
  <div class="page-shell">
    <h2>资源审核 (管理员)</h2>
    <el-card shadow="never">
      <el-empty v-if="pendingList.length === 0" description="暂无待审核资源" />
      <el-table v-else :data="pendingList" style="width: 100%">
        <el-table-column prop="title" label="标题" />
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column prop="uploader.displayName" label="上传者" width="150" />
        <el-table-column prop="createdAt" label="上传时间" width="180">
          <template #default="scope">
            {{ new Date(scope.row.createdAt).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button size="small" type="success" @click="audit(scope.row.id, 'approve')">通过</el-button>
            <el-button size="small" type="danger" @click="openReject(scope.row.id)">拒绝</el-button>
            <el-button size="small" link @click="download(scope.row.id)">下载预览</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="showReject" title="拒绝原因" width="30%">
      <el-input v-model="rejectReason" type="textarea" placeholder="请输入拒绝原因" />
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showReject = false">取消</el-button>
          <el-button type="primary" @click="confirmReject">确认拒绝</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { fetchPendingResources, auditResource, downloadResource } from '@/api'
import type { Resource } from '@/types'
import { ElMessage } from 'element-plus'

const pendingList = ref<Resource[]>([])
const showReject = ref(false)
const rejectReason = ref('')
const currentId = ref('')

onMounted(load)

async function load() {
  pendingList.value = await fetchPendingResources()
}

async function audit(id: string, action: 'approve' | 'reject', reason?: string) {
  try {
    await auditResource(id, action, reason)
    ElMessage.success(action === 'approve' ? '已通过' : '已拒绝')
    load()
  } catch (err: any) {
    ElMessage.error(err?.response?.data?.error || '操作失败')
  }
}

function openReject(id: string) {
  currentId.value = id
  rejectReason.value = ''
  showReject.value = true
}

async function confirmReject() {
  if (!rejectReason.value) {
    ElMessage.warning('请输入拒绝原因')
    return
  }
  await audit(currentId.value, 'reject', rejectReason.value)
  showReject.value = false
}

async function download(id: string) {
  try {
    await downloadResource(id)
  } catch (err) {
    ElMessage.error('下载失败')
  }
}
</script>
