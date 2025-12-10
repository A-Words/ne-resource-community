<template>
  <div class="page-shell auth">
    <el-card class="auth-card" shadow="never">
      <h2>{{ mode === 'login' ? '登录' : '注册' }}</h2>
      <el-form label-position="top" :model="form">
        <el-form-item label="邮箱">
          <el-input v-model="form.email" type="email" />
        </el-form-item>
        <el-form-item v-if="mode === 'register'" label="昵称">
          <el-input v-model="form.displayName" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="submit">{{ mode === 'login' ? '登录' : '注册' }}</el-button>
          <el-button link type="info" @click="toggle">{{ mode === 'login' ? '去注册' : '去登录' }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const mode = ref<'login' | 'register'>('login')
const loading = ref(false)
const router = useRouter()
const user = useUserStore()

const form = reactive({
  email: '',
  password: '',
  displayName: '',
})

function toggle() {
  mode.value = mode.value === 'login' ? 'register' : 'login'
}

async function submit() {
  loading.value = true
  try {
    if (mode.value === 'login') {
      await user.login(form.email, form.password)
    } else {
      await user.register(form.email, form.password, form.displayName)
    }
    ElMessage.success('操作成功')
    router.push('/')
  } catch (err: any) {
    ElMessage.error(err?.response?.data?.error || '请求失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth {
  display: flex;
  justify-content: center;
  padding-top: 60px;
}
.auth-card {
  width: 420px;
}
</style>
