import { defineStore } from 'pinia'
import { login, register } from '@/api'
import type { UserProfile } from '@/types'

interface AuthState {
  token: string | null
  profile: UserProfile | null
}

export const useUserStore = defineStore('user', {
  state: (): AuthState => ({
    token: localStorage.getItem('token'),
    profile: localStorage.getItem('profile') ? JSON.parse(localStorage.getItem('profile') as string) : null,
  }),
  getters: {
    isAuthed: (state) => !!state.token,
  },
  actions: {
    async login(email: string, password: string) {
      const res = await login(email, password)
      this.token = res.token
      this.profile = res.user
      localStorage.setItem('token', this.token || '')
      localStorage.setItem('profile', JSON.stringify(this.profile))
    },
    async register(email: string, password: string, displayName: string) {
      const res = await register(email, password, displayName)
      this.token = res.token
      this.profile = res.user
      localStorage.setItem('token', this.token || '')
      localStorage.setItem('profile', JSON.stringify(this.profile))
    },
    logout() {
      this.token = null
      this.profile = null
      localStorage.removeItem('token')
      localStorage.removeItem('profile')
    },
  },
})
