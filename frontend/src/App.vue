<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import ConfigPanel from './components/ConfigPanel.vue'
import ScanResults from './components/ScanResults.vue'

const activeTab = ref('config')
const config = ref<any>(null)
const scanResults = ref<any>(null)
const isScanning = ref(false)
const hasScanned = ref(false)
const isInitializing = ref(true)

// 动态导入 Wails 绑定
let GetConfig: any
let StartScan: any
let SelectDirectory: any
let AddScanPath: any

const loadConfig = async () => {
  try {
    if (!GetConfig) {
      console.warn('Wails 绑定未初始化，跳过配置加载')
      return
    }
    config.value = await GetConfig()
  } catch (error) {
    console.error('加载配置失败:', error)
    // 只在绑定已加载的情况下显示错误消息
    if (GetConfig) {
      ElMessage.error('加载配置失败: ' + error)
    }
  }
}

onMounted(async () => {
  // 等待 Wails 运行时准备好
  try {
    // 添加延迟，确保 window.go 对象已经注入
    await new Promise(resolve => setTimeout(resolve, 100))

    const module = await import('../wailsjs/go/main/App')
    GetConfig = module.GetConfig
    StartScan = module.StartScan
    SelectDirectory = module.SelectDirectory
    AddScanPath = module.AddScanPath

    await loadConfig()
  } catch (error: any) {
    console.error('加载 Wails 绑定失败:', error)
    // 只在非开发环境或真正的错误时显示消息
    if (error?.message && !error.message.includes('reading \'main\'')) {
      ElMessage.error('初始化失败: ' + error.message)
    }
  } finally {
    isInitializing.value = false
  }
})

const handleScan = async () => {
  isScanning.value = true
  try {
    scanResults.value = await StartScan()
    hasScanned.value = true
    activeTab.value = 'results'

    const count = scanResults.value?.totalCount || 0
    if (count > 0) {
      ElMessage.success(`扫描完成，找到 ${count} 个项目`)
    } else {
      ElMessage.info('扫描完成，未找到可清理的项目')
    }
  } catch (error) {
    console.error('扫描失败:', error)
    ElMessage.error('扫描失败: ' + error)
  } finally {
    isScanning.value = false
  }
}

const handleAddPath = async () => {
  try {
    const path = await SelectDirectory()

    if (path && path.trim() !== '') {
      await AddScanPath(path)
      ElMessage.success('路径已添加')
      await loadConfig()
      // 清空旧的扫描结果
      scanResults.value = null
    }
  } catch (error) {
    console.error('添加路径失败:', error)
    ElMessage.error('添加路径失败: ' + error)
  }
}

const handleReloadConfig = async () => {
  await loadConfig()
  // 清空旧的扫描结果，因为配置已改变
  const hadResults = scanResults.value !== null
  scanResults.value = null
  hasScanned.value = false

  // 如果当前在结果页，切换回配置页
  if (activeTab.value === 'results') {
    activeTab.value = 'config'
    if (hadResults) {
      ElMessage.info('配置已更新，请重新扫描')
    }
  }
}
</script>

<template>
  <div class="app-container">
    <el-container>
      <el-header class="app-header">
        <h1>🧹 Fast Clean X</h1>
        <p>快速清理项目构建文件</p>
      </el-header>

      <el-main>
        <!-- 初始化加载状态 -->
        <div v-if="isInitializing" class="loading-container">
          <el-icon class="is-loading" :size="50" color="#409eff">
            <Loading />
          </el-icon>
          <p style="margin-top: 20px; color: #606266;">正在初始化...</p>
        </div>

        <!-- 主内容 -->
        <el-tabs v-else v-model="activeTab" class="main-tabs">
          <el-tab-pane label="配置" name="config">
            <ConfigPanel
              :config="config"
              @scan="handleScan"
              @add-path="handleAddPath"
              @reload-config="handleReloadConfig"
              :is-scanning="isScanning"
              :has-scanned="hasScanned"
            />
          </el-tab-pane>

          <el-tab-pane name="results">
            <template #label>
              <span>
                扫描结果
                <el-badge
                  v-if="scanResults?.totalCount"
                  :value="scanResults.totalCount"
                  :max="99"
                  style="margin-left: 5px;"
                />
              </span>
            </template>
            <ScanResults :results="scanResults" />
          </el-tab-pane>
        </el-tabs>
      </el-main>
    </el-container>
  </div>
</template>

<style scoped>
.app-container {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.app-header {
  background: rgba(255, 255, 255, 0.95);
  text-align: center;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  height: auto !important;
  line-height: normal;
}

.app-header h1 {
  margin: 0;
  font-size: 32px;
  color: #333;
  line-height: 1.2;
}

.app-header p {
  margin: 8px 0 0;
  color: #666;
  font-size: 14px;
  line-height: 1.5;
}

.el-main {
  padding: 20px;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 200px);
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.main-tabs {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  min-height: calc(100vh - 200px);
}
</style>
