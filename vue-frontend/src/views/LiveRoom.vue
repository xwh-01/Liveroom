<template>
  <div class="live-room">
    <header class="live-room-header">LiveRoom Battle</header>

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
        <el-input v-model="roomId" placeholder="Room ID" size="small" style="width: 100px" />
        <el-input v-model="userId" placeholder="User ID" size="small" style="width: 120px" />
        <el-button
          :type="connected ? 'danger' : 'primary'"
          size="small"
          @click="connected ? disconnect() : connect()"
        >
          {{ connected ? '断开' : '连接' }}
        </el-button>
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
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount } from 'vue'
import { ElMessage } from 'element-plus'
import { LiveWSClient } from '../utils/ws.js'
import LivePlayer from '../components/LivePlayer.vue'
import DanmuPanel from '../components/DanmuPanel.vue'
import GiftPanel from '../components/GiftPanel.vue'
import RankBoard from '../components/RankBoard.vue'
import RoomStats from '../components/RoomStats.vue'
import SystemLog from '../components/SystemLog.vue'

const ws = new LiveWSClient()

const roomId = ref('1001')
const userId = ref('user_' + Math.random().toString(36).slice(2, 8))
const connected = ref(false)
const onlineCount = ref(0)
const limitedCount = ref(0)
const rankings = ref([])
const chatText = ref('')

const danmuPanelRef = ref(null)
const systemLogRef = ref(null)

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

function connect() {
  if (!roomId.value || !userId.value) {
    ElMessage.warning('请输入 Room ID 和 User ID')
    return
  }
  ws.connect(roomId.value, userId.value)
}

function disconnect() {
  onlineCount.value = 0
  rankings.value = []
  ws.disconnect()
}

function sendChat() {
  const text = chatText.value.trim()
  if (!text) return
  ws.sendChat(roomId.value, userId.value, text)
  chatText.value = ''
}

function sendGift(giftType) {
  ws.sendGift(roomId.value, userId.value, giftType)
}

async function fetchRoomState() {
  try {
    const res = await fetch(`/api/room/state?room_id=${roomId.value}`)
    const data = await res.json()
    if (data.code === 0) {
      onlineCount.value = data.data?.online_count || 0
      limitedCount.value = data.data?.limited_count || 0
    }
  } catch (e) {
    // ignore
  }
}

onBeforeUnmount(() => {
  ws.disconnect()
  ws.removeAllListeners()
})
</script>
