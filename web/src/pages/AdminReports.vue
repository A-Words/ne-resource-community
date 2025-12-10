<template>
  <div class="page-shell">
    <h2>举报处理 (管理员)</h2>
    <el-card shadow="never">
      <el-empty v-if="reports.length === 0" description="暂无待处理举报" />
      <el-table v-else :data="reports" style="width: 100%">
        <el-table-column prop="resource.title" label="被举报资源" />
        <el-table-column prop="user.displayName" label="举报人" width="150" />
        <el-table-column prop="reason" label="举报原因" />
        <el-table-column prop="createdAt" label="举报时间" width="180">
          <template #default="scope">
            {{ new Date(scope.row.createdAt).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="scope">
            <el-button size="small" type="primary" @click="resolve(scope.row.id)">标记已处理</el-button>
            <el-button size="small" link @click="goResource(scope.row.resourceId)">查看资源</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { fetchReports, resolveReport } from '@/api'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'

const reports = ref<any[]>([])
const router = useRouter()

onMounted(load)

async function load() {
  reports.value = await fetchReports()
}

async function resolve(id: string) {
  try {
    await resolveReport(id)
    ElMessage.success('已标记为已处理')
    load()
  } catch (err: any) {
    ElMessage.error(err?.response?.data?.error || '操作失败')
  }
}

function goResource(id: string) {
  router.push(`/resources/${id}`)
}
</script>
