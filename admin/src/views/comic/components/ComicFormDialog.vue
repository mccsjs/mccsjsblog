<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import type { CreateComicRequest } from '@/types/comic'

const props = defineProps<{
  visible: boolean
  mode: 'create' | 'edit'
  data: CreateComicRequest
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  submit: []
}>()

const formRef = ref()
const local = reactive({ ...props.data })

watch(() => props.visible, (v) => {
  if (v) Object.assign(local, props.data)
})

const rules = {
  name: [{ required: true, message: '请输入番剧名称', trigger: 'blur' }],
  url: [{ required: true, message: '请输入跳转链接', trigger: 'blur' }],
}

const onConfirm = async () => {
  await formRef.value?.validate()
  Object.assign(props.data, local)
  emit('submit')
}
</script>

<template>
  <el-dialog
    :model-value="visible"
    :title="mode === 'create' ? '新增番剧' : '编辑番剧'"
    width="500px"
    @update:model-value="emit('update:visible', $event)"
  >
    <el-form ref="formRef" :model="local" :rules="rules" label-width="100px">
      <el-form-item label="番剧名称" prop="name">
        <el-input v-model="local.name" placeholder="如：某科学的超电磁炮T" />
      </el-form-item>
      <el-form-item label="跳转链接" prop="url">
        <el-input v-model="local.url" placeholder="https://..." />
      </el-form-item>
      <el-form-item label="封面图片">
        <el-input v-model="local.cover" placeholder="封面图片URL" />
      </el-form-item>
      <el-form-item label="排序">
        <el-input-number v-model="local.sort" :min="0" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="emit('update:visible', false)">取消</el-button>
      <el-button type="primary" @click="onConfirm">确定</el-button>
    </template>
  </el-dialog>
</template>
