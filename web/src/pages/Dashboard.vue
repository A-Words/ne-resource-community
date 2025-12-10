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
          <h4>收藏与下载历史</h4>
          <el-empty description="功能占位，待后端接口完善" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const { profile } = storeToRefs(userStore)
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
