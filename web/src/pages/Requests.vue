<template>
  <div class="page-shell">
    <div class="header-row">
      <h2>资源求助</h2>
      <el-button type="primary" @click="showCreate = true">发布求助</el-button>
    </div>

    <el-row :gutter="16">
      <el-col :span="24">
        <el-card shadow="never">
          <el-table :data="requests" style="width: 100%">
            <el-table-column prop="title" label="标题" />
            <el-table-column prop="description" label="描述" show-overflow-tooltip />
            <el-table-column prop="bounty" label="悬赏积分" width="120">
              <template #default="scope">
                <div class="bounty-cell">
                  <span class="bounty-val">{{ scope.row.bounty }}</span>
                  <span class="bounty-unit">分</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="user.displayName" label="发布人" width="150">
              <template #default="scope">
                <span class="author-name">{{ scope.row.user?.displayName || 'Unknown' }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 'open' ? 'success' : 'info'" effect="dark" size="small">
                  {{ scope.row.status === 'open' ? '进行中' : '已结束' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="发布时间" width="180">
              <template #default="scope">
                <span class="muted">{{ new Date(scope.row.createdAt).toLocaleString() }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100" fixed="right">
              <template #default="scope">
                <el-button link type="primary" :disabled="scope.row.status !== 'open'">查看</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="showCreate" title="发布求助" width="500px" class="request-dialog">
      <el-form :model="form" label-position="top">
        <el-form-item label="标题">
          <el-input v-model="form.title" placeholder="简短描述你需要什么资源（如：Cisco IOS XR 7.0 镜像）" />
        </el-form-item>
        <el-form-item label="详细描述">
          <el-input v-model="form.description" type="textarea" :rows="4" placeholder="请详细说明资源的版本、适用设备、用途等信息，描述越详细越容易获得帮助。" />
        </el-form-item>
        <el-form-item label="悬赏积分">
          <div class="bounty-input-row">
            <el-input-number v-model="form.bounty" :min="0" :step="10" controls-position="right" style="width: 160px" />
            <div class="balance-info">
              <span class="label">可用积分:</span>
              <span class="value">{{ userStore.profile?.points || 0 }}</span>
            </div>
          </div>
          <div class="form-tip">高额悬赏能吸引更多大佬关注哦！</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreate = false">取消</el-button>
          <el-button type="primary" @click="handleCreate" :loading="submitting">发布</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { fetchRequests, createRequest } from '@/api'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import type { Request } from '@/types'

const userStore = useUserStore()
const requests = ref<Request[]>([])
const showCreate = ref(false)
const submitting = ref(false)

const form = reactive({
  title: '',
  description: '',
  bounty: 0
})

onMounted(loadRequests)

async function loadRequests() {
  try {
    requests.value = await fetchRequests()
  } catch (e) {
    console.error(e)
  }
}

async function handleCreate() {
  if (!form.title) {
    ElMessage.warning('请输入标题')
    return
  }
  if (form.bounty > (userStore.profile?.points || 0)) {
    ElMessage.error('积分不足')
    return
  }

  submitting.value = true
  try {
    await createRequest(form)
    ElMessage.success('发布成功')
    showCreate.value = false
    form.title = ''
    form.description = ''
    form.bounty = 0
    loadRequests()
    if (userStore.profile) {
        userStore.profile.points -= form.bounty
    }
  } catch (e) {
    ElMessage.error('发布失败')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.bounty-cell {
  display: flex;
  align-items: baseline;
  gap: 2px;
}
.bounty-val {
  color: var(--accent);
  font-weight: 700;
  font-size: 16px;
}
.bounty-unit {
  font-size: 12px;
  color: var(--muted);
}
.author-name {
  font-weight: 500;
  color: var(--text);
}
.bounty-input-row {
  display: flex;
  align-items: center;
  gap: 16px;
  width: 100%;
}
.balance-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  background: rgba(255, 255, 255, 0.05);
  padding: 4px 10px;
  border-radius: 4px;
  border: 1px solid var(--border);
}
.balance-info .label {
  color: var(--muted);
}
.balance-info .value {
  color: var(--accent);
  font-weight: 600;
}
.form-tip {
  font-size: 12px;
  color: var(--muted);
  margin-top: 6px;
}
</style>
