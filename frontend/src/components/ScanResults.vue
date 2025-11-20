<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { Delete, Folder, Document, ArrowRight, FolderOpened } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const props = defineProps<{
  results: any
}>()

const isCleaning = ref(false)
const selectedItems = ref<any[]>([])
const expandedProjects = ref<Set<string>>(new Set())

// 动态导入 Wails 绑定
let StartClean: any
let OpenFolder: any

onMounted(async () => {
  try {
    const module = await import('../../wailsjs/go/main/App')
    StartClean = module.StartClean
    OpenFolder = module.OpenFolder
  } catch (error) {
    console.error('加载 Wails 绑定失败:', error)
  }
})

// 按项目分组
const groupedItems = computed(() => {
  if (!props.results?.items) return []

  const groups = new Map<string, any>()

  props.results.items.forEach((item: any) => {
    const projectPath = item.projectPath || item.path

    if (!groups.has(projectPath)) {
      groups.set(projectPath, {
        projectPath,
        projectName: item.projectName || projectPath.split('/').pop(),
        items: [],
        totalSize: 0,
        totalFiles: 0,
        allSelected: false
      })
    }

    const group = groups.get(projectPath)
    group.items.push(item)
    group.totalSize += item.size
    group.totalFiles += item.fileCount
  })

  // 更新每个组的选中状态
  groups.forEach(group => {
    group.allSelected = group.items.every((item: any) => item.selected)
  })

  return Array.from(groups.values())
})

const totalSize = computed(() => {
  if (!props.results?.items) return '0 B'
  const total = props.results.items
    .filter((item: any) => item.selected)
    .reduce((sum: number, item: any) => sum + item.size, 0)
  return formatSize(total)
})

const selectedCount = computed(() => {
  return props.results?.items?.filter((item: any) => item.selected).length || 0
})

const formatSize = (bytes: number) => {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let size = bytes
  let unitIndex = 0
  
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  
  return `${size.toFixed(2)} ${units[unitIndex]}`
}

