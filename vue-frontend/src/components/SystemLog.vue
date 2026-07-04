<template>
  <div class="card-panel system-log">
    <div class="panel-title">系统消息</div>
    <div class="log-list" ref="logRef">
      <div v-if="logs.length === 0" class="log-empty">暂无消息</div>
      <div
        v-for="(item, idx) in logs"
        :key="idx"
        class="log-item"
        :class="'log-' + item.type"
      >
        <span class="log-time">{{ item.time }}</span>
        <span class="log-text">{{ item.text }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'

const logs = ref([])
const logRef = ref(null)
const MAX_LOGS = 50

function addLog(type, text) {
  const now = new Date()
  const time = now.toTimeString().slice(0, 8)
  logs.value.push({ type, text, time })
  if (logs.value.length > MAX_LOGS) {
    logs.value.shift()
  }
  nextTick(() => {
    if (logRef.value) {
      logRef.value.scrollTop = logRef.value.scrollHeight
    }
  })
}

defineExpose({ addLog })
</script>

<style scoped>
.system-log {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.log-list {
  flex: 1;
  overflow-y: auto;
  max-height: 200px;
  font-size: 12px;
  line-height: 1.8;
}

.log-list::-webkit-scrollbar { width: 4px; }
.log-list::-webkit-scrollbar-thumb { background: #c0c4cc; border-radius: 2px; }

.log-empty {
  color: #c0c4cc;
  text-align: center;
  margin-top: 20px;
}

.log-item {
  padding: 2px 0;
  border-bottom: 1px dashed #f0f0f0;
}

.log-time {
  color: #c0c4cc;
  margin-right: 8px;
}

.log-gift .log-text { color: #e6a23c; }
.log-system .log-text { color: #909399; }
.log-warn .log-text { color: #f56c6c; }
</style>
