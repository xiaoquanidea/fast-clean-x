<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { Delete, FolderAdd, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const props = defineProps<{
  config: any
  isScanning: boolean
  hasScanned?: boolean
}>()

const emit = defineEmits(['scan', 'add-path', 'reload-config'])

// 动态导入 Wails 绑定
let RemoveScanPath: any
let UpdateScanRule: any
let OpenFolder: any

onMounted(async () => {
  try {
    const module = await import('../../wailsjs/go/main/App')
    RemoveScanPath = module.RemoveScanPath
    UpdateScanRule = module.UpdateScanRule
    OpenFolder = module.OpenFolder
  } catch (error) {
    console.error('加载 Wails 绑定失败:', error)
  }
})

const handleRemovePath = async (path: string) => {
  try {
    if (!RemoveScanPath) {
      ElMessage.error('功能未就绪')
      return
    }
    await RemoveScanPath(path)
    ElMessage.success('路径已删除，请重新扫描')
    emit('reload-config')
  } catch (error) {
    console.error('删除路径失败:', error)
    ElMessage.error('删除路径失败')
  }
}

const handleRuleChange = async (ruleName: string, enabled: boolean) => {
  console.log('切换规则:', ruleName, enabled)
  try {
    if (!UpdateScanRule) {
      ElMessage.error('功能未就绪')
      return
    }
    await UpdateScanRule(ruleName, enabled)
    // 使用更简洁的提示
    ElMessage({
      message: `规则 ${ruleName} 已${enabled ? '启用' : '禁用'}`,
      type: 'success',
      duration: 2000
    })
    emit('reload-config')
  } catch (error) {
    console.error('更新规则失败:', error)
    ElMessage.error('更新规则失败')
  }
}

const handleOpenFolder = async (path: string) => {
  try {
    if (!OpenFolder) {
      ElMessage.error('功能未就绪')
      return
    }
    await OpenFolder(path)
    console.log('打开文件夹:', path)
  } catch (error) {
    console.error('打开文件夹失败:', error)
    ElMessage.error('打开文件夹失败: ' + error)
  }
}

const totalPaths = computed(() => props.config?.scanPaths?.length || 0)
const enabledRules = computed(() =>
  props.config?.scanRules?.filter((r: any) => r.enabled).length || 0
)
</script>

<template>
  <div class="config-panel">
    <el-row :gutter="20">
      <!-- 扫描路径 -->
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>📁 扫描路径</span>
              <el-button 
                type="primary" 
                :icon="FolderAdd" 
                size="small"
                @click="emit('add-path')"
              >
                添加路径
              </el-button>
            </div>
          </template>
          
          <div v-if="!config || totalPaths === 0" class="empty-state">
            <p>还没有添加扫描路径</p>
            <p class="hint">点击上方按钮添加要扫描的目录</p>
          </div>
          
          <el-scrollbar v-else height="400px">
            <div class="path-list">
              <div
                v-for="(path, index) in config.scanPaths"
                :key="index"
                class="path-item"
              >
                <span
                  class="path-text path-link"
                  @click="handleOpenFolder(path)"
                  :title="'点击打开: ' + path"
                >
                  {{ path }}
                </span>
                <el-button
                  type="danger"
                  :icon="Delete"
                  size="small"
                  text
                  @click="handleRemovePath(path)"
                />
              </div>
            </div>
          </el-scrollbar>
        </el-card>
      </el-col>
      
      <!-- 扫描规则 -->
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>⚙️ 扫描规则</span>
              <el-tag size="small">{{ enabledRules }} / {{ config?.scanRules?.length || 0 }} 已启用</el-tag>
            </div>
          </template>
          
          <el-scrollbar height="400px">
            <div class="rules-list">
              <div
                v-for="rule in config?.scanRules"
                :key="rule.name"
                class="rule-item"
              >
                <div class="rule-info">
                  <el-switch
                    v-model="rule.enabled"
                    @change="handleRuleChange(rule.name, rule.enabled)"
                  />
                  <div class="rule-details">
                    <div class="rule-name">{{ rule.name }}</div>
                    <div class="rule-desc">{{ rule.description }}</div>
                    <div class="rule-targets">
                      <el-tag
                        v-for="dir in rule.targetDirs"
                        :key="dir"
                        size="small"
                        type="info"
                      >
                        {{ dir }}
                      </el-tag>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </el-scrollbar>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 开始扫描按钮 -->
    <div class="scan-action">
      <el-alert
        v-if="hasScanned"
        title="配置已更新"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 20px; max-width: 600px; margin-left: auto; margin-right: auto;"
      >
        <template #default>
          扫描路径或规则已修改，建议重新扫描以获取最新结果
        </template>
      </el-alert>

      <el-button
        type="primary"
        size="large"
        :icon="Search"
        :loading="isScanning"
        :disabled="totalPaths === 0 || enabledRules === 0"
        @click="emit('scan')"
      >
        {{ isScanning ? '扫描中...' : (hasScanned ? '重新扫描' : '开始扫描') }}
      </el-button>
      <p v-if="totalPaths === 0" class="hint">请先添加扫描路径</p>
      <p v-else-if="enabledRules === 0" class="hint">请至少启用一个扫描规则</p>
      <p v-else-if="!hasScanned" class="hint">配置完成后，点击按钮开始扫描</p>
    </div>
  </div>
</template>

<style scoped>
.config-panel {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #909399;
}

.empty-state p {
  margin: 10px 0;
}

.hint {
  font-size: 12px;
  color: #909399;
  margin-top: 10px;
}

.path-list {
  padding: 10px;
}

.path-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  margin-bottom: 8px;
  background: #f5f7fa;
  border-radius: 4px;
  transition: all 0.3s;
}

.path-item:hover {
  background: #e8eaf0;
}

.path-text {
  flex: 1;
  font-size: 13px;
  color: #606266;
  word-break: break-all;
}

.path-text.path-link {
  color: #409eff;
  cursor: pointer;
  transition: color 0.2s;
}

.path-text.path-link:hover {
  color: #66b1ff;
  text-decoration: underline;
}

.path-text.path-link:active {
  color: #3a8ee6;
}

.rules-list {
  padding: 10px;
}

.rule-item {
  padding: 15px;
  margin-bottom: 12px;
  background: #f5f7fa;
  border-radius: 4px;
}

.rule-info {
  display: flex;
  gap: 15px;
  align-items: flex-start;
}

.rule-details {
  flex: 1;
}

.rule-name {
  font-weight: bold;
  color: #303133;
  margin-bottom: 5px;
}

.rule-desc {
  font-size: 12px;
  color: #606266;
  margin-bottom: 8px;
}

.rule-targets {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
}

.scan-action {
  margin-top: 30px;
  text-align: center;
}

.scan-action .el-button {
  min-width: 200px;
}
</style>

