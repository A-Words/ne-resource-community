<template>
  <div class="page-shell" v-if="resource">
    <el-breadcrumb separator="/">
      <el-breadcrumb-item to="/">资源广场</el-breadcrumb-item>
      <el-breadcrumb-item>{{ resource.title }}</el-breadcrumb-item>
    </el-breadcrumb>

    <el-row :gutter="16" style="margin-top: 12px">
      <el-col :md="16" :xs="24">
        <el-card shadow="never">
          <div class="title-row">
            <h2>{{ resource.title }}</h2>
            <el-tag type="info">{{ resource.type }}</el-tag>
          </div>
          <p class="muted">{{ resource.vendor }} {{ resource.deviceModel }}</p>
          <p>{{ resource.description }}</p>
          <el-descriptions :column="2" border size="small" class="meta">
            <el-descriptions-item label="协议">{{ resource.protocol || '-' }}</el-descriptions-item>
            <el-descriptions-item label="场景">{{ resource.scenario || '-' }}</el-descriptions-item>
            <el-descriptions-item label="标签">{{ resource.tags || '-' }}</el-descriptions-item>
            <el-descriptions-item label="下载">{{ resource.downloadCount }}</el-descriptions-item>
          </el-descriptions>
          <el-button type="primary" @click="download">下载</el-button>
        </el-card>

        <el-card shadow="never" style="margin-top: 12px">
          <h3>评价</h3>
          <div class="review-box">
            <el-rate v-model="review.score" />
            <el-input v-model="review.comment" type="textarea" placeholder="写下你的体验" />
            <el-button type="primary" size="small" @click="submitReview">提交评价</el-button>
          </div>
        </el-card>
      </el-col>

      <el-col :md="8" :xs="24">
        <el-card shadow="never">
          <h4>相似推荐</h4>
          <el-skeleton v-if="loadingRecommend" :rows="3" animated />
          <el-empty v-else-if="recommendations.length === 0" description="暂无" />
          <el-timeline v-else>
            <el-timeline-item v-for="item in recommendations" :key="item.id">
              <div class="rec-item" @click="go(item.id)">
                <div>{{ item.title }}</div>
                <small class="muted">{{ item.vendor }} / {{ item.type }}</small>
              </div>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue'
import { fetchResource, fetchRecommendations, downloadResource, submitReview } from '@/api'
import type { Resource } from '@/types'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const resource = ref<Resource | null>(null)
const recommendations = ref<Resource[]>([])
const loadingRecommend = ref(false)
const review = reactive({ score: 4, comment: '' })

async function load() {
  const id = route.params.id as string
  resource.value = await fetchResource(id)
  loadRecommend(id)
}

async function loadRecommend(id: string) {
  loadingRecommend.value = true
  recommendations.value = await fetchRecommendations(id)
  loadingRecommend.value = false
}

async function download() {
  if (!resource.value) return
  try {
    await downloadResource(resource.value.id)
    ElMessage.success('开始下载')
  } catch (err: any) {
    ElMessage.error(err?.response?.data?.error || '下载失败')
  }
}

async function submitReview() {
  if (!resource.value) return
  try {
    await submitReview(resource.value.id, { ...review })
    ElMessage.success('感谢反馈')
  } catch (err: any) {
    ElMessage.error(err?.response?.data?.error || '提交失败')
  }
}

function go(id: string) {
  router.push(`/resources/${id}`)
}

watch(
  () => route.params.id,
  () => load(),
)

onMounted(load)
</script>

<style scoped>
.title-row {
  display: flex;
  align-items: center;
  gap: 12px;
}
.meta {
  margin: 12px 0;
}
.muted {
  color: var(--muted);
}
.review-box {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.rec-item {
  cursor: pointer;
}
</style>
