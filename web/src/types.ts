export interface UserProfile {
  id: string
  email: string
  displayName: string
  role: string
  points: number
  level: number
}

export interface Resource {
  id: string
  title: string
  description: string
  type: string
  vendor: string
  deviceModel: string
  protocol: string
  scenario: string
  tags: string
  filePath: string
  fileName: string
  contentType: string
  downloadCount: number
  ratingAverage: number
  ratingCount: number
  uploaderId: string
  uploader?: UserProfile
  createdAt: string
  updatedAt: string
  parentId?: string
  version?: string
  externalLink?: string
}

export interface ResourcePayload {
  title: string
  description?: string
  type: string
  vendor?: string
  deviceModel?: string
  protocol?: string
  scenario?: string
  tags?: string
  file?: File
  externalLink?: string
  parentId?: string
  version?: string
}

export interface ReviewPayload {
  score: number
  comment?: string
}
