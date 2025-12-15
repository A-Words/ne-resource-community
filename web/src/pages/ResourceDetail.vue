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
            <el-tag type="success" effect="plain">v{{ resource.version || '1.0' }}</el-tag>
          </div>
          <p class="muted">{{ resource.vendor }} {{ resource.deviceModel }}</p>
          <p>{{ resource.description }}</p>
          <el-descriptions :column="2" border size="small" class="meta">
            <el-descriptions-item label="协议">{{ resource.protocol || '-' }}</el-descriptions-item>
            <el-descriptions-item label="场景">{{ resource.scenario || '-' }}</el-descriptions-item>
            <el-descriptions-item label="标签">{{ resource.tags || '-' }}</el-descriptions-item>
            <el-descriptions-item label="下载">{{ resource.downloadCount }}</el-descriptions-item>
          </el-descriptions>
          <div class="actions-row">
            <el-button type="primary" @click="download">
              {{ resource.externalLink ? '访问链接' : '下载' }}
            </el-button>
            <el-button @click="toggleFav">收藏</el-button>
            <el-button type="danger" link @click="showReport = true">举报</el-button>
            <el-button v-if="isUploader" type="warning" link @click="goUpdate">更新版本</el-button>
          </div>
        </el-card>

        <el-card shadow="never" style="margin-top: 12px" v-if="versions.length > 1">
          <h3>版本历史</h3>
          <el-table :data="versions" size="small">
            <el-table-column prop="version" label="版本" width="80" />
            <el-table-column prop="createdAt" label="时间">
              <template #default="scope">{{ new Date(scope.row.createdAt).toLocaleDateString() }}</template>
            </el-table-column>
            <el-table-column label="操作">
              <template #default="scope">
                <el-button link type="primary" @click="go(scope.row.id)" v-if="scope.row.id !== resource.id">查看</el-button>
                <span v-else>当前</span>
              </template>
            </el-table-column>
          </el-table>
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
        <el-card shadow="never" style="margin-bottom: 12px" v-if="userStore.isAuthed">
          <h4>学习进度</h4>
          <div class="progress-box">
            <el-progress :percentage="learningProgress.progress" />
            <div style="margin-top: 10px; display: flex; justify-content: space-between; align-items: center">
              <el-radio-group v-model="learningProgress.status" size="small" @change="saveProgress">
                <el-radio-button label="started">进行中</el-radio-button>
                <el-radio-button label="completed">已完成</el-radio-button>
              </el-radio-group>
              <el-input-number 
                v-model="learningProgress.progress" 
                :min="0" :max="100" 
                size="small" 
                style="width: 100px"
                @change="saveProgress" 
              />
            </div>
          </div>
        </el-card>

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

    <el-dialog v-model="showReport" title="举报资源" width="30%">
      <el-input v-model="reportReason" type="textarea" placeholder="请输入举报原因（如：内容错误、侵权、病毒等）" />
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showReport = false">取消</el-button>
          <el-button type="primary" @click="submitReport">提交</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch, computed } from 'vue'
import { fetchResource, fetchRecommendations, downloadResource, submitReview as apiSubmitReview, toggleFavorite, reportResource, fetchVersions, fetchProgress, updateProgress } from '@/api'
import type { Resource } from '@/types'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const resource = ref<Resource | null>(null)
const learningProgress = reactive({
  progress: 0,
  status: 'not_started'
})

watch(() => resource.value, async (val) => {
  if (val && userStore.isAuthed) {
    try {
      const res = await fetchProgress(val.id)
      if (res) {
        learningProgress.progress = res.progress
        learningProgress.status = res.status || 'started'
      }
    } catch (e) {
      // ignore error if progress not found or failed
    }
  }
})

async function saveProgress() {
  if (!resource.value) return
  await updateProgress(resource.value.id, learningProgress.progress, learningProgress.status)
}
const recommendations = ref<Resource[]>([])
const versions = ref<Resource[]>([])
const loadingRecommend = ref(false)
const review = reactive({ score: 4, comment: '' })
const showReport = ref(false)
const reportReason = ref('')

const isUploader = computed(() => {
  return resource.value && userStore.profile && resource.value.uploaderId === userStore.profile.id
})

async function load() {
  const id = route.params.id as string
  resource.value = await fetchResource(id)
  loadRecommend(id)
  loadVersions(id)
}

async function loadVersions(id: string) {
  versions.value = await fetchVersions(id)
}

async function loadRecommend(id: string) {
  loadingRecommend.value = true
  recommendations.value = await fetchRecommendations(id)
  loadingRecommend.value = false
}

async function download() {
  if (!resource.value) return
  
  if (resource.value.externalLink) {
    window.open(resource.value.externalLink, '_blank')
    return
  }

  try {
    await downloadResource(resource.value.id)
    ElMessage.success('开始下载')
  } catch (err: any) {
    ElMessage.error(err?.response?.data?.error || '下载失败')
  }
}

async function toggleFav() {
  if (!resource.value) return
  try {
    const res = await toggleFavorite(resource.value.id)
    ElMessage.success(res.status === 'added' ? '已收藏' : '已取消收藏')
  } catch (err: any) {
    ElMessage.error(err?.response?.data?.error || '操作失败')
  }
}

async function submitReview() {
  if (!resource.value) return
  try {
    await apiSubmitReview(resource.value.id, { ...review })
    ElMessage.success('感谢反馈')
  } catch (err: any) {
    ElMessage.error(err?.response?.data?.error || '提交失败')
  }
}

async function submitReport() {
  if (!resource.value) return
  if (!reportReason.value) {
    ElMessage.warning('请输入原因')
    return
  }
  try {
    await reportResource(resource.value.id, reportReason.value)
    ElMessage.success('举报已提交')
    showReport.value = false
    reportReason.value = ''
  } catch (err: any) {
    ElMessage.error(err?.response?.data?.error || '提交失败')
  }
}

function go(id: string) {
  router.push(`/resources/${id}`)
}

function goUpdate() {
  if (!resource.value) return
  router.push({
    path: '/upload',
    query: {
      parentId: resource.value.id,
      title: resource.value.title,
      type: resource.value.type,
      vendor: resource.value.vendor,
      deviceModel: resource.value.deviceModel,
      protocol: resource.value.protocol,
      scenario: resource.value.scenario,
      tags: resource.value.tags,
    }
  })
}

watch(
  () => route.params.id,
  (newId) => {
    if (newId) load()
  }
)

onMounted(load)
</script>

<style scoped>
.title-row {
  display: flex;
  align-items: center;
  gap: 12px;
}
.actions-row {
  margin-top: 16px;
  display: flex;
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
