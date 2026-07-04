<template>
  <div class="card-panel danmu-panel">
    <div class="panel-title">弹幕区</div>
    <div class="danmu-list" ref="listRef">
      <div v-if="danmuList.length === 0" class="danmu-empty">暂无弹幕</div>
      <div
        v-for="(item, idx) in danmuList"
        :key="idx"
        class="danmu-item"
        :class="{ 'danmu-enter': item.justEntered }"
        :style="{ color: item.color }"
      >
        <span class="danmu-time">{{ item.time }}</span>
        <span class="danmu-user">{{ item.user }}</span>
        <span class="danmu-content">{{ item.content }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'

const danmuList = ref([])
const listRef = ref(null)
const MAX_DANMU = 100

const colors = ['#e74c3c', '#3498db', '#2ecc71', '#f39c12', '#9b59b6', '#1abc9c', '#e67e22', '#e91e63']

function addDanmu(msg) {
  const item = {
    time: msg.data?.timestamp || '',
    user: msg.user_id || 'unknown',
    content: msg.data?.content || '',
    color: colors[Math.floor(Math.random() * colors.length)],
    justEntered: true
  }
  danmuList.value.push(item)
  if (danmuList.value.length > MAX_DANMU) {
    danmuList.value.shift()
  }
  scrollToBottom()
  setTimeout(() => { item.justEntered = false }, 500)
}

function scrollToBottom() {
  nextTick(() => {
    if (listRef.value) {
      listRef.value.scrollTop = listRef.value.scrollHeight
    }
  })
}

defineExpose({ addDanmu })
</script>

<style scoped>
.danmu-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.danmu-list {
  flex: 1;
  overflow-y: auto;
  max-height: 260px;
  font-size: 13px;
  line-height: 1.8;
}

.danmu-list::-webkit-scrollbar { width: 4px; }
.danmu-list::-webkit-scrollbar-thumb { background: #c0c4cc; border-radius: 2px; }

.danmu-empty {
  color: #c0c4cc;
  text-align: center;
  margin-top: 40px;
  font-size: 13px;
}

.danmu-item {
  padding: 3px 0;
  border-bottom: 1px solid #f5f5f5;
  word-break: break-all;
}

.danmu-item.danmu-enter {
  animation: danmuIn 0.3s ease-out;
}

.danmu-time {
  color: #c0c4cc;
  font-size: 11px;
  margin-right: 6px;
}

.danmu-user {
  font-weight: 600;
  margin-right: 6px;
}

@keyframes danmuIn {
  from { opacity: 0; transform: translateX(-20px); }
  to { opacity: 1; transform: translateX(0); }
}
</style>
