<template>
  <div class="page-shell">
    <h2>个人学习空间</h2>
    <el-row :gutter="16">
      <el-col :md="8" :xs="24">
        <el-card shadow="never">
          <h4>积分与等级</h4>
          <p class="muted">积分用于奖励贡献与互动</p>
          <div class="metric">{{ profile?.points ?? 0 }} 分</div>
          <div class="muted">等级：{{ profile?.role || 'user' }}</div>
        </el-card>
      </el-col>
      <el-col :md="16" :xs="24">
        <el-card shadow="never">
          <el-tabs v-model="activeTab">
            <el-tab-pane label="我的收藏" name="favorites">
              <el-empty v-if="favorites.length === 0" description="暂无收藏" />
              <el-table v-else :data="favorites" style="width: 100%">
                <el-table-column prop="title" label="标题" />
                <el-table-column prop="type" label="类型" width="100" />
                <el-table-column label="操作" width="100">
                  <template #default="scope">
                    <el-button link type="primary" @click="goDetail(scope.row.id)">查看</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-tab-pane>
            <el-tab-pane label="下载历史" name="downloads">
              <el-empty v-if="downloads.length === 0" description="暂无下载记录" />
              <el-table v-else :data="downloads" style="width: 100%">
                <el-table-column prop="title" label="标题" />
                <el-table-column prop="type" label="类型" width="100" />
                <el-table-column label="操作" width="100">
                  <template #default="scope">
                    <el-button link type="primary" @click="goDetail(scope.row.id)">查看</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-tab-pane>
            <el-tab-pane label="安全设置" name="security">
              <el-form :model="passwordForm" label-width="100px" style="max-width: 400px">
                <el-form-item label="当前密码">
                  <el-input v-model="passwordForm.oldPassword" type="password" show-password />
                </el-form-item>
                <el-form-item label="新密码">
                  <el-input v-model="passwordForm.newPassword" type="password" show-password />
                </el-form-item>
                <el-form-item label="确认新密码">
                  <el-input v-model="passwordForm.confirmPassword" type="password" show-password />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" @click="handleChangePassword">修改密码</el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/stores/user'
import { fetchFavorites, fetchDownloads, changePassword } from '@/api'
import type { Resource } from '@/types'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const { profile } = storeToRefs(userStore)
const router = useRouter()
const route = useRoute()

const activeTab = ref(route.query.tab as string || 'favorites')
const favorites = ref<Resource[]>([])
const downloads = ref<Resource[]>([])
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

onMounted(async () => {
  favorites.value = await fetchFavorites()
  downloads.value = await fetchDownloads()
})

function goDetail(id: string) {
  router.push(`/resource/${id}`)
}

async function handleChangePassword() {
  if (!passwordForm.value.oldPassword || !passwordForm.value.newPassword) {
    ElMessage.error('请填写完整')
    return
  }
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    ElMessage.error('两次输入的密码不一致')
    return
  }
  if (passwordForm.value.newPassword.length < 6) {
    ElMessage.error('新密码长度不能少于6位')
    return
  }
  try {
    await changePassword(passwordForm.value.oldPassword, passwordForm.value.newPassword)
    ElMessage.success('密码修改成功')
    passwordForm.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '修改失败')
  }
}
</script>

<style scoped>
.metric {
  font-size: 32px;
  font-weight: 700;
  color: #22d3ee;
}
.muted {
  color: var(--muted);
}
</style>
