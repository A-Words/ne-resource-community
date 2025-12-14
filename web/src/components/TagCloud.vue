<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { fetchPopularTags } from '@/api'
import { useRouter } from 'vue-router'

const tags = ref<{ tag: string; count: number }[]>([])
const router = useRouter()

onMounted(async () => {
  try {
    tags.value = await fetchPopularTags()
  } catch (e) {
    console.error('Failed to fetch tags', e)
  }
})

const handleTagClick = (tag: string) => {
  router.push({ name: 'Home', query: { q: tag } })
}

const getFontSize = (count: number) => {
  if (tags.value.length === 0) return '1rem'
  const max = Math.max(...tags.value.map(t => t.count))
  const min = Math.min(...tags.value.map(t => t.count))
  if (max === min) return '1rem'
  // Scale between 0.8rem and 2rem
  const size = 0.8 + ((count - min) / (max - min)) * 1.2
  return `${size}rem`
}
</script>

<template>
  <div class="tag-cloud card" v-if="tags.length > 0">
    <h3>热门标签</h3>
    <div class="tags">
      <span
        v-for="t in tags"
        :key="t.tag"
        class="tag"
        :style="{ fontSize: getFontSize(t.count) }"
        @click="handleTagClick(t.tag)"
      >
        {{ t.tag }}
      </span>
    </div>
  </div>
</template>

<style scoped>
.tag-cloud {
  padding: 1.5rem;
  background: var(--panel);
  border: 1px solid var(--border);
  border-radius: 8px;
  margin-bottom: 2rem;
}
h3 {
  margin: 0 0 1rem 0;
  font-size: 1.2rem;
  color: var(--text);
  border-bottom: 1px solid var(--border);
  padding-bottom: 0.5rem;
}
.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.8rem;
  align-items: center;
}
.tag {
  cursor: pointer;
  color: var(--accent);
  transition: all 0.3s ease;
  padding: 0.3rem 0.8rem;
  background: rgba(34, 211, 238, 0.1);
  border: 1px solid rgba(34, 211, 238, 0.2);
  border-radius: 20px;
  line-height: 1.2;
}
.tag:hover {
  color: #0f172a;
  background: var(--accent);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(34, 211, 238, 0.3);
}
</style>
