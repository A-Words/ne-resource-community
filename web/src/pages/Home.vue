<template>
  <div class="page-shell">
    <div class="hero">
      <div>
        <div class="gradient-pill">网络主题资源共享站</div>
        <div class="hero-title">聚焦网络工程的高质量资源库</div>
        <div class="hero-sub">工具、配置模板、文档、课程一站式搜索与下载</div>
      </div>
      <el-input v-model="filters.search" placeholder="关键词 / 协议 / 厂商" clearable @keyup.enter="load">
        <template #append>
          <el-button type="primary" @click="load">搜索</el-button>
        </template>
      </el-input>
    </div>

    <el-card shadow="never" class="panel">
      <el-tabs v-model="filters.type" @tab-change="load" class="type-tabs">
        <el-tab-pane label="全部资源" name="" />
        <el-tab-pane label="网络工具" name="网络工具" />
        <el-tab-pane label="配置模板" name="配置模板" />
        <el-tab-pane label="文档资料" name="文档资料" />
        <el-tab-pane label="学习资源" name="学习资源" />
      </el-tabs>

      <el-form inline :model="filters" label-width="80px" class="filters">
        <template v-if="filters.type !== '网络工具'">
          <el-form-item label="厂商">
            <el-input v-model="filters.vendor" placeholder="如 Cisco / Huawei" clearable />
          </el-form-item>
          <el-form-item label="设备型号">
            <el-input v-model="filters.device" placeholder="如 AR / ASR" clearable />
          </el-form-item>
          <el-form-item label="协议">
            <el-input v-model="filters.protocol" placeholder="BGP / OSPF / VXLAN" clearable />
          </el-form-item>
          <el-form-item label="场景">
            <el-input v-model="filters.scenario" placeholder="数据中心 / 骨干" clearable />
          </el-form-item>
        </template>
        <el-form-item>
          <el-button type="primary" @click="load">筛选</el-button>
          <el-button @click="reset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-row :gutter="16">
        <el-col v-for="item in resources" :key="item.id" :xs="24" :sm="12" :md="8">
          <el-card class="res-card" shadow="hover">
            <div class="card-header">
              <div class="title">{{ item.title }}</div>
              <el-tag size="small" type="info">{{ item.type }}</el-tag>
            </div>
            <div class="meta">{{ item.vendor }} {{ item.deviceModel }}</div>
            <div class="desc">{{ item.description || '暂无描述' }}</div>
            <div class="footer">
              <div class="tags">{{ item.tags }}</div>
              <div class="stats">
                <el-rate :model-value="item.ratingAverage" disabled allow-half />
                <span>{{ item.downloadCount }} 下载</span>
              </div>
            </div>
            <div class="actions">
              <el-button link type="primary" @click="goDetail(item.id)">查看</el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { fetchResources } from '@/api'
import type { Resource } from '@/types'
import { useRouter } from 'vue-router'

const router = useRouter()
const resources = ref<Resource[]>([])

const filters = reactive({
  search: '',
  type: '',
  vendor: '',
  device: '',
  protocol: '',
  scenario: '',
})

async function load() {
  resources.value = await fetchResources({ ...filters })
}

function reset() {
  filters.search = ''
  filters.type = ''
  filters.vendor = ''
  filters.device = ''
  filters.protocol = ''
  filters.scenario = ''
  load()
}

function goDetail(id: string) {
  router.push(`/resources/${id}`)
}

onMounted(load)
</script>

<style scoped>
.type-tabs {
  margin-bottom: 16px;
}

.hero {
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 18px;
  align-items: center;
  margin-bottom: 16px;
}

.panel {
  margin-top: 8px;
}

.filters {
  margin-bottom: 12px;
}

.res-card {
  margin-bottom: 14px;
  min-height: 220px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-weight: 700;
  color: #e0f2fe;
}

.meta {
  color: var(--muted);
  margin-top: 4px;
}

.desc {
  margin: 8px 0;
  color: #cbd5e1;
  min-height: 40px;
}

.footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: var(--muted);
  font-size: 13px;
}

.tags {
  color: #a5b4fc;
}

.actions {
  margin-top: 8px;
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 900px) {
  .hero {
    grid-template-columns: 1fr;
  }
}
</style>
