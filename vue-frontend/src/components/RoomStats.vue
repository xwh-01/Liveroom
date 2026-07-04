<template>
  <div class="card-panel room-stats">
    <div class="panel-title">房间状态</div>
    <div class="stats-grid">
      <div class="stat-item">
        <div class="stat-label">房间号</div>
        <div class="stat-value">{{ roomId }}</div>
      </div>
      <div class="stat-item">
        <div class="stat-label">在线人数</div>
        <div class="stat-value highlight">{{ onlineCount }}</div>
      </div>
      <div class="stat-item">
        <div class="stat-label">限流次数</div>
        <div class="stat-value warn">{{ limitedCount }}</div>
      </div>
      <div class="stat-item">
        <div class="stat-label">连接状态</div>
        <div class="stat-value">
          <el-tag :type="connected ? 'success' : 'info'" size="small">
            {{ connected ? '已连接' : '未连接' }}
          </el-tag>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  roomId: { type: String, default: '-' },
  onlineCount: { type: Number, default: 0 },
  limitedCount: { type: Number, default: 0 },
  connected: { type: Boolean, default: false }
})
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.stat-item {
  text-align: center;
}

.stat-label {
  font-size: 12px;
  color: #909399;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
}

.stat-value.highlight {
  color: #409eff;
}

.stat-value.warn {
  color: #e6a23c;
}
</style>
