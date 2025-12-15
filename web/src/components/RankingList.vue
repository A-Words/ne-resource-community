<template>
  <el-card class="ranking-card" shadow="never">
    <template #header>
      <div class="card-header">
        <span>下载排行榜</span>
      </div>
    </template>
    <div v-for="(item, index) in resources" :key="item.id" class="ranking-item" @click="goDetail(item.id)">
      <div class="rank-num" :class="{ 'top-3': index < 3 }">{{ index + 1 }}</div>
      <div class="rank-info">
        <div class="rank-title" :title="item.title">{{ item.title }}</div>
        <div class="rank-count">{{ item.downloadCount }} 次下载</div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { fetchResources } from '@/api'
import type { Resource } from '@/types'
import { useRouter } from 'vue-router'

const router = useRouter()
const resources = ref<Resource[]>([])

onMounted(async () => {
  resources.value = await fetchResources({ sort: 'downloads', limit: 10 })
})

function goDetail(id: string) {
  router.push(`/resources/${id}`)
}
</script>

<style scoped>
.ranking-card {
  margin-bottom: 20px;
}
:deep(.el-card__header) {
  border-bottom: 1px solid var(--border);
  padding: 12px 20px;
}
:deep(.el-card__body) {
  padding: 0 20px;
}
.card-header {
  font-weight: bold;
  color: var(--text);
}
.ranking-item {
  display: flex;
  align-items: center;
  padding: 12px 0;
  cursor: pointer;
  border-bottom: 1px solid var(--border);
  transition: background-color 0.2s;
}
.ranking-item:last-child {
  border-bottom: none;
}
.ranking-item:hover {
  background-color: rgba(255, 255, 255, 0.02);
}
.rank-num {
  width: 24px;
  height: 24px;
  line-height: 24px;
  text-align: center;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  margin-right: 12px;
  font-weight: bold;
  color: var(--muted);
  font-size: 12px;
}
.rank-num.top-3 {
  background-color: var(--accent);
  color: #0f172a;
}
.rank-info {
  flex: 1;
  overflow: hidden;
}
.rank-title {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: 14px;
  color: var(--text);
  margin-bottom: 4px;
}
.rank-count {
  font-size: 12px;
  color: var(--muted);
}
</style>
