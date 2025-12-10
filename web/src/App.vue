<template>
  <el-container>
    <el-header class="header">
      <div class="logo">NE Resource</div>
      <div class="nav-actions">
        <RouterLink to="/">资源广场</RouterLink>
        <RouterLink to="/upload">上传资源</RouterLink>
        <RouterLink to="/dashboard">个人空间</RouterLink>
      </div>
      <div class="user-area">
        <el-button v-if="!userStore.isAuthed" type="primary" @click="router.push('/auth')">登录 / 注册</el-button>
        <div v-else class="user-tag">
          <span>{{ userStore.profile?.displayName }}</span>
          <el-button link type="danger" @click="userStore.logout">退出</el-button>
        </div>
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

const userStore = useUserStore()
const router = useRouter()
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
}
</style>