const handleClean = async () => {
  if (selectedCount.value === 0) {
    ElMessage.warning('请至少选择一个项目进行清理')
    return
  }

  if (!StartClean) {
    ElMessage.error('功能未就绪')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除 ${selectedCount.value} 个项目吗？总大小: ${totalSize.value}`,
      '确认清理',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    isCleaning.value = true
    await StartClean(props.results.items)
    ElMessage.success('清理完成！')

    // 重新加载页面
    setTimeout(() => {
      window.location.reload()
    }, 1500)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('清理失败:', error)
      ElMessage.error('清理失败，请重试')
    }
  } finally {
    isCleaning.value = false
  }
}

const toggleSelection = (item: any) => {
  item.selected = !item.selected
}

const selectAll = () => {
  props.results?.items?.forEach((item: any) => {
    item.selected = true
  })
}

const deselectAll = () => {
  props.results?.items?.forEach((item: any) => {
    item.selected = false
  })
}

const getTypeTag = (type: string) => {
  const tags: Record<string, { type: string; effect: string }> = {
    'Maven': { type: 'danger', effect: 'dark' },
    'Gradle': { type: 'success', effect: 'dark' },
    'Node.js': { type: 'primary', effect: 'dark' },
    'Python': { type: 'warning', effect: 'dark' },
    'Rust': { type: 'info', effect: 'dark' },
    'Go': { type: 'primary', effect: 'dark' },
  }
  return tags[type] || { type: 'info', effect: 'dark' }
}

const handleOpenFolder = async (path: string) => {
  try {
    if (!OpenFolder) {
      ElMessage.error('功能未就绪')
      return
    }
    await OpenFolder(path)
  } catch (error) {
    console.error('打开文件夹失败:', error)
    ElMessage.error('打开文件夹失败: ' + error)
  }
}

// 切换项目展开/折叠
const toggleProject = (projectPath: string) => {
  if (expandedProjects.value.has(projectPath)) {
    expandedProjects.value.delete(projectPath)
  } else {
    expandedProjects.value.add(projectPath)
  }
}

// 切换项目所有项的选中状态
const toggleProjectSelection = (group: any) => {
  const newState = !group.allSelected
  group.items.forEach((item: any) => {
    item.selected = newState
  })
}

// 展开所有项目
const expandAll = () => {
  groupedItems.value.forEach(group => {
    expandedProjects.value.add(group.projectPath)
  })
}

// 折叠所有项目
const collapseAll = () => {
  expandedProjects.value.clear()
}
</script>

<template>
  <div class="scan-results">
    <div v-if="!results || !results.items || results.items.length === 0" class="empty-state">
      <el-empty>
        <template #image>
          <el-icon :size="100" color="#909399">
            <Folder />
          </el-icon>
        </template>
        <template #description>
          <p style="font-size: 16px; color: #606266; margin-bottom: 10px;">
            {{ !results ? '还没有扫描结果' : '未找到可清理的项目' }}
          </p>
          <p style="font-size: 14px; color: #909399;">
            {{ !results ? '请在"配置"页面添加扫描路径并开始扫描' : '当前扫描路径下没有匹配的构建目录' }}
          </p>
        </template>
      </el-empty>
    </div>
    
    <div v-else>
      <!-- 统计信息 -->
      <el-row :gutter="20" class="stats">
        <el-col :span="4">
          <el-statistic title="项目数" :value="groupedItems.length">
            <template #prefix>
              <el-icon><Folder /></el-icon>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="4">
          <el-statistic title="构建目录" :value="results.totalCount">
            <template #prefix>
              <el-icon><Document /></el-icon>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="4">
          <el-statistic title="已选择" :value="selectedCount">
            <template #prefix>
              <el-icon><Delete /></el-icon>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="4">
          <el-statistic title="总大小" :value="totalSize">
            <template #prefix>
              <el-icon><Delete /></el-icon>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="8">
          <div class="action-buttons">
            <el-button size="small" @click="selectAll" type="primary">全选</el-button>
            <el-button size="small" @click="deselectAll">取消全选</el-button>
            <el-button size="small" @click="expandAll">全部展开</el-button>
            <el-button size="small" @click="collapseAll">全部折叠</el-button>
          </div>
        </el-col>
      </el-row>
      
      <!-- 结果列表 - 按项目分组 -->
      <el-card shadow="hover" class="results-card">
        <el-scrollbar height="500px">
          <div class="project-groups">
            <div
              v-for="group in groupedItems"
              :key="group.projectPath"
              class="project-group"
            >
              <!-- 项目头部 -->
              <div class="project-header" @click="toggleProject(group.projectPath)">
                <div class="project-header-left">
                  <el-icon class="expand-icon" :class="{ expanded: expandedProjects.has(group.projectPath) }">
                    <ArrowRight />
                  </el-icon>
                  <el-checkbox
                    :model-value="group.allSelected"
                    @change="toggleProjectSelection(group)"
                    @click.stop
                  />
                  <el-icon class="project-icon">
                    <Folder />
                  </el-icon>
                  <span class="project-name">{{ group.projectName }}</span>
                  <el-tag size="small" type="info">{{ group.items.length }} 项</el-tag>
                </div>
                <div class="project-header-right">
                  <span class="project-size">{{ formatSize(group.totalSize) }}</span>
                  <span class="project-files">{{ group.totalFiles }} 文件</span>
                  <el-button
                    size="small"
                    text
                    @click.stop="handleOpenFolder(group.projectPath)"
                    title="打开项目文件夹"
                  >
                    <el-icon><FolderOpened /></el-icon>
                  </el-button>
                </div>
              </div>

              <!-- 项目详情（可折叠） -->
              <el-collapse-transition>
                <div v-show="expandedProjects.has(group.projectPath)" class="project-items">
                  <div
                    v-for="item in group.items"
                    :key="item.path"
                    class="project-item"
                    :class="{ selected: item.selected }"
                  >
                    <div class="item-left">
                      <el-checkbox
                        :model-value="item.selected"
                        @change="toggleSelection(item)"
                      />
                      <el-tag
                        :type="getTypeTag(item.type).type"
                        :effect="getTypeTag(item.type).effect"
                        size="small"
                      >
                        {{ item.type }}
                      </el-tag>
                      <span
                        class="item-path path-link"
                        @click="handleOpenFolder(item.path)"
                        :title="'点击打开: ' + item.path"
                      >
                        {{ item.path }}
                      </span>
                    </div>
                    <div class="item-right">
                      <span class="item-size">{{ item.sizeReadable }}</span>
                      <span class="item-files">{{ item.fileCount }} 文件</span>
                    </div>
                  </div>
                </div>
              </el-collapse-transition>
            </div>
          </div>
        </el-scrollbar>
      </el-card>
      
      <!-- 清理按钮 -->
      <div class="clean-action">
        <el-button 
          type="danger" 
          size="large"
          :icon="Delete"
          :loading="isCleaning"
          :disabled="selectedCount === 0"
          @click="handleClean"
        >
          {{ isCleaning ? '清理中...' : `清理选中项 (${totalSize})` }}
        </el-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.scan-results {
  padding: 20px;
}

.empty-state {
  padding: 100px 20px;
  text-align: center;
}

.stats {
  margin-bottom: 20px;
}

.action-buttons {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  align-items: center;
  padding-top: 20px;
  flex-wrap: wrap;
}

.results-card {
  margin-bottom: 20px;
}

.selected-row {
  background-color: #f0f9ff !important;
}

.size-text {
  font-weight: bold;
  color: #f56c6c;
}

.path-link {
  color: #409eff;
  cursor: pointer;
  transition: all 0.2s;
  text-decoration: none;
  display: inline-block;
}

.path-link:hover {
  color: #66b1ff;
  text-decoration: underline;
}

.path-link:active {
  color: #3a8ee6;
}

.clean-action {
  text-align: center;
  margin-top: 20px;
}

.clean-action .el-button {
  min-width: 250px;
}

:deep(.el-statistic__head) {
  font-size: 14px;
  color: #909399;
}

:deep(.el-statistic__number) {
  font-size: 24px;
  font-weight: bold;
}

/* 项目分组样式 */
.project-groups {
  padding: 10px;
}

.project-group {
  margin-bottom: 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s;
}

.project-group:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.project-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eef5 100%);
  cursor: pointer;
  user-select: none;
  transition: all 0.2s;
}

.project-header:hover {
  background: linear-gradient(135deg, #e8eef5 0%, #dce4f0 100%);
}

.project-header-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.expand-icon {
  transition: transform 0.3s;
  color: #606266;
}

.expand-icon.expanded {
  transform: rotate(90deg);
}

.project-icon {
  color: #409eff;
  font-size: 18px;
}

.project-name {
  font-weight: 600;
  font-size: 15px;
  color: #303133;
}

.project-header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.project-size {
  font-weight: bold;
  color: #f56c6c;
  font-size: 14px;
}

.project-files {
  color: #909399;
  font-size: 13px;
}

.project-items {
  background: #fff;
}

.project-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px 10px 48px;
  border-top: 1px solid #f0f0f0;
  transition: all 0.2s;
}

.project-item:hover {
  background: #f5f7fa;
}

.project-item.selected {
  background: #f0f9ff;
}

.item-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
}

.item-path {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 13px;
}

.item-right {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-shrink: 0;
}

.item-size {
  font-weight: 600;
  color: #f56c6c;
  font-size: 13px;
  min-width: 80px;
  text-align: right;
}

.item-files {
  color: #909399;
  font-size: 12px;
  min-width: 70px;
  text-align: right;
}
</style>

