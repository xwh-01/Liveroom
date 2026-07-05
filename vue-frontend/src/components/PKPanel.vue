<template>
  <div class="card-panel pk-panel">
    <div class="panel-title">红蓝 PK</div>
    <div v-if="!pkState || pkState.status === 'waiting'" class="pk-waiting">
      <p>等待 PK 开始...</p>
      <div class="pk-team-select">
        <el-button type="danger" :disabled="userTeam" @click="joinTeam('red')">
          {{ userTeam === 'red' ? '🔴 你已加入红队' : '🔴 加入红队' }}
        </el-button>
        <el-button type="primary" :disabled="userTeam" @click="joinTeam('blue')">
          {{ userTeam === 'blue' ? '🔵 你已加入蓝队' : '🔵 加入蓝队' }}
        </el-button>
      </div>
    </div>
    <div v-else class="pk-active">
      <div class="pk-score-bar">
        <div class="pk-red-bar" :style="{ width: redPercent + '%' }">
          <span class="pk-red-score">{{ pkState.red_score }}</span>
        </div>
        <div class="pk-blue-bar" :style="{ width: bluePercent + '%' }">
          <span class="pk-blue-score">{{ pkState.blue_score }}</span>
        </div>
      </div>
      <div class="pk-info">
        <div class="pk-team-info red">
          <span class="team-label">🔴 红队</span>
          <span class="team-users">{{ pkState.red_users }} 人</span>
        </div>
        <div class="pk-status-info">
          <span v-if="pkState.status === 'running' && pkState.remaining_seconds > 0" class="pk-timer">
            剩余 {{ formatTime(pkState.remaining_seconds) }}
          </span>
          <span v-else-if="pkState.status === 'ended'" class="pk-result">
            <template v-if="pkState.winner === 'red'">🔴 红队胜利！</template>
            <template v-else-if="pkState.winner === 'blue'">🔵 蓝队胜利！</template>
            <template v-else>⚪ 平局</template>
          </span>
        </div>
        <div class="pk-team-info blue">
          <span class="team-label">🔵 蓝队</span>
          <span class="team-users">{{ pkState.blue_users }} 人</span>
        </div>
      </div>
      <div v-if="!userTeam && pkState.status === 'running'" class="pk-team-select">
        <el-button type="danger" size="small" @click="joinTeam('red')">🔴 加入红队</el-button>
        <el-button type="primary" size="small" @click="joinTeam('blue')">🔵 加入蓝队</el-button>
      </div>
      <div v-else-if="userTeam" class="pk-my-team">
        你正在支持
        <span :class="userTeam === 'red' ? 'red-text' : 'blue-text'">
          {{ userTeam === 'red' ? '红队' : '蓝队' }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  ws: { type: Object, required: true },
  roomId: { type: String, required: true },
  userId: { type: String, required: true },
  connected: { type: Boolean, default: false }
})

const emit = defineEmits(['team-changed'])

const pkState = ref(null)
const userTeam = ref('')

const redPercent = computed(() => {
  if (!pkState.value) return 50
  const total = pkState.value.red_score + pkState.value.blue_score
  if (total === 0) return 50
  return Math.round((pkState.value.red_score / total) * 100)
})

const bluePercent = computed(() => {
  if (!pkState.value) return 50
  const total = pkState.value.red_score + pkState.value.blue_score
  if (total === 0) return 50
  return Math.round((pkState.value.blue_score / total) * 100)
})

function joinTeam(team) {
  if (!props.connected) return
  props.ws.send({ type: 'join_team', room_id: props.roomId, user_id: props.userId, data: { team } })
}

function formatTime(seconds) {
  const m = Math.floor(seconds / 60)
  const s = seconds % 60
  return `${m}:${String(s).padStart(2, '0')}`
}

function onPKState(state) {
  pkState.value = state
}

function onSystem(msg) {
  const content = msg.data?.content || ''
  if (content.includes('已加入红队')) {
    userTeam.value = 'red'
    localStorage.setItem(`pk_team_${props.roomId}`, 'red')
    emit('team-changed', 'red')
  } else if (content.includes('已加入蓝队')) {
    userTeam.value = 'blue'
    localStorage.setItem(`pk_team_${props.roomId}`, 'blue')
    emit('team-changed', 'blue')
  }
}

function init() {
  const saved = localStorage.getItem(`pk_team_${props.roomId}`)
  if (saved) userTeam.value = saved
}

init()

defineExpose({ onPKState, onSystem, userTeam })
</script>

<style scoped>
.pk-panel { }
.pk-waiting { text-align: center; padding: 16px 0; color: #909399; font-size: 13px; }
.pk-team-select { display: flex; gap: 8px; justify-content: center; margin-top: 8px; }

.pk-score-bar {
  display: flex;
  height: 28px;
  border-radius: 14px;
  overflow: hidden;
  margin-bottom: 8px;
  background: #f0f0f0;
}
.pk-red-bar {
  background: linear-gradient(90deg, #ff4757, #ff6b81);
  display: flex;
  align-items: center;
  justify-content: flex-start;
  padding-left: 8px;
  transition: width 0.3s;
}
.pk-blue-bar {
  background: linear-gradient(90deg, #3742fa, #5352ed);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding-right: 8px;
  transition: width 0.3s;
}
.pk-red-score, .pk-blue-score { color: #fff; font-weight: 700; font-size: 13px; }

.pk-info { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; font-size: 12px; }
.pk-team-info.red { color: #ff4757; }
.pk-team-info.blue { color: #3742fa; }
.team-label { font-weight: 600; margin-right: 4px; }
.pk-timer { color: #ffa502; font-weight: 600; }
.pk-result { font-weight: 700; font-size: 13px; }
.pk-my-team { text-align: center; font-size: 12px; color: #909399; margin-top: 4px; }
.red-text { color: #ff4757; font-weight: 600; }
.blue-text { color: #3742fa; font-weight: 600; }
</style>
