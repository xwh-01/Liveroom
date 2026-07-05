<template>
  <div class="room-lobby">
    <header class="lobby-header">
      <h1>LiveRoom Battle</h1>
      <el-button type="primary" @click="showCreateDialog = true">创建房间</el-button>
    </header>

    <div class="lobby-body">
      <div v-if="loading" class="lobby-loading">加载中...</div>
      <div v-else-if="rooms.length === 0" class="lobby-empty">暂无直播间，快来创建一个吧！</div>
      <div v-else class="room-grid">
        <div v-for="room in rooms" :key="room.room_id" class="room-card">
          <div class="room-card-header">
            <span class="room-title">{{ room.title }}</span>
            <el-tag :type="room.status === 'live' ? 'success' : 'info'" size="small">
              {{ room.status === 'live' ? '直播中' : '已关闭' }}
            </el-tag>
          </div>
          <div class="room-card-body">
            <div class="room-info">
              <span class="info-label">房主</span>
              <span>{{ room.owner_name }}</span>
            </div>
            <div class="room-stats">
              <div class="stat-item">
                <span class="stat-value">{{ room.online_count }}</span>
                <span class="stat-label">在线</span>
              </div>
              <div class="stat-item">
                <span class="stat-value">{{ room.chat_count }}</span>
                <span class="stat-label">弹幕</span>
              </div>
              <div class="stat-item">
                <span class="stat-value">{{ room.gift_count }}</span>
                <span class="stat-label">礼物</span>
              </div>
            </div>
          </div>
          <div class="room-card-footer">
            <span class="room-id">ID: {{ room.room_id }}</span>
            <el-button type="primary" size="small" @click="enterRoom(room.room_id)">进入</el-button>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="showCreateDialog" title="创建直播间" width="400px">
      <el-form :model="createForm" label-width="80px">
        <el-form-item label="房间标题" required>
          <el-input v-model="createForm.title" placeholder="输入房间标题" />
        </el-form-item>
        <el-form-item label="房主昵称">
          <el-input v-model="createForm.owner_name" placeholder="anonymous" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="doCreateRoom" :loading="creating">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()

const rooms = ref([])
const loading = ref(true)
const showCreateDialog = ref(false)
const creating = ref(false)
const createForm = ref({ title: '', owner_name: '' })

async function fetchRooms() {
  try {
    const res = await fetch('/api/rooms?limit=50')
    const data = await res.json()
    if (data.code === 0) {
      rooms.value = data.data || []
    }
  } catch (e) {
    ElMessage.error('加载房间列表失败')
  } finally {
    loading.value = false
  }
}

async function doCreateRoom() {
  if (!createForm.value.title.trim()) {
    ElMessage.warning('请输入房间标题')
    return
  }
  creating.value = true
  try {
    const res = await fetch('/api/rooms', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        title: createForm.value.title,
        owner_name: createForm.value.owner_name || 'anonymous'
      })
    })
    const data = await res.json()
    if (data.code === 0) {
      ElMessage.success('房间创建成功')
      showCreateDialog.value = false
      createForm.value = { title: '', owner_name: '' }
      router.push(`/room/${data.data.room_id}`)
    } else {
      ElMessage.error(data.msg || '创建失败')
    }
  } catch (e) {
    ElMessage.error('创建房间失败')
  } finally {
    creating.value = false
  }
}

function enterRoom(roomId) {
  router.push(`/room/${roomId}`)
}

onMounted(fetchRooms)
</script>

<style scoped>
.room-lobby {
  min-height: 100vh;
  background: #f0f2f5;
}

.lobby-header {
  height: 56px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.15);
}

.lobby-header h1 {
  font-size: 20px;
  font-weight: 700;
  letter-spacing: 2px;
}

.lobby-body {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.lobby-loading, .lobby-empty {
  text-align: center;
  color: #909399;
  margin-top: 80px;
  font-size: 14px;
}

.room-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.room-card {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.06);
  display: flex;
  flex-direction: column;
  gap: 12px;
  transition: box-shadow 0.2s;
}

.room-card:hover {
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
}

.room-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.room-title {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.room-card-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.room-info {
  font-size: 13px;
  color: #606266;
}

.info-label {
  color: #909399;
  margin-right: 8px;
}

.room-stats {
  display: flex;
  gap: 24px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: 18px;
  font-weight: 700;
  color: #303133;
}

.stat-label {
  font-size: 12px;
  color: #909399;
}

.room-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.room-id {
  font-size: 12px;
  color: #c0c4cc;
}
</style>
