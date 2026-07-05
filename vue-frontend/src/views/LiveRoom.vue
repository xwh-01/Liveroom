<template>
  <div class="live-room">
    <header class="live-room-header">
      <span class="header-back" @click="$router.push('/rooms')">◂ 大厅</span>
      <span class="header-title">{{ roomTitle }}</span>
      <span class="header-anchor">{{ roomAnchor }}</span>
    </header>

    <div v-if="roomError" class="room-error">
      <p>{{ roomError }}</p>
      <el-button type="primary" @click="$router.push('/rooms')">返回大厅</el-button>
    </div>

    <template v-else>
      <div class="live-room-body">
        <div class="live-room-left">
          <LivePlayer :default-room="roomId" />
          <DanmuPanel ref="danmuPanelRef" />
        </div>

        <div class="live-room-right">
          <RoomStats
            :online-count="onlineCount"
            :limited-count="limitedCount"
            :connected="connected"
            :room-id="roomId"
          />
          <RankBoard :rankings="rankings" />
          <SystemLog ref="systemLogRef" />
        </div>
      </div>

      <footer class="live-room-footer">
        <div class="connect-group">
          <el-button
            :type="connected ? 'danger' : 'primary'"
            size="small"
            @click="connected ? disconnect() : connect()"
          >
            {{ connected ? '断开' : '连接' }}
          </el-button>
          <span class="room-id-label">房间: {{ roomId }}</span>
          <span class="user-id-label">用户: {{ userId }}</span>
        </div>

        <div class="chat-group">
          <el-input
            v-model="chatText"
            placeholder="输入弹幕..."
            size="small"
            @keyup.enter="sendChat"
            :disabled="!connected"
          />
          <el-button type="primary" size="small" @click="sendChat" :disabled="!connected">
            发送
          </el-button>
        </div>

        <GiftPanel :connected="connected" @send-gift="sendGift" />
      </footer>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { ElMessage } from 'element-plus'
import { LiveWSClient } from '../utils/ws.js'
import LivePlayer from '../components/LivePlayer.vue'
import DanmuPanel from '../components/DanmuPanel.vue'
import GiftPanel from '../components/GiftPanel.vue'
import RankBoard from '../components/RankBoard.vue'
import RoomStats from '../components/RoomStats.vue'
import SystemLog from '../components/SystemLog.vue'

const props = defineProps({
  roomId: { type: String, required: true }
})

const ws = new LiveWSClient()

const roomTitle = ref('加载中...')
const roomAnchor = ref('')
const roomError = ref('')
const userId = ref(getOrCreateUserId())
const connected = ref(false)
const onlineCount = ref(0)
const limitedCount = ref(0)
const rankings = ref([])
const chatText = ref('')

const danmuPanelRef = ref(null)
const systemLogRef = ref(null)

function getOrCreateUserId() {
  const stored = localStorage.getItem('liveroom_user_id')
  if (stored) return stored
  const id = 'guest_' + Math.random().toString(36).slice(2, 10)
  localStorage.setItem('liveroom_user_id', id)
  return id
}

ws.on('open', () => {
  connected.value = true
  ElMessage.success('已连接直播间')
  systemLogRef.value?.addLog('system', `用户 ${userId.value} 进入直播间`)
  fetchRoomState()
})

ws.on('close', () => {
  connected.value = false
  ElMessage.info('已断开连接')
})

ws.on('chat', (msg) => {
  danmuPanelRef.value?.addDanmu(msg)
})

ws.on('gift', (msg) => {
  const giftLabel = msg.data?.gift_type === 'rocket' ? '火箭 (100分)' : '小心心 (10分)'
  systemLogRef.value?.addLog('gift', `${msg.data?.sender || msg.user_id} 送出 ${giftLabel}`)
})

ws.on('rank', (msg) => {
  rankings.value = msg.data?.rankings || []
})

ws.on('online', (msg) => {
  onlineCount.value = msg.data?.count || 0
})

ws.on('system', (msg) => {
  const content = msg.data?.content || ''
  systemLogRef.value?.addLog('system', content)
  if (content.includes('限流')) {
    fetchRoomState()
  }
})

async function fetchRoomInfo() {
  try {
    const res = await fetch(`/api/rooms/${props.roomId}`)
    const data = await res.json()
    if (data.code !== 0 || !data.data) {
      roomError.value = '房间不存在或已关闭'
      return
    }
    if (data.data.status === 'closed') {
      roomError.value = '房间已关闭'
      return
    }
    roomTitle.value = data.data.title || '直播间'
    roomAnchor.value = data.data.anchor_name ? `主播: ${data.data.anchor_name}` : ''
  } catch (e) {
    roomError.value = '无法获取房间信息'
  }
}

function connect() {
  ws.connect(props.roomId, userId.value)
}

function disconnect() {
  onlineCount.value = 0
  rankings.value = []
  ws.disconnect()
}

function sendChat() {
  const text = chatText.value.trim()
  if (!text) return
  ws.sendChat(props.roomId, userId.value, text)
  chatText.value = ''
}

function sendGift(giftType) {
  ws.sendGift(props.roomId, userId.value, giftType)
}

async function fetchRoomState() {
  try {
    const res = await fetch(`/api/room/state?room_id=${props.roomId}`)
    const data = await res.json()
    if (data.code === 0) {
      onlineCount.value = data.data?.online_count || 0
      limitedCount.value = data.data?.limited_count || 0
    }
  } catch (e) {
    // ignore
  }
}

onMounted(async () => {
  await fetchRoomInfo()
  if (!roomError.value) {
    connect()
  }
})

onBeforeUnmount(() => {
  ws.disconnect()
  ws.removeAllListeners()
})
</script>

<style scoped>
.live-room-header {
  position: relative;
  justify-content: center;
}
.header-back {
  position: absolute;
  left: 24px;
  cursor: pointer;
  font-size: 14px;
  opacity: 0.8;
  transition: opacity 0.2s;
}
.header-back:hover {
  opacity: 1;
}
.header-title {
  font-size: 16px;
}
.header-anchor {
  position: absolute;
  right: 24px;
  font-size: 13px;
  opacity: 0.8;
}

.room-error {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  color: #909399;
  font-size: 14px;
}

.room-id-label, .user-id-label {
  font-size: 12px;
  color: #909399;
  margin-left: 8px;
}
</style>
