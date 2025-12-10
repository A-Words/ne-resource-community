<template>
  <div class="page-shell">
    <el-card shadow="never">
      <h2>{{ isUpdate ? '更新资源版本' : '上传新资源' }}</h2>
      <p class="muted">文件将存储在服务器本地磁盘，审核后对外展示。</p>
      <el-form :model="form" label-width="120px" class="upload-form">
        <el-form-item label="标题">
          <el-input v-model="form.title" placeholder="如 GNS3 实验拓扑" />
        </el-form-item>
        <el-form-item label="版本号" v-if="isUpdate">
          <el-input v-model="form.version" placeholder="如 1.1, 2.0" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.type" placeholder="选择类型">
            <el-option v-for="t in types" :key="t" :label="t" :value="t" />
          </el-select>
        </el-form-item>
        <template v-if="form.type !== '网络工具'">
          <el-form-item label="厂商">
            <el-input v-model="form.vendor" placeholder="Cisco / Huawei" />
          </el-form-item>
          <el-form-item label="设备型号">
            <el-input v-model="form.deviceModel" placeholder="ASR9K / CE6800" />
          </el-form-item>
          <el-form-item label="协议">
            <el-input v-model="form.protocol" placeholder="BGP / OSPF / VXLAN" />
          </el-form-item>
          <el-form-item label="场景">
            <el-input v-model="form.scenario" placeholder="数据中心 / 骨干 / 实验" />
          </el-form-item>
        </template>
        <el-form-item label="标签">
          <el-input v-model="form.tags" placeholder="逗号分隔：BGP,Route Reflector" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" rows="4" />
        </el-form-item>
        
        <el-form-item label="资源来源">
          <el-radio-group v-model="sourceType">
            <el-radio label="file">上传文件</el-radio>
            <el-radio label="link">外部链接</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="文件" v-if="sourceType === 'file'">
          <el-upload drag :auto-upload="false" :on-change="onFileChange" accept="*/*">
            <i class="el-icon-upload" />
            <div class="el-upload__text">拖拽或点击上传文件</div>
          </el-upload>
          <div v-if="form.file" class="file-name">{{ form.file.name }}</div>
        </el-form-item>

        <el-form-item label="链接地址" v-if="sourceType === 'link'">
          <el-input v-model="form.externalLink" placeholder="https://..." />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="submit">提交</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted, computed } from 'vue'
import { uploadResource } from '@/api'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()
const types = ['网络工具', '配置模板', '文档资料', '学习资源']
const submitting = ref(false)
const sourceType = ref('file')

const form = reactive({
  title: '',
  type: '',
  vendor: '',
  deviceModel: '',
  protocol: '',
  scenario: '',
  tags: '',
  description: '',
  version: '1.0',
  parentId: '',
  file: null as File | null,
  externalLink: '',
})

const isUpdate = computed(() => !!form.parentId)

onMounted(() => {
  if (route.query.parentId) {
    form.parentId = route.query.parentId as string
    form.title = route.query.title as string
    form.type = route.query.type as string
    form.vendor = route.query.vendor as string
    form.deviceModel = route.query.deviceModel as string
    form.protocol = route.query.protocol as string
    form.scenario = route.query.scenario as string
    form.tags = route.query.tags as string
    form.version = '1.1' // Default increment suggestion
  }
})

function onFileChange(file: any) {
  form.file = file.raw
}

async function submit() {
  if (sourceType.value === 'file' && !form.file) {
    ElMessage.error('请上传文件')
    return
  }
  if (sourceType.value === 'link' && !form.externalLink) {
    ElMessage.error('请填写外部链接')
    return
  }

  submitting.value = true
  try {
    await uploadResource({
      title: form.title,
      type: form.type,
      vendor: form.vendor,
      deviceModel: form.deviceModel,
      protocol: form.protocol,
      scenario: form.scenario,
      tags: form.tags,
      description: form.description,
      version: form.version,
      parentId: form.parentId,
      file: sourceType.value === 'file' ? form.file! : undefined,
      externalLink: sourceType.value === 'link' ? form.externalLink : undefined,
    })
    ElMessage.success('提交成功，等待审核/发布')
    router.push('/')
  } catch (err: any) {
    ElMessage.error(err?.response?.data?.error || '提交失败')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.muted {
  color: var(--muted);
  margin-top: -6px;
}

.upload-form {
  margin-top: 12px;
}

.file-name {
  color: #a5b4fc;
}
</style>
