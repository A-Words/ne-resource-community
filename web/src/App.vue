<template>
  <el-container>
    <el-header class="header">
      <div class="logo">NE Resource</div>
      <div class="nav-actions">
        <RouterLink to="/">资源广场</RouterLink>
        <RouterLink to="/upload">上传资源</RouterLink>
      </div>
      <div class="user-area">
        <el-button v-if="!userStore.isAuthed" type="primary" @click="router.push('/auth')">登录 / 注册</el-button>
        <el-dropdown v-else trigger="click" @command="handleCommand">
          <div class="user-tag">
            <el-avatar :size="32" :icon="UserFilled" class="user-avatar" />
            <span class="user-name">{{ userStore.profile?.displayName }}</span>
            <el-tag size="small" type="warning" effect="dark" round>Lv.{{ userStore.profile?.level || 1 }}</el-tag>
            <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="/dashboard" :icon="DataLine">个人空间</el-dropdown-item>
              <el-dropdown-item command="/dashboard?tab=security" :icon="Lock">安全设置</el-dropdown-item>
              <template v-if="userStore.profile?.role === 'admin'">
                <el-dropdown-item command="/admin/audit" :icon="Monitor">审核后台</el-dropdown-item>
                <el-dropdown-item command="/admin/reports" :icon="Monitor">举报处理</el-dropdown-item>
              </template>
              <el-dropdown-item divided command="logout" :icon="SwitchButton">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>
    <el-main>
      <RouterView />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { useUserStore } from './stores/user'
import { ArrowDown, UserFilled, SwitchButton, DataLine, Monitor, Lock } from '@element-plus/icons-vue'

const userStore = useUserStore()
const router = useRouter()

const handleCommand = (command: string) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/auth')
  } else {
    router.push(command)
  }
}
</script>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  background: rgba(17, 24, 39, 0.6);
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  backdrop-filter: blur(8px);
  z-index: 10;
}
.logo {
  font-weight: 700;
  letter-spacing: 0.6px;
  color: var(--accent);
}
.nav-actions {
  display: flex;
  gap: 16px;
  font-weight: 600;
}
.nav-actions a {
  color: var(--text);
}
.nav-actions a.router-link-active {
  color: var(--accent);
}
.user-area {
  display: flex;
  align-items: center;
  gap: 12px;
}
.user-tag {
  display: flex;
  gap: 8px;
  align-items: center;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: background-color 0.2s;
}
.user-tag:hover {
  background: rgba(255, 255, 255, 0.05);
}
.user-name {
  font-weight: 600;
  color: var(--text);
}
.user-avatar {
  background: var(--accent);
  color: #0f172a;
}
</style>
