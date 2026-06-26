<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getComics, createComic, updateComic, deleteComic } from '@/api/comic'
import type { Comic, CreateComicRequest, UpdateComicRequest } from '@/types/comic'
import ComicFormDialog from './components/ComicFormDialog.vue'

const list = ref<Comic[]>([])
const total = ref(0)
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)

const formVisible = ref(false)
const formMode = ref<'create' | 'edit'>('create')
const editId = ref<number | null>(null)
const formData = ref<CreateComicRequest>({ name: '', url: '', cover: '', sort: 0 })

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getComics({ page: page.value, page_size: pageSize.value })
    list.value = res.list || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  formMode.value = 'create'
  editId.value = null
  formData.value = { name: '', url: '', cover: '', sort: 0 }
  formVisible.value = true
}

const handleEdit = (row: Comic) => {
  formMode.value = 'edit'
  editId.value = row.id
  formData.value = { name: row.name, url: row.url, cover: row.cover, sort: row.sort }
  formVisible.value = true
}

const handleDelete = async (row: Comic) => {
  try {
    await ElMessageBox.confirm(`确定删除「${row.name}」吗？`, '确认删除', { type: 'warning' })
    await deleteComic(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch { /* cancelled */ }
}

const handleSubmit = async () => {
  try {
    if (formMode.value === 'create') {
      await createComic(formData.value)
      ElMessage.success('创建成功')
    } else {
      await updateComic(editId.value!, formData.value as UpdateComicRequest)
      ElMessage.success('更新成功')
    }
    formVisible.value = false
    fetchList()
  } catch { /* error handled by request interceptor */ }
}

onMounted(fetchList)
</script>

<template>
  <div class="comic-container">
    <div class="page-header">
      <h2>番剧管理</h2>
      <el-button type="primary" @click="handleCreate">新增番剧</el-button>
    </div>

    <el-table :data="list" v-loading="loading" stripe>
      <el-table-column label="封面" width="100">
        <template #default="{ row }">
          <el-image
            v-if="row.cover"
            :src="row.cover"
            fit="cover"
            style="width:60px;height:80px;border-radius:6px"
          />
          <span v-else style="color:#999">无</span>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="150" />
      <el-table-column prop="url" label="跳转链接" min-width="200" show-overflow-tooltip />
      <el-table-column prop="sort" label="排序" width="80" />
      <el-table-column label="操作" width="160">
        <template #default="{ row }">
          <el-button type="primary" size="small" link @click="handleEdit(row)">编辑</el-button>
          <el-button type="danger" size="small" link @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-if="total > pageSize"
      v-model:current-page="page"
      :page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="fetchList"
      style="margin-top:16px;justify-content:flex-end"
    />

    <ComicFormDialog
      v-model:visible="formVisible"
      :mode="formMode"
      :data="formData"
      @submit="handleSubmit"
    />
  </div>
</template>

<style scoped>
.comic-container {
  padding: 20px;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.page-header h2 {
  margin: 0;
}
</style>
